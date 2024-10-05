package api

import (
	"fmt"
	"log-aggregator/aggregator/internal"
	"net/http"
)

const defaultListenAddr = ":8005"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	wp       *internal.WorkerPool
	handlers *Handlers
}

func NewServer(cfg Config, wp *internal.WorkerPool) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	handlers := NewHandlers(wp)
	return &Server{Config: cfg, wp: wp, handlers: handlers}
}

// Start starts the server
func (s *Server) Start() error {
	http.HandleFunc("/logs", s.handlers.HandleLog) // Use the handler

	fmt.Printf("Starting server on %s\n", s.ListenAddr)
	return http.ListenAndServe(s.ListenAddr, nil)
}
