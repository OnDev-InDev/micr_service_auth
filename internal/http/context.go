package http

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	RoleKey      contextKey = "role"
	RequestIDKey contextKey = "request_id"
)
