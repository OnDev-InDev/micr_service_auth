package logger

import (
	"log/slog"
	"os"
)

func Init() {
	base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(&ContextHandler{next: base})
	slog.SetDefault(logger)
}
