package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/exp/rand"
)

// LogAggregatorURL is the URL of the log aggregator service
const LogAggregatorURL = "http://localhost:8005/logs/batch"

// LogMessage represents the structure of the log message
type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

// SendLog sends a batch of log messages to the log aggregator service
func SendLog(logs []LogMessage) error {
	jsonData, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", LogAggregatorURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status is not OK (200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send log, status code: " + resp.Status)
	}

	log.Printf("Logs sent successfully: %v", logs)
	return nil
}

// LogProducer generates log messages and sends them to the log aggregator
func LogProducer() {
	// Create a new random generator with a seed based on the current time
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Define log levels
	logLevels := []string{"INFO", "WARN", "ERROR"}

	// Prepare a batch of log messages
	var logs []LogMessage
	for i := 0; i < 10; i++ { // Simulate creating 10 log messages
		message := fmt.Sprintf("Log message %d", r.Intn(100))
		level := logLevels[r.Intn(len(logLevels))] // Randomly select a log level
		timestamp := time.Now().UTC()              // Get the current UTC time

		logs = append(logs, LogMessage{
			Timestamp: timestamp,
			Level:     level,
			Message:   message,
		})
	}

	operation := func() error {
		return SendLog(logs) // Send the batch of log messages
	}

	// Create an exponential backoff with custom settings if needed
	backoffStrategy := backoff.NewExponentialBackOff()
	backoffStrategy.MaxElapsedTime = 1 * time.Minute // Maximum time to retry

	// Use exponential backoff for retrying
	err := backoff.Retry(operation, backoffStrategy)
	if err != nil {
		log.Printf("Failed to send logs after retries: %v", err)
	}
}
