package conf

import (
	"github.com/gps/gps-tracker/db"
	"github.com/sirupsen/logrus"
)

type Config interface {
	Log() *logrus.Entry
	DB() *db.DB
	Tracker() *Tracker
}

type ConfigImpl struct {
	log     *logrus.Entry
	db      *db.DB
	tracker *Tracker
}

func New() Config {
	return &ConfigImpl{}
}
