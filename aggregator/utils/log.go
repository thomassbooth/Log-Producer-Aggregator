package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ValidateRequest checks if the request method is valid.
func ValidateRequest(w http.ResponseWriter, r *http.Request, expectedMethod string) error {
	if r.Method != expectedMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return fmt.Errorf("invalid request method: %s, expected: %s", r.Method, expectedMethod)
	}
	return nil
}

// DecodeJSON decodes a JSON request body into the provided structure.
func DecodeJSON(body io.ReadCloser, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}

// RespondWithJSON prepares a JSON response with the specified status code.
func RespondWithJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Write the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// ParseLogQueryParams extracts and validates the startTime, endTime, and logLevel from the query parameters.
func ParseLogQueryParams(r *http.Request) (time.Time, time.Time, string, error) {
	queryParams := r.URL.Query()

	// Get and validate the startTime parameter
	startTimeStr := queryParams.Get("startTime")
	var startTime time.Time
	var err error
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return time.Time{}, time.Time{}, "", errors.New("invalid startTime. Expected format: RFC3339")
		}
	}

	// Get and validate the endTime parameter if found
	endTimeStr := queryParams.Get("endTime")
	var endTime time.Time
	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return time.Time{}, time.Time{}, "", errors.New("invalid endTime. Expected format: RFC3339")
		}
	}
	logLevel := queryParams.Get("logLevel")

	return startTime, endTime, logLevel, nil
}
