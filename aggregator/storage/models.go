package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

// LogMessage represents a log entry in the database
type LogMessage struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Message string             `bson:"message"`
	Level   string             `bson:"level"`
	Time    primitive.DateTime `bson:"time"`
}
