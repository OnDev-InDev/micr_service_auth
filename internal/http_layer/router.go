package http_layer

import (
	"net/http"
)

func Router() {
	http.HandleFunc("/post/signin", AuthHandler)
	http.HandleFunc("/protected", ProtectedHandler)
	http.HandleFunc("/post/logout", LogoutHandler)
}
