package service

import (
	"context"
	"micr_service_auth/internal/storage/models"
	"time"

	"github.com/google/uuid"
)




type SessionRepository interface {
	CreateSessionRepo(ctx context.Context, session models.Session) error
	GetSessionRepo(ctx context.Context, value string) (models.Session, error)
	DeleteSessionRepo(ctx context.Context, value string)
}



type SessionService struct {
	sessionRepo SessionRepository
}

// генерация ID
func generateSessionID() string {
	newID := uuid.New().String()
	return newID
}

// создаем сессию
func (s *SessionService) CreateSession(ctx context.Context, userID string) string {
	sessionID := generateSessionID()

	//описываем новую сессию
	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Role:      "user",
	}

	if err := s.sessionRepo.CreateSessionRepo(ctx, session); err != nil {

	}
	return sessionID
}

func (s *SessionService) DeleteSession(ctx context.Context, value string) {
	s.sessionRepo.DeleteSessionRepo(ctx, value)
}

// смотрим куки session
func (s *SessionService) CheckSession(ctx context.Context, value string) (models.Session, error) {
	// Получаем JSON из Redis
	sessionJSON, err := s.sessionRepo.GetSessionRepo(ctx, value)
	if err != nil {
		return models.Session{}, err
	}

	if time.Now().After(sessionJSON.ExpiresAt) {
		return models.Session{}, err
	}

	return sessionJSON, nil
}
