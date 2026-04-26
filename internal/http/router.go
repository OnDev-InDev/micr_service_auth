package http

import (
	"net/http"
)

type Server struct {
	router *http.ServeMux
}

func NewServer(h *Handler, auth *AuthMiddleware) *Server {

	mux := http.NewServeMux()

	mux.HandleFunc("/post/signin", h.AuthHandler)

	mux.Handle("/protected",
	  auth.RequireAuth(
		  auth.RequireRole("user", http.HandlerFunc(h.ProtectedHandler)),
	),
)

	mux.HandleFunc("/post/logout", h.LogoutHandler)

	return &Server{
		router: mux,
	}
}

func (s *Server) Router() http.Handler {
	return s.router
}
