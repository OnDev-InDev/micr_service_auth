package config

import (
	"os"
)

type Config struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
}

func Load() Config {
	return Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
	}
}
