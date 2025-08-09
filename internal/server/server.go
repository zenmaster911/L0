package server

import (
	"net/http"
	"time"
)

type Server struct {
	HttpServer *http.Server
}

func (s *Server) run(port string, handler http.Handler) error {
	s.HttpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return s.HttpServer.ListenAndServe()
}
