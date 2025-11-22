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
	cfg := &Config{
		Server: ServerConfig{
			Port:     os.Getenv("SERVER_PORT"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
		Postgres: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       os.Getenv("POSTGRES_DB"),
		},
	}

	if cfg.Server.LogLevel == "" {
		cfg.Server.LogLevel = "info"
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (p PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DB)
}

func (c *Config) validate() error {
	var errs []error

	if c.Server.Port == "" {
		errs = append(errs, errors.New("SERVER_PORT is required"))
	}
	if c.Postgres.Host == "" {
		errs = append(errs, errors.New("POSTGRES_HOST is required"))
	}
	if c.Postgres.Port == "" {
		errs = append(errs, errors.New("POSTGRES_PORT is required"))
	}
	if c.Postgres.User == "" {
		errs = append(errs, errors.New("POSTGRES_USER is required"))
	}
	if c.Postgres.Password == "" {
		errs = append(errs, errors.New("POSTGRES_PASSWORD is required"))
	}
	if c.Postgres.DB == "" {
		errs = append(errs, errors.New("POSTGRES_DB is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
