package config

import (
	"github.com/caarlos0/env/v6"
)

// Config menyimpan semua konfigurasi aplikasi, termasuk SMTPConfig
type Config struct {
	SMTP SMTPConfig `envPrefix:"SMTP_"`
}

// SMTPConfig menyimpan konfigurasi SMTP untuk pengiriman email
type SMTPConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"` // Host SMTP
	Port     int    `env:"PORT" envDefault:"587"`       // Port SMTP (default: 587 untuk TLS)
	Email    string `env:"EMAIL"`                       // Email pengirim
	Password string `env:"PASSWORD"`                    // Password email
}

// LoadConfig membaca konfigurasi dari environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
