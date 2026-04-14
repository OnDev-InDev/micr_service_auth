package service

import (
	"context"
	"micr_service_auth/internal/storage/models"
)

type SessionRepository interface {
	CreateSessionRepo(ctx context.Context, session models.Session) error
	GetSessionRepo(ctx context.Context, value string) (models.Session, error)
	DeleteSessionRepo(ctx context.Context, value string)
}
