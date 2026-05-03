package app

import (
	"github.com/OnDev-InDev/micr_service_auth/internal/config"
	"github.com/OnDev-InDev/micr_service_auth/internal/http"
	"github.com/OnDev-InDev/micr_service_auth/internal/logger"
	"github.com/OnDev-InDev/micr_service_auth/internal/storage/repository/postgres"
	"github.com/OnDev-InDev/micr_service_auth/internal/storage/repository/redis"
	"github.com/OnDev-InDev/micr_service_auth/internal/usecase"
	"github.com/jinzhu/gorm"
	redisClient "github.com/redis/go-redis/v9"
)

type App struct {
	DB          *gorm.DB
	RedisClient *redisClient.Client
	Router      *http.Router
}

func Run() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger.Init()

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	rdb, err := redis.NewRedisClient()
	if err != nil {
		return nil, err
	}

	userRepo := postgres.NewPostgresRepo(db)
	sessionRepo := redis.NewRedisSessionRepo(rdb)

	//authService := service.NewAuthService(userRepo)
	//sessionService := service.NewSessionService(sessionRepo)

	uc := usecase.NewUsecase(userRepo, sessionRepo)

	handler := http.NewHandler(uc)
	authMiddleware := http.NewAuthMiddleware(uc)

	router := http.NewRouter(handler, authMiddleware)

	return &App{
		DB:          db,
		RedisClient: rdb,
		Router:      router,
	}, nil
}

func (a *App) Shutdown() error {
	if sqlDB := a.DB.DB(); sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	return a.RedisClient.Close()
}
