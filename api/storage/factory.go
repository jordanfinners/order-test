package storage

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is a storage client for mongodb
type Client struct {
	db *mongo.Database
}

func new(databaseName string) Client {
	databaseConnectionString, ok := os.LookupEnv("DATABASE_CONNECTION_STRING")
	if !ok {
		log.Fatalf("Missing environment variable DATABASE_CONNECTION_STRING")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(databaseConnectionString).SetDirect(true)
	c, err := mongo.NewClient(clientOptions)

	err = c.Connect(ctx)

	if err != nil {
		log.Fatalf("Unable to initialize connection %v", err)
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Unable to connect %v", err)
	}
	return Client{
		db: c.Database(databaseName),
	}
}

// NewClient creates a fresh database client
func NewClient() Client {
	databaseName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		log.Fatalf("Missing environment variable DATABASE_NAME")
	}

	return new(databaseName)
}
