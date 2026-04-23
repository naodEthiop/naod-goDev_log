package main

import (
	"fmt"
	"net/http"

	"github.com/naodEthiop/naod-goDev_log.git/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", handlers.CreateUser)
	mux.HandleFunc("GET /users", handlers.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", handlers.GetUser)
	mux.HandleFunc("DELETE /users/{id}", handlers.DeleteUser)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}