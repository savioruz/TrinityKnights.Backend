package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewLogrus creates a new logrus logger
func NewLogrus(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	logLevel := viper.GetInt("APP_LOG_LEVEL")
	log.SetLevel(logrus.Level(logLevel))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
