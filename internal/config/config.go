package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/OnDev-InDev/micr_service_auth/internal/errors"
)

type ConfigPostgres struct {
	PostgresUser     string
	PostgresPassword string
	PostgresName     string
	PostgresHost     string
	PostgresPort     string
}

type ConfigRedis struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
}

func LoadPostgres() (ConfigPostgres, error) {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	cfg := ConfigPostgres{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASS"),
		PostgresName:     os.Getenv("POSTGRES_NAME"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
	}

	if err := validate(cfg); err != nil {
		return ConfigPostgres{}, err
	}

	return cfg, nil
}

func LoadRedis() (ConfigRedis, error) {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}
	cfg := ConfigRedis{
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),
	}

	return cfg, nil
}

func validate(cfg ConfigPostgres) error {
	var missing []string

	if cfg.PostgresHost == "" {
		missing = append(missing, "POSTGRES_HOST")
	}
	if cfg.PostgresPort == "" {
		missing = append(missing, "POSTGRES_PORT")
	}
	if cfg.PostgresUser == "" {
		missing = append(missing, "POSTGRES_USER")
	}
	if cfg.PostgresName == "" {
		missing = append(missing, "POSTGRES_NAME")
	}

	if len(missing) > 0 {
		return errors.Wrap(errors.CodeValidationError, fmt.Sprintf("missing required env vars: %v", missing), nil)
	}

	return nil
}
