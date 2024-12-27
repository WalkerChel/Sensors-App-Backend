package server

import (
	"net"
	"time"
)

// Option represents function for assigning
// custom values to the Server.
type Option func(*Server)

// Assign a custom address to the server.
func AddAddress(host, port string) Option {
	return func(s *Server) {
		s.httpServer.Addr = net.JoinHostPort(host, port)
	}
}

// Assign a custom read timeout to the server.
func AddReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.httpServer.ReadTimeout = t
	}
}

// Assign a custom write timeout to the server.
func AddWriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.httpServer.WriteTimeout = t
	}
}

// Assign a custom shutdown timeout to the server.
func AddShutdownTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = t
	}
}
