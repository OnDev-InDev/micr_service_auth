package connection_db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	Client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := Client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	return Client, nil
}
