package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	ID        string
	Username  string
	ExpiresAt time.Time
	Role      string
}

var RedisClient *redis.Client

func SetSession_InRedis(ctx context.Context, session Session) error {
	ttl := time.Until(session.ExpiresAt)

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = RedisClient.Set(ctx, "session_id:"+session.ID, data, ttl).Err()
	if err != nil {
		fmt.Println("Error write", err)
		return err
	}
	fmt.Println("Success writer into Redis!")
	return nil
}

func GetSession_FromRedis(ctx context.Context, sessionID string) (Session, error) {
	var session Session

	key := "session_id:" + sessionID

	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error receiving ", err)
		return session, err
	}

	err = json.Unmarshal([]byte(value), &session)
	if err != nil {
		return session, err
	}

	fmt.Println("Success receipt ")
	return session, nil
}

func DeleteSession_FromRedis(ctx context.Context, sessionID string) error {
	key := "session_id:" + sessionID
	deleteCount, err := RedisClient.Del(ctx, key).Result()
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
