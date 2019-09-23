package conf

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/caarlos0/env"
	"github.com/gps/gps-traking/db"
)

type DB struct {
	Name string `env:"GPS_TRACKING_DB_NAME,required"`
	Host string `env:"GPS_TRACKING_DB_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"GPS_TRACKING_DB_PORT" envDefault:"27017"`
	//User     string `env:"GPS_TRACKING_DB_USER,required"`
	//Password string `env:"GPS_TRACKING_DB_PASSWORD,required"`
}

func (d DB) Info() string {
	return fmt.Sprintf("mongodb://%s:%d", d.Host, d.Port)
}

func (d DB) GetName() string {
	return d.Name
}

func (c *ConfigImpl) DB() *db.DB {
	if c.db != nil {
		return c.db
	}

	database := &DB{}
	if err := env.Parse(database); err != nil {
		panic(err)
	}

	repo, err := db.New(database.Info(), database.GetName())
	if err != nil {
		panic(errors.Wrap(err, "failed to setup db"))
	}

	c.db = repo

	return c.db
}
