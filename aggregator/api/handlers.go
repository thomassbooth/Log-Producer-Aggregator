package api

import (
	"fmt"
	"log-aggregator/aggregator/internal"
	"log-aggregator/aggregator/utils"
	"net/http"
	"time"
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

// HandleHealthCheck handles health check requests
func (h *Handlers) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	activeWorkers := h.wp.ActiveWorkers()
	queuedTasks := h.wp.QueuedTasks()

	// Basic health check logic
	if activeWorkers == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("No active workers"))
		return
	} else if queuedTasks > 100 {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(fmt.Sprintf("Too many queued tasks: %d", queuedTasks)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("OK - Active workers: %d, Queued tasks: %d", activeWorkers, queuedTasks)))
}

func (h *Handlers) HandleLogRetrieval(w http.ResponseWriter, r *http.Request) {
	// Validate the request method
	if err := utils.ValidateRequest(w, r, http.MethodGet); err != nil {
		return
	}

	//parses our query
	startTime, endTime, logLevel, err := utils.ParseLogQueryParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a channel to receive the result of the log retrieval
	resultChannel := make(chan []utils.LogMessage)

	// Create the fetch job with the result channel
	job := utils.Job{
		Type:      utils.FetchJob, // This job is to fetch logs
		Result:    resultChannel,
		StartTime: startTime,
		EndTime:   endTime,
		LogLevel:  logLevel,
	}

	// Add the job to the worker pool
	if err := h.circuitBreaker.Call(func() error {
		h.wp.AddJob(job)
		return nil
	}); err != nil {
		http.Error(w, "Service unavailable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Wait for the result from the worker
	select {
	case fetchedLogs := <-resultChannel:
		fmt.Println(fetchedLogs)
		if len(fetchedLogs) == 0 {
			utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"message": "No logs found"})
		} else {
			utils.RespondWithJSON(w, http.StatusOK, fetchedLogs)
		}
	case <-time.After(10 * time.Second): // Timeout to avoid long waits
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Timeout while fetching logs"})
	}

}

// Stores a batch of log messages
func (h *Handlers) HandleBatchLog(w http.ResponseWriter, r *http.Request) {
	// Validate the request method and decode the JSON request body
	fmt.Println("message recieved")
	if err := utils.ValidateRequest(w, r, http.MethodPost); err != nil {
		return
	}

	var logBatch []utils.LogMessage
	// Check if valid JSON is passed in the format we need
	if err := utils.DecodeJSON(r.Body, &logBatch); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Create a job for storing logs
	storeJob := utils.Job{
		Type:   utils.StoreJob,
		Logs:   logBatch,
		Result: nil,
	}

	if err := h.circuitBreaker.Call(func() error {
		h.wp.AddJob(storeJob)
		return nil
	}); err != nil {
		http.Error(w, "Service unavailable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Respond with a success message
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Log batch accepted"})
}
