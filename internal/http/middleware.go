package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"log/slog"

	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/errors"
)

type SessionChecker interface {
	GetSession(ctx context.Context, id string) (domain.Session, error)
}

type AuthMiddleware struct {
	sessionChecker SessionChecker
}

func NewAuthMiddleware(sessionChecker SessionChecker) *AuthMiddleware {
	return &AuthMiddleware{
		sessionChecker: sessionChecker,
	}
}

func (m *AuthMiddleware) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_id")
		if err != nil {
			writeError(w, r, errors.New(errors.CodeUnauthorized, "missing session"))
			return
		}

		session, err := m.sessionChecker.GetSession(r.Context(), cookie.Value)
		if err != nil {
			writeError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		ctx = context.WithValue(ctx, roleKey, session.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger := LoggerFromContext(r.Context())

		userRole, ok := RoleFromContext(r.Context())
		if !ok || userRole == "" {
			logger.Warn("missing role")
			writeError(w, r, errors.New(errors.CodeUnauthorized, "missing role"))
			return
		}

		if userRole != role {
			logger.Warn("insufficient permissions",
				slog.String("required_role", role),
				slog.String("user_role", userRole),
			)
			writeError(w, r, errors.New(errors.CodeForbidden, "insufficient permissions"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := generateID()

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		lrw := &responseWriter{
			ResponseWriter: w,
			status:         200,
		}

		next.ServeHTTP(lrw, r)

		slog.Info("http_request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", lrw.status),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
