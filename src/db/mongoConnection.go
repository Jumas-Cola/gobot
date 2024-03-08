package db

import (
	"context"
	"log"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDbClient *mongo.Client

func ConnectToMongoDB() {
	connectURI := os.Getenv("MONGODB_CONNECT_URI")
	clientOptions := options.Client().ApplyURI(connectURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connected to MongoDB!")
	MongoDbClient = client
}
