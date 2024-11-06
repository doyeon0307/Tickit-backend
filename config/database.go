package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Database, error) {
	connPattern := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(connPattern)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(5000)*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database("tickit"), err
}
