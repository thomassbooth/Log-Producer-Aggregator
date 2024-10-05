package shared

import "time"

// LogMessage represents a log entry
type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}
