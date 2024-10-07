package main

import (
	"log"
	"log-aggregator/aggregator/api"
	"log-aggregator/aggregator/storage"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := storage.NewStorage("mongodb://mongodb:27017", "logdb", "logs")
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	logs, _ := db.GetLogMessages()

	defer db.Close() // Ensure MongoDB connection is closed on shutdown
	// Notify the channel when an interrupt or terminate signal is received
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// Create the server
	server := api.NewServer(api.Config{})
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
	log.Println("Server stopped successfully")
}
