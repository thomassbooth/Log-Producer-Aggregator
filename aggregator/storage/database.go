package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage holds the MongoDB client and database information
type Storage struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewStorage initializes a new Storage instance and connects to MongoDB
// reference https://stackoverflow.com/questions/71893934/how-to-connect-to-mongodb-running-inside-one-container-from-golang-app-container
func NewStorage(uri, dbName, collectionName string) (*Storage, error) {
	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(uri)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check if the connection was successful
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &Storage{
		client:     client,
		collection: collection,
	}, nil
}

// Close closes the MongoDB client connection
func (s *Storage) Close() error {
	return s.client.Disconnect(context.TODO())
}
