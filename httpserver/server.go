package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type server struct {
	http.Server
	log *logrus.Logger
}

func NewServer(log *logrus.Logger, handler http.Handler, opts ...Option) *server {
	srv := server{
		Server: http.Server{
			Handler: handler,
		},
		log: log,
	}

	for _, opt := range opts {
		opt(&srv)
	}

	return &srv
}

func (s *server) Start() {
	go func() {
		s.log.Infof("Server is running on %s", s.Addr)
		err := s.ListenAndServe() // NOTE: always returns a non-nil error.
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("Server graceful shutdown.")
		} else {
			s.log.Error("Server shutdown error:", err)
		}
	}()
}

func (s *server) Stop(ctx context.Context) {
	err := s.Shutdown(ctx)
	if err != nil {
		s.log.Error("Server shutdown error:", err)
	}
}

// NOTE: Options

type Option func(*server)

func WithPort(port int) Option {
	return func(srv *server) {
		srv.Addr = fmt.Sprintf("%s:%d", srv.Addr, port)
	}
}
