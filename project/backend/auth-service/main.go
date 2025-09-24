package main

import (
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middleware"
	"log"
	"net/http"
)

func main() {

	database.Connect()

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/check-username", handlers.CheckUsernameHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/firebase-login", handlers.FirebaseLoginHandler)
	mux.HandleFunc("/complete-profile", handlers.CompleteProfileHandler)
	mux.Handle("/profile", middleware.JWTAuth(http.HandlerFunc(handlers.GetCurrentUser)))

	handler := corsMiddleware(mux)

	log.Println("Auth service listening on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
