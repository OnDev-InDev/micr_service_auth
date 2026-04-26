package http

import (
	"context"
	"net/http"

	"micr_service_auth/internal/errors"
	"micr_service_auth/internal/usecase"
)

type AuthMiddleware struct {
	authUC *usecase.AuthUsecase
}

func NewAuthMiddleware(authUC *usecase.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		authUC: authUC,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_id")
		if err != nil {
			writeError(w, errors.New(errors.CodeInvalidCredentials, "unauthorized"))
			return
		}

		session, err := m.authUC.GetSession(r.Context(), cookie.Value)
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
			writeError(w, errors.New(errors.CodeInvalidCredentials, "unauthorized"))
			return
		}

		if userRole != role {
			writeError(w, errors.New(errors.CodeInvalidCredentials, "forbidden"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
