package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var RedisClient *redis.Client

func SetSession_InRedis(ctx context.Context, sessionID string, username string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	err := RedisClient.Set(ctx, sessionID, username, ttl).Err()
	if err != nil {
		fmt.Println("Error write", err)
	}
	fmt.Println("Success writer into Redis!")
	return nil
}

func GetSession_FromRedis(ctx context.Context, sessionID string) (string, error) {
	value, err := RedisClient.Get(ctx, sessionID).Result()
	if err != nil {
		fmt.Println("Error receiving ", err)
	}

	fmt.Println("Success receipt ")
	return value, nil
}

func DeleteSession_FromRedis(ctx context.Context, sessionID string) error {
	deleteCount, err := RedisClient.Del(ctx, sessionID).Result()
	if err != nil {
		fmt.Println("Error delete", err)
		return err
	}

	// значит у сессии истек срок хранения
	if deleteCount == 0 {
		fmt.Println("Session not found")
	}

	fmt.Println("Success delete", deleteCount)
	return nil
}

func ConnectionRedis(ctx context.Context) {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connection Redis", err)
		return
	}

	fmt.Println("Connection Redis OK!")
}
