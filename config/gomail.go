package config

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type SMTP struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func NewGomail(viper *viper.Viper) *gomail.Dialer {
	smtp := &SMTP{
		Host:     viper.GetString("SMTP_HOST"),
		Port:     viper.GetInt("SMTP_PORT"),
		Username: viper.GetString("SMTP_USERNAME"),
		Password: viper.GetString("SMTP_PASSWORD"),
	}

	return gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
}

