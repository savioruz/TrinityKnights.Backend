package config

import (
	"github.com/sirupsen/logrus"
)

// NewLogrus creates a new logrus logger
func NewLogrus() *logrus.Logger {
	log := logrus.New()

	log.SetLevel(6)
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
