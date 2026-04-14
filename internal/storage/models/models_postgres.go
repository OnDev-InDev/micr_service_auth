package models


import (
	"time"
)

type User struct {
	ID              uint
	Email           string     `gorm:  "primaryKey"`
	PasswordHash    string     `gorm:  "uniqueIndex;not null"`
  Created_At      time.Time  `gorm:  "not null"`
}