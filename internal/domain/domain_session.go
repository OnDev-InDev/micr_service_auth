package domain

import "time"

type Session struct {
	ID        string
	UserID    string
	Role      string
	ExpiresAt time.Time
}
