package main

import (
	"log"
	"log-aggregator/aggregator/api"
	"os"
	"os/signal"
	"syscall"
)

var config = api.Config{
	ListenAddr: ":8005",
	DSN:        "mongodb://mongodb:27017",
}

func main() {
	// Notify the channel when an interrupt or terminate signal is received
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// Create the server
	server := api.NewServer(config)
	go func() {
		if err := server.Start(); err != nil {
			// this closes our application on Fatal error
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	// Block the main thread, waiting for a signal to close our application
	sig := <-signals
	log.Printf("Received signal: %v. Shutting down...", sig)
	server.Stop()
}
