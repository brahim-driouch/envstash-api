package server

import "net/http"

type Server struct {
	server *http.Server
	addr   string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr: s.addr,
	}
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}

	return nil
}
