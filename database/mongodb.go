package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client *mongo.Client
}

func ConnectMongoDB() (*MongoDB, error) {
	log.Println("Connecting to MongoDB!")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return &MongoDB{Client: client}, nil
}

func (db *MongoDB) Disconnect() {
	err := db.Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
