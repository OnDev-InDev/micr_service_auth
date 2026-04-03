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

type Session struct {
	ID        string
	Username  string
	ExpiresAt time.Time
}

var ctx = context.Background()

// смотрим куки session
func checkSession(r *http.Request) (Session, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return Session{}, false
	}

	// Получаем JSON строку из Redis
	sessionJSON, err := repository.GetSession_FromRedis(ctx, cookie.Value)
	if err != nil {
		log.Printf("Error getting session from Redis: %v", err)
		return Session{}, false
	}

	// Декодируем JSON в структуру Session
	var session Session
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		log.Printf("Error unmarshaling session: %v", err)
		return Session{}, false
	}
	return session, true
}

// создаем куки
func createSession(w http.ResponseWriter, username string) {
	sessionID := generateSessionID()
	//описываем новую сессию
	session := Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	// // в мапу добавляем сессию созданную
	// sessions[sessionID] = session
	//полагаю должен быть ID постгри, как ссылка на пользователя. а роль в редис тоже хранить ?
	repository.SetSession_InRedis(ctx, sessionID, username, session.ExpiresAt)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,                     //защита от XSS
		SameSite: http.SameSiteDefaultMode, //защита от CSRF
		MaxAge:   1800,                     //30 мин
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
	defer r.Body.Close()

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
	createSession(w, creds.Username)

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
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	if cookie, err := r.Cookie("session_id"); err == nil {
		repository.DeleteSession_FromRedis(ctx, cookie.Value)
	}

	jsonResponse(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})

}

func main() {
	http.HandleFunc("/post/signin", authHandler)
	http.HandleFunc("/protected", protectedHandler)
	http.HandleFunc("/post/logout", logoutHandler)

	repository.ConnectionRedis(ctx)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
