package postgres

import (
	"time"
)

type UserModel struct {
	ID           string    `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Created_At   time.Time `gorm:"not null"`
	Role         string    `gorm:"not null"`
}
