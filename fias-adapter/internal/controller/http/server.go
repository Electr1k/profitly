package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"fias-adapter/config"
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func NewServer(cfg config.HttpServer, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         net.JoinHostPort(cfg.Address, cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}

func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
