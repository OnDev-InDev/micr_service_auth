package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"micr_service_auth/internal/storage/connection_db"
	"micr_service_auth/internal/storage/models"
	"time"
)


type RedisSessionRepo struct{}




func (r *RedisSessionRepo)CreateSessionRepo(ctx context.Context, session models.Session) error {
	ttl := time.Until(session.ExpiresAt)

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = connection_db.RedisClient.Set(ctx, "session_id:"+session.ID, data, ttl).Err()
	if err != nil {
		fmt.Println("Error write", err)
		return err
	}
	fmt.Println("Success writer into Redis!")
	return nil
}

func (r *RedisSessionRepo)GetSessionRepo(ctx context.Context, sessionID string) (models.Session, error) {
	var session models.Session

	key := "session_id:" + sessionID

	value, err := connection_db.RedisClient.Get(ctx, key).Result()
	if err != nil {
		//fmt.Println("Error receiving ", err)
		return session, err
	}

	err = json.Unmarshal([]byte(value), &session)
	if err != nil {
		return session, err
	}

	//fmt.Println("Success receipt ")
	return session, nil
}

func (r *RedisSessionRepo)DeleteSessionRepo(ctx context.Context, sessionID string) error {
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
