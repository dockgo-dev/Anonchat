package lib

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	// Config aggregates application configuration loaded from YAML.
	Config struct {
		Server      ServerConfig  `yaml:"server"`
		AuthService ServiceConfig `yaml:"authorization-sercice"`
	}

	// ServerConfig describes the gateway server settings.
	ServerConfig struct {
		Addr       string             `yaml:"addr"`
		Timeouts   ServerTimeouts     `yaml:"timeouts"`
		RateLimits RateLimitsConfig   `yaml:"rate-limits"`
		Auth       AuthConfig         `yaml:"auth"`
		Logger     ServerLoggerConfig `yaml:"logger"`
	}

	// ServerTimeouts contains read/write timeouts.
	ServerTimeouts struct {
		Write time.Duration `yaml:"write"`
		Read  time.Duration `yaml:"read"`
	}

	// RateLimitsConfig controls per-client throttling options.
	RateLimitsConfig struct {
		MaxRequests int           `yaml:"max-requests"`
		UpdateIn    time.Duration `yaml:"update-in"`
		Mode        string        `yaml:"mode"`
	}

	AuthConfig struct {
		Mode string `yaml:"mode"`
	}

	// ServerLoggerConfig toggles gateway logging behaviour.
	ServerLoggerConfig struct {
		Level string `yaml:"level"`
		Mode  string `yaml:"mode"`
	}

	// ServiceConfig declares downstream service connectivity options.
	ServiceConfig struct {
		Addr string `yaml:"addr"`
		Mode string `yaml:"mode"`
	}
)

// LoadConfig reads YAML configuration from the provided path.
func LoadConfig() (*Config, error) {
	path := "config/config.yaml"

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := new(Config)
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

// MustLoadConfig is a helper that wraps LoadConfig and panics on error.
func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println("aw, config error:", err)
		os.Exit(1)
	}

	return cfg
}
