package app

import (
	"micr_service_auth/internal/config"
	"micr_service_auth/internal/http"
	"micr_service_auth/internal/service/auth"
	"micr_service_auth/internal/service/session"
	"micr_service_auth/internal/storage/connection"
	"micr_service_auth/internal/storage/repository/postgres"
	"micr_service_auth/internal/storage/repository/redis"
	"micr_service_auth/internal/usecase"

	"github.com/jinzhu/gorm"
)

type App struct {
	DB      *gorm.DB
	Usecase *usecase.AuthUsecase
	Server  *http.Server
}

func Init() (*App, error) {

	cfg := config.Load()

	db, err := connection.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := connection.NewRedisClient()
	if err != nil {
		return nil, err
	}

	userRepo := postgres.NewPostgresRepo(db)
	sessionRepo := redis.NewRedisSessionRepo(redisClient)

	authService := auth.NewAuthService(userRepo)
	sessionService := session.NewSessionService(sessionRepo)

	uc := usecase.NewUsecase(authService, sessionService)

	handler := http.NewHandler(uc)

	authMiddleware := http.NewAuthMiddleware(uc)

	server := http.NewServer(handler, authMiddleware)

	return &App{
		DB:      db,
		Usecase: uc,
		Server:  server,
	}, nil
}
