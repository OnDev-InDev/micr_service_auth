package models

import "time"

type Session struct {
	ID        string
	Username  string
	ExpiresAt time.Time
	Role      string
}
