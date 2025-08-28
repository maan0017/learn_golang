package main

import "net/http"

type Server struct {
	ListenAddr string
}

func NewServer(addr string) *Server {
	return &Server{
		ListenAddr: addr,
	}
}

func (s *Server) StartServer() error {
	return http.ListenAndServe(s.ListenAddr, nil)
}
