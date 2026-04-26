package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
}

func Load() Config {
	// грузим .env только если не прод
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("no .env file found (ok for production)")
		}
	}

	cfg := Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
	}

	validate(cfg)

	return cfg
}

func validate(cfg Config) {
	if cfg.DBHost == "" ||
		cfg.DBPort == "" ||
		cfg.DBUser == "" ||
		cfg.DBName == "" {
		log.Fatal("database config is not set properly (env variables missing)")
	}
}
