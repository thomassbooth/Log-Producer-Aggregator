package utils

import "time"

type JobType int

const (
	FetchJob JobType = iota
	StoreJob
)

type Job struct {
	Type      JobType           `json:"type"`
	Logs      []LogMessage      `json:"logs"`
	Result    chan []LogMessage `json:"-"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	LogLevel  string            `json:"log_level"`
}

type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}
