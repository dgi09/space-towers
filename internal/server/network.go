package server

type WsNetwork struct {
	ws *ws
}

func NewWsNetwork(ws *ws) *WsNetwork {
	return &WsNetwork{
		ws: ws,
	}
}

func (n *WsNetwork) SendPkg(connectionId string, data []byte) {
	n.ws.SendMsg(connectionId, data)
}
