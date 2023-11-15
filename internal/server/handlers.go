package server

import "net/http"

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	s.ws.Accept(r, w)
}
