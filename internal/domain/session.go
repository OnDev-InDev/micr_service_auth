package domain

import "time"

type Session struct {
	ID        string
	UserID    string
	Role      string
	CreatedAt time.Time
	ExpiresAt time.Time
}
