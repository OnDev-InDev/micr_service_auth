package http_layer

import (
	"context"
	"encoding/json"
	"micr_service_auth/internal/service"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Handler struct {
	session *service.SessionService
	auth    *service.AuthService
}

func (h *Handler) createCookie(w http.ResponseWriter, username string, ctx context.Context) {
	sessionID := h.session.CreateSession(username, ctx)
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

// один обработчик ответов
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// логинимся
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Use POST"})
		return
	}
	//defer r.Body.Close()
	ctx := r.Context()
	var creds Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	if creds.Username == "" || creds.Password == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
		return
	}

	//идентификация
	if !h.auth.AuthenticateUser(creds.Username, creds.Password) {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}
	//создаем сессию
	h.createCookie(w, creds.Username, ctx)

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
		jsonResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	session, ok := h.session.CheckSession(ctx, cookie.Value)
	if !ok {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"username": session.Username,
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.session.DeleteSession(ctx, cookie.Value)
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
