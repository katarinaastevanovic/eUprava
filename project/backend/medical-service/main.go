package main

import (
	"log"
	"medical-service/database"
	"medical-service/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	database.Connect()

	r := mux.NewRouter()
	r.HandleFunc("/medical-records", handlers.CreateMedicalRecord).Methods("POST")
	r.HandleFunc("/medical-record/{userId}", handlers.GetMedicalRecord).Methods("GET")
	r.HandleFunc("/medical-record/{userId}", handlers.UpdateMedicalRecord).Methods("PUT")
	r.HandleFunc("/medical-record/full/{userId}", handlers.GetFullMedicalRecord).Methods("GET")

	handler := corsMiddleware(r)

	log.Println("Medical service listening on port 8082")
	if err := http.ListenAndServe(":8082", handler); err != nil {
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
