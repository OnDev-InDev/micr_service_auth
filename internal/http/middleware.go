package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

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

		ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, RoleKey, session.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userRole, ok := r.Context().Value(RoleKey).(string)
		if !ok || userRole == "" {
			writeError(w, r, errors.New(errors.CodeUnauthorized, "missing role"))
			return
		}

		if userRole != role {
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

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
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

		duration := time.Since(start)

		requestID, _ := r.Context().Value(RequestIDKey).(string)

		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%s",
			requestID,
			r.Method,
			r.URL.Path,
			lrw.status,
			duration,
		)
	})
}
