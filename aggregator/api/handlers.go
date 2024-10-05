package api

import (
	"encoding/json"
	"fmt"
	"log-aggregator/aggregator/internal"
	"net/http"
)

type Handlers struct {
	wp *internal.WorkerPool
}

func NewHandlers(wp *internal.WorkerPool) *Handlers {
	return &Handlers{wp: wp}
}

// handle log endpoint
func (h *Handlers) HandleLog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received log message")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logMsg internal.LogMessage
	if err := json.NewDecoder(r.Body).Decode(&logMsg); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Add log message to the worker pool
	h.wp.AddJob(logMsg)

	// Prepare a response
	response := map[string]string{
		"status":  "success",
		"message": "Log message accepted",
	}

	// Set the response header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	// Write the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
