package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/service/session"
	"micr_service_auth/internal/storage/connection_db"
	"micr_service_auth/internal/storage/models"
	"time"

	"github.com/redis/go-redis/v9"
)





type RedisSessionRepo struct{}



func sessionKey(id string) string {
	return "session_id:" + id
}








func (r *RedisSessionRepo) CreateSessionRepo(ctx context.Context, s domain.Session) error {
	ttl := time.Until(s.ExpiresAt)
	key := "session_id:" + s.ID
	
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	err = connection_db.RedisClient.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set session: %w", err)
	}
	return nil
}












func (r *RedisSessionRepo) GetSessionRepo(ctx context.Context, sessionID string) (domain.Session, error) {
  var ses domain.Session
	
	if sessionID == "" {
		return ses, fmt.Errorf("empty session id")
	}

	key := "session_id:" + sessionID

	value, err := connection_db.RedisClient.Get(ctx, key).Result()
	if err != nil {
		// ключа нет — это нормальная ситуация
		if errors.Is(err, redis.Nil) {
			return ses, session.ErrSessionNotFound
		}

		// реальная ошибка Redis
		return ses, fmt.Errorf("redis get failed for key %s: %w", key, err)
	}

	// парсим JSON
	if err := json.Unmarshal([]byte(value), &ses); err != nil {
		return ses, fmt.Errorf("failed to unmarshal session (key %s): %w", key, err)
	}

	return ses, nil
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
