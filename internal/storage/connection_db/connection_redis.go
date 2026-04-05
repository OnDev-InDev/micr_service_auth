package connection_db

import (
	"context"
	"fmt"
  "github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client


func ConnectionRedis(ctx context.Context) {
	RedisClient = redis.NewClient(&redis.Options{
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
