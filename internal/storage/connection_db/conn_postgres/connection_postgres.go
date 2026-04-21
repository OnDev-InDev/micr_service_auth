package conn_postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBName,
		cfg.DBPass,
	)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres connect failed: %w", err)
	}

	sqlDB := db.DB()

	if sqlDB != nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)

		if err := sqlDB.Ping(); err != nil {
			return nil, fmt.Errorf("postgres ping failed: %w", err)
		}
	}

	return db, nil
}
