package http_layer

import (
	"net/http"
)

type Server struct {
	router *http.ServeMux
}

func NewServer(h *Handler) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/post/signin", h.AuthHandler)
	mux.HandleFunc("/protected", h.ProtectedHandler)
	mux.HandleFunc("/post/logout", h.LogoutHandler)

	return &Server{router: mux}
}

func (s *Server) Router() http.Handler {
	return s.router
}
