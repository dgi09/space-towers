package server

import (
	"fmt"
	"net/http"
)

type Opts struct {
	Port    int
	Binding Binding
}

type Server struct {
	port    int
	ws      *ws
	binding Binding
}

func New(opts Opts) *Server {
	s := &Server{
		port:    opts.Port,
		binding: opts.Binding,
	}
	ws := newWS(WsOpts{
		MsgHandler:              s.wsMsgHandler,
		ConnectionOpenedHandler: s.wsOpenConnHandler,
		ConnectionClosedHandler: s.wsCloseConnHandler,
	})

	s.ws = ws

	return s
}

func (s *Server) wsMsgHandler(id string, data []byte) error {
	s.binding.OnPlayerMessage(id, data)

	return nil
}

func (s *Server) wsOpenConnHandler(id string) error {
	s.binding.OnPlayerConnected(id)
	return nil
}

func (s *Server) wsCloseConnHandler(id string) error {
	s.binding.OnPlayerDisconnected(id)

	return nil
}

func (s *Server) GetNetwork() *WsNetwork {
	return NewWsNetwork(s.ws)
}

func (s *Server) Start() error {
	fmt.Println("Starting server...")

	s.binding.Init(s.GetNetwork())
	s.ws.Start()

	s.binding.Start()

	http.HandleFunc("/", http.FileServer(http.Dir("public")).ServeHTTP)
	http.HandleFunc("/ws", s.handleWS)

	host := "192.168.100.9"
	return http.ListenAndServe(fmt.Sprint(host, ":", s.port), nil)
}
