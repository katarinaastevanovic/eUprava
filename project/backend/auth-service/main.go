package main

import (
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/check-username", handlers.CheckUsernameHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/firebase-login", handlers.FirebaseLoginHandler).Methods("POST")
	r.HandleFunc("/complete-profile", handlers.CompleteProfileHandler).Methods("POST")
	r.Handle("/profile", middleware.JWTAuth(http.HandlerFunc(handlers.GetCurrentUser))).Methods("GET")
	r.HandleFunc("/patients/{userId}", handlers.GetPatientHandler).Methods("GET")
	r.HandleFunc("/users/doctors", handlers.GetAllDoctorsHandler).Methods("GET")
	r.HandleFunc("/users/students", handlers.GetAllStudentsHandler).Methods("GET")
	r.HandleFunc("/users/{userId}", handlers.GetUserByIdHandler).Methods("GET")

	handler := corsMiddleware(r)

	log.Println("Auth service listening on port 8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
