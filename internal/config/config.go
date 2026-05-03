package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/OnDev-InDev/micr_service_auth/internal/errors"
)

type Config struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
}

func Load() (Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	cfg := Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validate(cfg Config) error {
	var missing []string

	if cfg.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DBPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return errors.Wrap(errors.CodeValidationError, fmt.Sprintf("missing required env vars: %v", missing), nil,)
	}

	return nil
}
