package api

import (
	"fmt"
	"log-aggregator/aggregator/internal"
	"log-aggregator/aggregator/utils"
	"log-aggregator/shared"
	"net/http"
)

// Handlers struct
type Handlers struct {
	wp             *internal.WorkerPool
	circuitBreaker *internal.CircuitBreaker // Add circuit breaker field
}

// NewHandlers initializes the Handlers with a WorkerPool and CircuitBreaker
func NewHandlers(wp *internal.WorkerPool, cb *internal.CircuitBreaker) *Handlers {
	return &Handlers{
		wp:             wp,
		circuitBreaker: cb, // Initialize circuit breaker
	}
}

// HandleLog handles incoming log messages
func (h *Handlers) HandleLog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received log message")

	// Validate the request method and decode the JSON request body
	if err := utils.ValidateRequest(w, r, http.MethodPost); err != nil {
		return
	}

	var logMsg shared.LogMessage
	if err := utils.DecodeJSON(r.Body, &logMsg); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.circuitBreaker.Call(func() error {
		h.wp.AddJob(logMsg)
		return nil
	}); err != nil {
		http.Error(w, "Service unavailable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Prepare and respond with a success message
	utils.RespondWithJSON(w, http.StatusAccepted, map[string]string{"status": "success", "message": "Log message accepted"})
}

// HandleHealthCheck handles health check requests
func (h *Handlers) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	activeWorkers := h.wp.ActiveWorkers()
	queuedTasks := h.wp.QueuedTasks()

	// Basic health check logic
	if activeWorkers == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("No active workers"))
		return
	} else if queuedTasks > 100 { // Example threshold
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf("Too many queued tasks: %d", queuedTasks)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("OK - Active workers: %d, Queued tasks: %d", activeWorkers, queuedTasks)))
}
