package models

import "time"

type (
	LocalConfig struct {
		Server   ServerConfig   `yaml:"server"`
		Postgres PostgresConfig `yaml:"postgres"`
	}

	ServerConfig struct {
		Addr     string         `yaml:"addr"`
		Password string         `yaml:"password"`
		Timeouts TimeoutsConfig `yaml:"timeouts"`
	}
	TimeoutsConfig struct {
		Write time.Duration `yaml:"write"`
		Read  time.Duration `yaml:"read"`
	}

	PostgresConfig struct {
		Addr string `yaml:"addr"`
		User string `yaml:"user"`
		DB   string `yaml:"db"`
	}
)
