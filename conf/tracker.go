package conf

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Tracker struct {
	Port int `env:"GPS_TRACKING_PORT" envDefault:"8080"`
}

func (d Tracker) Info() string {
	return fmt.Sprintf(":%d", d.Port)
}

func (c *ConfigImpl) Tracker() *Tracker {
	if c.tracker != nil {
		return c.tracker
	}

	tracker := &Tracker{}
	if err := env.Parse(tracker); err != nil {
		panic(err)
	}

	c.tracker = tracker

	return c.tracker
}
