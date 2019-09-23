package migrations

import (
	"context"

	"github.com/gps/gps-traking/models"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		trackers := []interface{}{
			models.RegistredTracker{
				IMEI:   "087079763391",
				Type:   "chinese",
				UserID: 1,
			},
		}

		_, err := db.Collection("registred-trackers").InsertMany(context.TODO(), trackers)
		if err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		filter := bson.D{
			{"imei", "087079763391"},
		}

		_, err := db.Collection("registred-trackers").DeleteMany(context.TODO(), filter)
		if err != nil {
			return err
		}

		return nil
	})
}
