package http

import (
	"net/http"
)

type Router struct {
	router *http.ServeMux
}

func NewRouter(h *Handler, auth *AuthMiddleware) *Router {

	mux := http.NewServeMux()

	mux.HandleFunc("/post/signin", h.AuthHandler)

	mux.Handle("/protected",
		auth.RequireSession(
			auth.RequireRole("user", http.HandlerFunc(h.ProtectedHandler)),
		),
	)

	mux.HandleFunc("/post/logout", h.LogoutHandler)

	return &Router{
		router: mux,
	}
}

func (s *Router) Router() http.Handler {
	return s.router
}
