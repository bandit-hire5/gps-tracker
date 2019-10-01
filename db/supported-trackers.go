package db

import (
	"context"

	"github.com/gps/gps-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) FindAllSupportedTrackers() ([]*models.SupportedTracker, error) {
	opts := options.Find()
	opts.SetLimit(2)
	filter := bson.M{}

	var results []*models.SupportedTracker

	cursor, err := d.db.Collection("supported-trackers").Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var elem models.SupportedTracker

		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(context.TODO())

	return results, nil
}
