package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"micr_service_auth/internal/storage/connection_db"
	"micr_service_auth/internal/storage/models"
	"time"

	"github.com/redis/go-redis/v9"
)





type RedisSessionRepo struct{}

func (r *RedisSessionRepo) CreateSessionRepo(ctx context.Context, session models.Session) error {
	ttl := time.Until(session.ExpiresAt)

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	key := "session_id:" + session.ID

	err = connection_db.RedisClient.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("repo Write to Redis:", err)
	}
	return nil
}

func (r *RedisSessionRepo) GetSessionRepo(ctx context.Context, sessionID string) (models.Session, error) {
	var session models.Session

	if sessionID == "" {
		return session, fmt.Errorf("empty session id")
	}

	key := "session_id:" + sessionID

	value, err := connection_db.RedisClient.Get(ctx, key).Result()
	if err != nil {
		// ключа нет — это нормальная ситуация
		if errors.Is(err, redis.Nil) {
			return session, ErrSessionNotFound
		}

		// реальная ошибка Redis
		return session, fmt.Errorf("redis get failed for key %s: %w", key, err)
	}

	// парсим JSON
	if err := json.Unmarshal([]byte(value), &session); err != nil {
		return session, fmt.Errorf("failed to unmarshal session (key %s): %w", key, err)
	}

	return session, nil
}

func (r *RedisSessionRepo) DeleteSessionRepo(ctx context.Context, sessionID string) error {
	key := "session_id:" + sessionID
	deleteCount, err := connection_db.RedisClient.Del(ctx, key).Result()
	if err != nil {
		//fmt.Println("Error delete", err)
		return err
	}

	// значит у сессии истек срок хранения
	if deleteCount == 0 {
		//fmt.Println("Session not found")
		return err
	}

	//fmt.Println("Success delete", deleteCount)
	return nil
}
