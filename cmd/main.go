package main

import (
	"context"
	"log"
	"micr_service_auth/internal/app"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	application, err := app.Run()
	if err != nil {
		log.Fatal(err)
	}

	// НАСТОЯЩИЙ HTTP SERVER
	server := &http.Server{
		Addr:    ":8080",
		Handler: application.Router.Router(),
	}

	// запуск сервера
	go func() {
		log.Println("Server started on :8080")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// ловим сигналы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. HTTP shutdown (реальный сервер)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	// 2. зависимости
	if err := application.Shutdown(); err != nil {
		log.Printf("App shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
