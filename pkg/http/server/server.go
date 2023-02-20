package server

import (
	"context"
	"miborchestrator/internal/handler"

	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler *handler.Handler) chan error {
	serverErr := make(chan error)
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler.Router,
		MaxHeaderBytes: 1 << 20, //1 mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	return serverErr
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
