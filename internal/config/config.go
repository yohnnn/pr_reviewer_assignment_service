package config

import (
	"errors"
	"fmt"
	"os"
)

type (
	ServerConfig struct {
		Port     string
		LogLevel string
	}

	PostgresConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DB       string
	}

	Config struct {
		Server   ServerConfig
		Postgres PostgresConfig
	}
)

func Load() (*Config, error) {
	var errs []error

	required := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			errs = append(errs, fmt.Errorf("environment variable %q is required", key))
			return ""
		}
		return val
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:     required("SERVER_PORT"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
		Postgres: PostgresConfig{
			Host:     required("POSTGRES_HOST"),
			Port:     required("POSTGRES_PORT"),
			User:     required("POSTGRES_USER"),
			Password: required("POSTGRES_PASSWORD"),
			DB:       required("POSTGRES_DB"),
		},
	}

	if cfg.Server.LogLevel == "" {
		cfg.Server.LogLevel = "info"
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
}

func (p PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DB)
}
