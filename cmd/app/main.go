package main

import (
	"context"
	"log"
	"micr_service_auth/internal/http_layer"
	"micr_service_auth/internal/storage/connection_db"
	"net/http"
)


func main() {

	connection_db.ConnectionRedis(context.Background())
	http_layer.Router()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
