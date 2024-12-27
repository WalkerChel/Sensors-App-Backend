package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultHost            = "localhost"
	defaultPort            = "8080"
	defaultShutdownTimeout = 3 * time.Second
)

// Server represents server to handle requests.
type Server struct {
	httpServer      *http.Server
	notifier        chan error
	shutdownTimeout time.Duration
}

// Returns a New instance of Server.
func New(handler http.Handler, options ...Option) Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         net.JoinHostPort(defaultHost, defaultPort),
	}

	errChan := make(chan error, 1)

	s := Server{
		httpServer:      httpServer,
		notifier:        errChan,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range options {
		opt(&s)
	}

	return s
}

// Starts server in seperate goroutine.
func (s *Server) Start() {
	go func() {
		log.Printf("Server started on %s", s.GetAdrr())
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.notifier <- err
		}
	}()
}

// Returns servers error notify channel.
func (s *Server) Notify() <-chan error {
	return s.notifier
}

// Closes Server.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// Returns the address of the server.
func (s *Server) GetAdrr() string {
	return s.httpServer.Addr
}
