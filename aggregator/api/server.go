package api

import (
	"fmt"
	"log-aggregator/aggregator/internal"
	"net/http"
)

const defaultListenAddr = ":8005"

// Config holds the configuration for the server.
type Config struct {
	ListenAddr string
}

// Server struct holds the server's configuration, worker pool, and handlers.
type Server struct {
	Config
	wp             *internal.WorkerPool
	handlers       *Handlers
	circuitBreaker *internal.CircuitBreaker // Add circuit breaker field
}

// NewServer initializes a new server with the given configuration and worker pool.
func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	wp := internal.NewWorkerPool(5)
	cb := internal.NewCircuitBreaker(3, 10) // Create a new circuit breaker

	handlers := NewHandlers(wp, cb) // Pass the circuit breaker to handlers
	return &Server{Config: cfg, wp: wp, handlers: handlers, circuitBreaker: cb}
}

// Start starts the server and listens for incoming requests and signals.
func (s *Server) Start() error {
	// Setup HTTP server and routes
	srv := &http.Server{Addr: s.ListenAddr}
	http.HandleFunc("/logs", s.handlers.HandleLog)
	http.HandleFunc("/health", s.handlers.HandleHealthCheck)

	fmt.Printf("Starting server on %s\n", s.ListenAddr)
	// If the server fails to start, return the error
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %v", err)
	}

	return nil
}

// Stop gracefully stops the worker pool
func (s *Server) Stop() {
	//close the worker pool
	s.wp.Stop()
	fmt.Println("Server stopped gracefully")
}
