package http

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(h *Handler, auth *AuthMiddleware) *Router {

	mux := http.NewServeMux()

	// auth endpoints
	mux.HandleFunc("/signin", h.AuthHandler)
	mux.HandleFunc("/logout", h.LogoutHandler)

	// protected chain
	protected := http.HandlerFunc(h.ProtectedHandler)

	mux.Handle("/protected",
		auth.RequireSession(
			auth.RequireRole("user", protected),
		),
	)

	return &Router{
		mux: mux,
	}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
