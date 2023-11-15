package server

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsMsgHandler func(id string, data []byte) error

type WsConnectionHandler func(id string) error

type rawMsg struct {
	Type int
	Data []byte
}

type msg = []byte

type wsMsg struct {
	Src  string
	Data []byte
}

type connection struct {
	Id            string
	Connection    *websocket.Conn
	TermChan      chan struct{}
	TermComplChan chan struct{}
	OutChan       chan msg
}

type WsOpts struct {
	MsgHandler              WsMsgHandler
	ConnectionOpenedHandler WsConnectionHandler
	ConnectionClosedHandler WsConnectionHandler
}

type ws struct {
	upgrader    websocket.Upgrader
	connections map[string]*connection
	recChan     chan wsMsg

	msgHandler        WsMsgHandler
	connOpenedHandler WsConnectionHandler
	connClosedHandler WsConnectionHandler
}

func newWS(opts WsOpts) *ws {
	return &ws{
		connections: make(map[string]*connection),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		recChan:           make(chan wsMsg),
		msgHandler:        opts.MsgHandler,
		connOpenedHandler: opts.ConnectionOpenedHandler,
		connClosedHandler: opts.ConnectionClosedHandler,
	}
}

func (ws *ws) Start() {
	go ws.listenForMsg()
}

func (ws *ws) Stop() {
	close(ws.recChan)
}

func (ws *ws) Accept(r *http.Request, w http.ResponseWriter) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	c := &connection{
		Id:            id,
		Connection:    conn,
		TermChan:      make(chan struct{}),
		TermComplChan: make(chan struct{}),
		OutChan:       make(chan msg),
	}

	ws.connections[id] = c

	go ws.execConnection(c)

	return nil
}

func (ws *ws) CloseConnection(id string) error {
	c, ok := ws.connections[id]

	if !ok {
		return errors.New("connection not found")
	}

	c.TermChan <- struct{}{}

	<-c.TermComplChan

	return nil
}

func (ws *ws) SendMsg(id string, data []byte) error {
	c, ok := ws.connections[id]

	if !ok {
		return errors.New("connection not found")
	}

	c.OutChan <- data

	return nil
}

func (ws *ws) onConnectionClosed(c *connection) {
	delete(ws.connections, c.Id)

	ws.connClosedHandler(c.Id)
}

func (ws *ws) listenForMsg() {
	for msg := range ws.recChan {
		ws.msgHandler(msg.Src, msg.Data)
	}
}

func (ws *ws) execConnection(c *connection) {
	ws.connOpenedHandler(c.Id)

	readChan := make(chan rawMsg)

	go func(ch chan rawMsg) {
		for {
			t, msg, err := c.Connection.ReadMessage()
			if err != nil {
				close(ch)
				return
			}

			ch <- rawMsg{
				Type: t,
				Data: msg,
			}
		}
	}(readChan)

	for {
		select {
		case msg := <-c.OutChan:
			c.Connection.WriteMessage(websocket.BinaryMessage, msg)

		case <-c.TermChan:
			close(readChan)
			c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
			c.Connection.Close()
			ws.onConnectionClosed(c)
			c.TermComplChan <- struct{}{}
			return

		case msg, opened := <-readChan:
			if !opened || msg.Type == websocket.CloseMessage {
				c.Connection.Close()
				ws.onConnectionClosed(c)
				c.TermComplChan <- struct{}{}
				return
			}

			ws.recChan <- wsMsg{
				Src:  c.Id,
				Data: msg.Data,
			}
		}
	}
}
