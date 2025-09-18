package main

import (
	"log"
	"net/http"

	"auth-service/database"
	"auth-service/handlers"
)

func main() {
	database.Connect()

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/check-username", handlers.CheckUsernameHandler)

	//mux.HandleFunc("/login", handlers.LoginHandler)

	log.Println("🚀 Auth service listening on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
