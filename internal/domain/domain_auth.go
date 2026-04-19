package domain

import "time"


type Session struct {
	ID         string
	UserID     string
	Role       string
	ExpiresAt  time.Time
}



type User struct {
	ID             string
	Email          string
	PasswordHash   string
	Role           string
}