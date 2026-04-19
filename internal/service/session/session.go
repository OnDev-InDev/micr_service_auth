package session

import (
	"context"
	"time"
  "micr_service_auth/internal/domain"
	"github.com/google/uuid"
)




type SessionRepository interface {
	CreateSessionRepo(ctx context.Context, session domain.Session) error
	GetSessionRepo(ctx context.Context, value string) (domain.Session, error)
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
func (s *SessionService) CreateSession(ctx context.Context, userID string) (string, error) {
	sessionID := generateSessionID()

	//описываем новую сессию
	session := domain.Session{
		ID:        sessionID,
		UserID:    userID,
		Role:      "user",
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := s.sessionRepo.CreateSessionRepo(ctx, session); err != nil {
    return "", err
	}
	return sessionID, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, value string) {
	s.sessionRepo.DeleteSessionRepo(ctx, value)
}

// смотрим куки session
func (s *SessionService) CheckSession(ctx context.Context, value string) (domain.Session, error) {
	// Получаем JSON из Redis
	sessionJSON, err := s.sessionRepo.GetSessionRepo(ctx, value)
	if err != nil {
		return domain.Session{}, ErrSessionNotFound
	}

	if time.Now().After(sessionJSON.ExpiresAt) {
		return domain.Session{}, ErrSessionExpired
	}

	return sessionJSON, nil
}
