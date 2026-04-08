package service

import (
	"context"
	"log"
	"micr_service_auth/internal/storage/models"
	"micr_service_auth/internal/storage/repository/redis"
	"time"

	"github.com/google/uuid"
)

// генерация ID
func generateSessionID() string {
	newID := uuid.New().String()
	return newID
}

// создаем сессию
func CreateSession(username string, ctx context.Context) string {
	sessionID := generateSessionID()

	//описываем новую сессию
	session := models.Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Role:      "user",
	}

	if err := redis.SetSession_InRedis(ctx, session); err != nil {
		//http.Error(w, "Internal error", 500)

	}
	return sessionID
}

func DeleteSession(ctx context.Context, value string) {
	redis.DeleteSession_FromRedis(ctx, value)
}

// смотрим куки session
func CheckSession(ctx context.Context, value string) (models.Session, bool) {
	// Получаем JSON из Redis
	sessionJSON, err := redis.GetSession_FromRedis(ctx, value)
	if err != nil {
		log.Printf("Error getting session from Redis: %v", err)
		return models.Session{}, false
	}

	if time.Now().After(sessionJSON.ExpiresAt) {
		return models.Session{}, false
	}

	return sessionJSON, true
}
