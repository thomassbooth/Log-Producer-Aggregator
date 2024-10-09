package storage

import (
	"context"
	"fmt"
	"log-aggregator/aggregator/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertLogMessages inserts multiple LogMessages into the MongoDB collection
func (s *Storage) InsertLogMessages(logs []utils.LogMessage) error {
	var logEntries []interface{} // Create a slice to hold the log entries

	// Iterate over the logs and create LogEntry documents
	for _, log := range logs {
		logEntry := LogEntry{
			ID:      primitive.NewObjectID(),
			Message: log.Message,
			Level:   log.Level,
			Time:    primitive.NewDateTimeFromTime(log.Timestamp.UTC()),
		}
		logEntries = append(logEntries, logEntry) // Append each log entry to the slice
	}

	// Insert all log entries into the collection
	_, err := s.collection.InsertMany(context.TODO(), logEntries)
	if err != nil {
		return fmt.Errorf("failed to insert log messages: %v", err)
	}
	return nil
}

// GetLogMessages retrieves log messages from the collection filtered by time range and log level
func (s *Storage) GetLogMessages(startTime, endTime time.Time, logLevel string) ([]utils.LogMessage, error) {
	var utilsLogs []utils.LogMessage
	// Create the filter based on the provided parameters

	filter, err := s.buildFilter(startTime, endTime, logLevel)
	if err != nil {
		return nil, err
	}

	fmt.Println("Filter:", filter)

	// Find log messages with the specified filter
	cursor, err := s.collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Failed to find log messages")
		return nil, fmt.Errorf("failed to find log messages: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Decode each log message and convert it to utils.LogMessage in one loop
	for cursor.Next(context.TODO()) {
		var log LogEntry
		if err := cursor.Decode(&log); err != nil {
			return nil, fmt.Errorf("failed to decode log message: %v", err)
		}
		// Directly append to utilsLogs
		utilsLogs = append(utilsLogs, utils.LogMessage{
			Timestamp: log.Time.Time(), // Format timestamp if needed
			Level:     log.Level,
			Message:   log.Message,
		})
	}

	fmt.Println(utilsLogs)
	// Check for any cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return utilsLogs, nil
}

// buildFilter constructs a filter for log messages based on the provided time range and log level
func (s *Storage) buildFilter(startTime, endTime time.Time, logLevel string) (bson.D, error) {
	filter := bson.D{}

	// Check if startTime and endTime are provided and not zero values
	if !startTime.IsZero() && !endTime.IsZero() {
		// Create a filter for the time range
		filter = append(filter, bson.E{Key: "time", Value: bson.D{
			{Key: "$gte", Value: primitive.NewDateTimeFromTime(startTime)},
			{Key: "$lte", Value: primitive.NewDateTimeFromTime(endTime)},
		}})
	}

	// Check if logLevel is provided and not empty
	if logLevel != "" {
		// Add a filter for the log level
		filter = append(filter, bson.E{Key: "level", Value: logLevel})
	}

	// If no filters are applied, retrieve all logs
	if len(filter) == 0 {
		filter = bson.D{{}} // Empty filter to match all documents
	}

	return filter, nil
}
