package redis

import (
	"context"
	"encoding/json"
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSessionRepo struct {
	client *redis.Client
}

func NewRedisSessionRepo(client *redis.Client) *RedisSessionRepo {
	return &RedisSessionRepo{
		client: client,
	}
}

func sessionKey(id string) string {
	return "session_id:" + id
}

func (r *RedisSessionRepo) Create(ctx context.Context, s domain.Session) error {
	ttl := time.Until(s.ExpiresAt)
	key := sessionKey(s.ID)

	data, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(errors.CodeInternal, "marshal session", err)
	}

	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return errors.Wrap(errors.CodeInternal, "redis set session", err)
	}
	return nil
}

func (r *RedisSessionRepo) Get(ctx context.Context, sessionID string) (domain.Session, error) {
	var ses domain.Session

	key := sessionKey(sessionID)

	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		// ключа нет
		if err == redis.Nil {
			return ses, errors.New(errors.CodeSessionNotFound, "key not found")
		}

		// реальная ошибка Redis
		return ses, errors.Wrap(errors.ErrDB, "redis get failed", err)
	}

	// парсим JSON
	if err := json.Unmarshal([]byte(value), &ses); err != nil {
		return ses, errors.Wrap(errors.CodeInternal, "unmarshal session", err)
	}

	return ses, nil
}

func (r *RedisSessionRepo) Delete(ctx context.Context, sessionID string) error {
	key := sessionKey(sessionID)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errors.Wrap(errors.ErrDB, "delete session failed", err)
	}
	return nil
}
