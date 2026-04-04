package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	repository "service_auth/repository"
	"time"

	"github.com/google/uuid"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// смотрим куки session
func checkSession(r *http.Request) (repository.Session, bool) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return repository.Session{}, false
	}

	// Получаем JSON из Redis
	sessionJSON, err := repository.GetSession_FromRedis(ctx, cookie.Value)
	if err != nil {
		log.Printf("Error getting session from Redis: %v", err)
		return repository.Session{}, false
	}

	if time.Now().After(sessionJSON.ExpiresAt) {
		return repository.Session{}, false
	}

	return sessionJSON, true
}

// создаем куки
func createSession(w http.ResponseWriter, username string, ctx context.Context) {
	sessionID := generateSessionID()

	//описываем новую сессию
	session := repository.Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Role:      "user",
	}

	if err := repository.SetSession_InRedis(ctx, session); err != nil {
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
func authHandler(w http.ResponseWriter, r *http.Request) {
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
	if !repository.IdentificationRepo(creds.Username, creds.Password) {
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
func protectedHandler(w http.ResponseWriter, r *http.Request) {
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err == nil {
		repository.DeleteSession_FromRedis(ctx, cookie.Value)
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

func main() {
	repository.ConnectionRedis(context.Background())

	http.HandleFunc("/post/signin", authHandler)
	http.HandleFunc("/protected", protectedHandler)
	http.HandleFunc("/post/logout", logoutHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
