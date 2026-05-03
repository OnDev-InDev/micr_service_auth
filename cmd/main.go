package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OnDev-InDev/micr_service_auth/internal/app"
)

func main() {

	application, err := app.Run()
	if err != nil {
		slog.Error("failed to start app", slog.Any("error", err))
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: application.Router.Handler(),
	}

	// запуск сервера
	go func() {
		slog.Info("server started", slog.String("addr", ":8080"))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server crashed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// ловим сигналы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// HTTP shutdown
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("http shutdown error", slog.Any("error", err))
	}

	// shutdown зависимостей
	if err := application.Shutdown(); err != nil {
		slog.Error("app shutdown error", slog.Any("error", err))
	}

	slog.Info("server stopped")
}
