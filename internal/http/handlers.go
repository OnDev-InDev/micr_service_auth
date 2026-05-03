package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/OnDev-InDev/micr_service_auth/internal/errors"
	"github.com/OnDev-InDev/micr_service_auth/internal/usecase"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUsecase interface {
	Login(ctx context.Context, input usecase.LoginInput) (string, error)
	Logout(ctx context.Context, sessionID string) error
}

type Handler struct {
	usecaseAuth AuthUsecase
}

func NewHandler(usecaseAuth AuthUsecase) *Handler {
	return &Handler{
		usecaseAuth: usecaseAuth,
	}
}

func (h *Handler) createCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,                              //защита от XSS
		SameSite: http.SameSiteStrictMode,           //защита от CSRF
		MaxAge:   int((30 * time.Minute).Seconds()), //30 мин
		//Secure:   true,
	}
	http.SetCookie(w, cookie)
}

// логинимся
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, errors.New(errors.CodeMethodNotAllowed, "method not allowed"))
		return
	}
	//defer r.Body.Close()
	ctx := r.Context()
	var creds Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, r, errors.New(errors.CodeValidationError, "invalid request body"))
		return
	}

	input := usecase.LoginInput{
		Email:    creds.Email,
		Password: creds.Password,
	}

	//идентификация
	sessionID, err := h.usecaseAuth.Login(ctx, input)
	if err != nil {
		writeError(w, r, err)
		return
	}

	//создаем сессию
	h.createCookie(w, sessionID)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged in",
	})
}

// ручка для авторизованных
func (h *Handler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	userID, _ := UserIDFromContext(r.Context())
	role, _ := RoleFromContext(r.Context())

	jsonResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"user_id": userID,
		"role":    role,
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, _ := r.Cookie("session_id")
	if cookie != nil {
		if err := h.usecaseAuth.Logout(ctx, cookie.Value); err != nil {
			writeError(w, r, err)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
