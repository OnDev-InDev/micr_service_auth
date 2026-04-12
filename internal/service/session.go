package service

import (
	"context"
	"log"
	"micr_service_auth/internal/storage/models"
	"time"

	"github.com/google/uuid"
)


type SessionService struct {
	sessionRepo SessionRepository
}

// генерация ID
func generateSessionID() string {
	newID := uuid.New().String()
	return newID
}

// создаем сессию
func (s *SessionService)CreateSession(username string, ctx context.Context) string {
	sessionID := generateSessionID()

	//описываем новую сессию
	session := models.Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Role:      "user",
	}
	if err := s.sessionRepo.CreateSessionRepo(session, ctx); err != nil {
		//http.Error(w, "Internal error", 500)

	}
	return sessionID
}

func (s *SessionService)DeleteSession(ctx context.Context, value string) {
	s.sessionRepo.DeleteSessionRepo(ctx, value)
}

// смотрим куки session
func (s *SessionService)CheckSession(ctx context.Context, value string) (models.Session, bool) {
	// Получаем JSON из Redis
	sessionJSON, err := s.sessionRepo.GetSessionRepo(ctx, value)
	if err != nil {
		log.Printf("Error getting session from Redis: %v", err)
		return models.Session{}, false
	}

	if time.Now().After(sessionJSON.ExpiresAt) {
		return models.Session{}, false
	}

	return sessionJSON, true
}
