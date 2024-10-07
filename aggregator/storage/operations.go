package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InsertLogMessage inserts a LogMessage into the MongoDB collection
func (s *Storage) InsertLogMessage(message, level string) error {
	// Create a new LogMessage document
	logMsg := LogMessage{
		ID:      primitive.NewObjectID(),
		Message: message,
		Level:   level,
		Time:    primitive.NewDateTimeFromTime(time.Now()),
	}

	// Insert the log message into the collection
	_, err := s.collection.InsertOne(context.TODO(), logMsg)
	if err != nil {
		return fmt.Errorf("failed to insert log message: %v", err)
	}

	return nil
}

// GetLogMessages retrieves all log messages from the collection
func (s *Storage) GetLogMessages() ([]LogMessage, error) {
	var logs []LogMessage

	// Find all log messages
	cursor, err := s.collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find log messages: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Decode each log message and append it to the list
	for cursor.Next(context.TODO()) {
		var log LogMessage
		if err := cursor.Decode(&log); err != nil {
			return nil, fmt.Errorf("failed to decode log message: %v", err)
		}
		logs = append(logs, log)
	}

	// Check for any cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return logs, nil
}
