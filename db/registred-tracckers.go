package db

import (
	"context"

	"github.com/gps/gps-traking/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *DB) FindRegistredTrackerByIMEI(imei string) (*models.RegistredTracker, error) {
	filter := bson.D{
		{"imei", imei},
	}

	var tracker models.RegistredTracker

	err := d.db.Collection("registred-trackers").FindOne(context.TODO(), filter).Decode(&tracker)
	if err != nil {
		return nil, err
	}

	return &tracker, nil
}
