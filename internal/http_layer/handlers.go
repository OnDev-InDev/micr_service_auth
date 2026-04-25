package http_layer

import (
	"context"
	"encoding/json"
	"micr_service_auth/internal/domain"
	"micr_service_auth/internal/errors"
	"micr_service_auth/internal/usecase"
	"net/http"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUsecase interface {
	Login(ctx context.Context, input usecase.LoginInput) (string, error)
	GetSession(ctx context.Context, sessionID string) (domain.Session, error)
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

func (h *Handler) createCookie(ctx context.Context, w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,                    //защита от XSS
		SameSite: http.SameSiteStrictMode, //защита от CSRF
		MaxAge:   1800,                    //30 мин
		//Secure:   true,
	}
	http.SetCookie(w, cookie)
}

// логинимся
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, errors.New(errors.CodeValidationError, "user Post"))
		return
	}
	//defer r.Body.Close()
	ctx := r.Context()
	var creds Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, errors.New(errors.CodeValidationError, "invalid JSON body"))
		return
	}

	input := usecase.LoginInput{
		Email:    creds.Email,
		Password: creds.Password,
	}

	//идентификация
	sessionID, err := h.usecaseAuth.Login(ctx, input)
	if err != nil {
		writeError(w, err)
		return
	}

	//создаем сессию
	h.createCookie(ctx, w, sessionID)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged in",
	})
}

// ручка для авторизованных
func (h *Handler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Извлекаем куку "session_id" из запроса
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Если куки нет в запросе, считаем пользователя неавторизованным
		writeError(w, errors.New(errors.CodeInvalidCredentials, "unauthorized"))
		return
	}

	ses, err := h.usecaseAuth.GetSession(ctx, cookie.Value)
	if err != nil {
		writeError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"role":    ses.Role,
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, _ := r.Cookie("session_id")
	if cookie != nil {
		if err := h.usecaseAuth.Logout(ctx, cookie.Value); err != nil {
			writeError(w, err)
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
