package conf

import (
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Log struct {
	Lvl string `env:"LOG_LEVEL" envDefault:"info"`
}

func (l *Log) GetLogEntry() *logrus.Entry {
	level, _ := logrus.ParseLevel(l.Lvl)

	logger := logrus.New()
	logger.SetLevel(level)

	return logrus.NewEntry(logger)
}

func (c *ConfigImpl) Log() *logrus.Entry {
	if c.log != nil {
		return c.log
	}

	log := &Log{}
	if err := env.Parse(log); err != nil {
		panic(err)
	}

	c.log = log.GetLogEntry()

	return c.log
}
