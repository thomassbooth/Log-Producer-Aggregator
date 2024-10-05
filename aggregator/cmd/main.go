package main

import (
	"fmt"
	"log"
	"log-aggregator/aggregator/api"
	"log-aggregator/aggregator/internal"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	// config.Load()

	// Initialize the worker pool with the desired number of workers (e.g., 5)
	wp := internal.NewWorkerPool(5)

	// Create a new server
	srv := api.NewServer(api.Config{}, wp)

	// Create a channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	fmt.Println("API started")

	// Wait for an interrupt signal
	<-stop
	fmt.Println("Shutting down server...")

	// Stop the worker pool gracefully
	wp.Stop()
}
