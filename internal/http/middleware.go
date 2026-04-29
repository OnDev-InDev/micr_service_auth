package http

import (
	"context"
	"net/http"

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
			writeError(w, errors.New(errors.CodeUnauthorized, "missing session"))
			return
		}

		session, err := m.sessionChecker.GetSession(r.Context(), cookie.Value)
		if err != nil {
			writeError(w, err)
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
			writeError(w, errors.New(errors.CodeUnauthorized, "missing role"))
			return
		}

		if userRole != role {
			writeError(w, errors.New(errors.CodeForbidden, "insufficient permissions"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
