package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/OnDev-InDev/micr_service_auth/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.ConfigRedis) (*redis.Client, error) {
	redisURL := fmt.Sprintf(
		"redis://%s:%s",
		cfg.RedisHost,
		cfg.RedisPort,
	)

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	opt.DialTimeout = 3 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return client, nil
}
