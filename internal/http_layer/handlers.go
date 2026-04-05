package http_layer

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
  "micr_service_auth/internal/storage/repository/redis"
	"micr_service_auth/internal/storage/repository/postgres"
	"micr_service_auth/internal/storage/models"
	"github.com/google/uuid"
)

// смотрим куки session
func checkSession(r *http.Request) (models.Session, bool) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return models.Session{}, false
	}

	// Получаем JSON из Redis
	sessionJSON, err := redis.GetSession_FromRedis(ctx, cookie.Value)
	if err != nil {
		log.Printf("Error getting session from Redis: %v", err)
		return models.Session{}, false
	}

	if time.Now().After(sessionJSON.ExpiresAt) {
		return models.Session{}, false
	}

	return sessionJSON, true
}

// создаем куки
func createSession(w http.ResponseWriter, username string, ctx context.Context) {
	sessionID := generateSessionID()

	//описываем новую сессию
	session := models.Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Role:      "user",
	}

	if err := redis.SetSession_InRedis(ctx, session); err != nil {
		http.Error(w, "Internal error", 500)
		return
	}

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

// генерация ID
func generateSessionID() string {
	newID := uuid.New().String()
	return newID
}

// один обработчик ответов
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// логинимся
func AuthHandler(w http.ResponseWriter, r *http.Request) {
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
	if !postgres.IdentificationRepo(creds.Username, creds.Password) {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}
	//создаем сессию
	createSession(w, creds.Username, ctx)

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged in",
	})
}

// ручка для авторизованных
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := checkSession(r)
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err == nil {
		redis.DeleteSession_FromRedis(ctx, cookie.Value)
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
