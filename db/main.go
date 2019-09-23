package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	db *mongo.Database
}

func (d DB) Mongo() *mongo.Database {
	return d.db
}

func New(link, dbName string) (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(link))
	if err != nil {
		return nil, err
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &DB{db: db}, err
}
