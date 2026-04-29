package app

import (
	"micr_service_auth/internal/config"
	"micr_service_auth/internal/http"
	"micr_service_auth/internal/storage/connection"
	"micr_service_auth/internal/storage/repository/postgres"
	"micr_service_auth/internal/storage/repository/redis"
	"micr_service_auth/internal/usecase"

	"github.com/jinzhu/gorm"
	redisClient "github.com/redis/go-redis/v9"
)

type App struct {
	DB          *gorm.DB
	RedisClient *redisClient.Client
	Router      *http.Router
}

func Run() (*App, error) {
	cfg := config.Load()

	// --- Postgres ---
	db, err := connection.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	// --- Redis ---
	rdb, err := connection.NewRedisClient()
	if err != nil {
		return nil, err
	}

	// --- Repos ---
	userRepo := postgres.NewPostgresRepo(db)
	sessionRepo := redis.NewRedisSessionRepo(rdb)

	// --- Services ---
	//authService := service.NewAuthService(userRepo)
	//sessionService := service.NewSessionService(sessionRepo)

	// --- Usecase ---
	uc := usecase.NewUsecase(userRepo, sessionRepo)

	// --- HTTP layer ---
	handler := http.NewHandler(uc)
	authMiddleware := http.NewAuthMiddleware(uc)

	router := http.NewRouter(handler, authMiddleware)

	return &App{
		DB:          db,
		RedisClient: rdb,
		Router:      router,
	}, nil
}

// Shutdown 
func (a *App) Shutdown() error {
	if sqlDB := a.DB.DB(); sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	return a.RedisClient.Close()
}
