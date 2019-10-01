package migrations

import (
	"context"

	"github.com/gps/gps-tracker/models"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		trackers := []interface{}{
			models.SupportedTracker{
				Name:    "Chinese Trackers",
				Type:    "chinese",
				Pattern: `^\(0[\d]{11}[\w]{2}[\d]{2}`,
			},
		}

		_, err := db.Collection("supported-trackers").InsertMany(context.TODO(), trackers)
		if err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		filter := bson.D{
			{"name", "Chinese Trackers"},
		}

		_, err := db.Collection("supported-trackers").DeleteMany(context.TODO(), filter)
		if err != nil {
			return err
		}

		return nil
	})
}
