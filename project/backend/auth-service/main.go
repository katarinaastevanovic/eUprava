package main

import (
	"auth-service/database"
	"auth-service/handlers"
	"log"
	"net/http"
)

func main() {

	database.Connect()

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/api/check-username", handlers.CheckUsernameHandler)

	handler := corsMiddleware(http.DefaultServeMux)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
