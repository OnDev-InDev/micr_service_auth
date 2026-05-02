package logger

import (
	"context"
	"log/slog"

	"micr_service_auth/internal/http"
)

type ContextHandler struct {
	next slog.Handler
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {

	// всегда есть
	if requestID, ok := http.RequestIDFromContext(ctx); ok {
		r.AddAttrs(slog.String("request_id", requestID))
	}

	// появляется только после auth middleware
	if userID, ok := http.UserIDFromContext(ctx); ok {
		r.AddAttrs(slog.String("user_id", userID))
	}

	return h.next.Handle(ctx, r)
}

func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{next: h.next.WithAttrs(attrs)}
}

func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{next: h.next.WithGroup(name)}
}
