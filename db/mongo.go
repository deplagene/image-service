package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient(uri string) (*mongo.Client, error) {
	const op = "db.NewMongoClient"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client, nil
}
