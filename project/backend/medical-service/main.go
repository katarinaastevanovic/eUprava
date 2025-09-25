package main

import (
	"log"
	"medical-service/database"
	"medical-service/handlers"
	"net/http"
)

func main() {

	database.Connect()

	mux := http.NewServeMux()
	mux.HandleFunc("/medical-records", handlers.CreateMedicalRecord)

	handler := corsMiddleware(mux)

	log.Println("Auth service listening on port 8082")
	if err := http.ListenAndServe(":8082", handler); err != nil {
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
