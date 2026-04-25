package main

import (
	"log"
	"micr_service_auth/internal/app"
	"net/http"
)

func main() {

	application, err := app.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server started on :8080")

	log.Fatal(http.ListenAndServe(
		":8080",
		application.Server.Router(),
	))
}
