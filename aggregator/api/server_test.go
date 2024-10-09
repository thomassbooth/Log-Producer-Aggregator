package api_test

import (
	"log-aggregator/aggregator/api"
	"net/http"
	"testing"
)

func TestServerStart(t *testing.T) {
	cfg := api.Config{ListenAddr: ":8080"}
	server := api.NewServer(cfg)

	// Create a goroutine to run the server
	go func() {
		if err := server.Start(); err != nil {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	// Allow some time for the server to start
	defer func() {
		http.DefaultServeMux = http.NewServeMux() // Reset the default mux
		server.Stop()
	}()

	// Make a test request to one of the handlers
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

func TestServerStop(t *testing.T) {
	cfg := api.Config{ListenAddr: ":8081"}
	server := api.NewServer(cfg)

	go server.Start()

	// Stop the server
	server.Stop()

	// Test that the worker pool stops
	if server.Wp.ActiveWorkers() != 0 {
		t.Errorf("Expected active workers to be 0 after stopping")
	}

	// Reset the handlers
	http.DefaultServeMux = http.NewServeMux()
}
