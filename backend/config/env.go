package config

import "github.com/kelseyhightower/envconfig"

// Config represents application configuration.
type Config struct {
	Database DatabaseConfig
}

// DatabaseConfig bundles database related environment variables.
type DatabaseConfig struct {
	User     string `envconfig:"DB_USERNAME" default:"root"`
	Password string `envconfig:"DB_PASSWORD" default:"password"`
	Host     string `envconfig:"DB_HOST" default:"db"`
	Port     int    `envconfig:"DB_PORT" default:"3306"`
	Name     string `envconfig:"DB_DATABASE" default:"test"`
}

// Load reads environment variables into Config using envconfig.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", &cfg.Database); err != nil {
		return nil, err
	}
	return cfg, nil
}
