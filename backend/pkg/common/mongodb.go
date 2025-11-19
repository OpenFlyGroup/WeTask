package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

// ? InitMongoDB initializes MongoDB connection
func InitMongoDB() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = fmt.Sprintf(
			"mongodb://%s:%s/%s",
			getEnv("MONGO_HOST", "localhost"),
			getEnv("MONGO_PORT", "27017"),
			getEnv("MONGO_DB", "kanban"),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// ? Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "kanban"
	}

	MongoDB = client.Database(dbName)
	log.Println("MongoDB connected successfully")
	return nil
}
