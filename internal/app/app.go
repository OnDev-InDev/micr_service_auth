package app

import (
	"micr_service_auth/internal/config"
	"micr_service_auth/internal/http_layer"
	"micr_service_auth/internal/service/auth"
	"micr_service_auth/internal/service/session"
	"micr_service_auth/internal/storage/connection_db/conn_postgres"
	"micr_service_auth/internal/storage/repository/conn_redis"
	"micr_service_auth/internal/storage/repository/postgres"
	"micr_service_auth/internal/usecase"

	"github.com/jinzhu/gorm"
)

type App struct {
	DB      *gorm.DB
	Usecase *usecase.AuthUsecase
	Server  *http_layer.Server
}

func Init() (*App, error) {

	cfg := config.Load()

	// db
	db, err := conn_postgres.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	// redis
	redisClient := conn_redis.NewRedisClient(cfg)

	// repos
	userRepo := postgres.NewPostgresRepo(db)
	sessionRepo := conn_redis.NewRedisSessionRepo(redisClient)

	// services
	authService := auth.NewAuthService(userRepo)
	sessionService := session.NewSessionService(sessionRepo)

	// usecase
	uc := usecase.NewUsecase(authService, sessionService)

	// handler
	handler := http_layer.NewHandler(uc)

	// server (ВОТ ТУТ)
	server := http_layer.NewServer(handler)

	return &App{
		DB:      db,
		Usecase: uc,
		Server:  server,
	}, nil
}
