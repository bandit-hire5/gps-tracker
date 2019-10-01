package db

import (
	"context"

	"github.com/gps/gps-tracker/models"
)

func (d *DB) AddNewTrackerLog(tracker *models.Tracker) error {
	_, err := d.db.Collection("tracker-logs").InsertOne(context.TODO(), tracker)
	if err != nil {
		return err
	}

	return nil
}
