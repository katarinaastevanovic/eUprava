package main

import (
	"log"
	"net/http"
	"os"
	"school-service/database"
	"school-service/handlers"
	"school-service/services"

	"github.com/gorilla/mux"
)

func main() {

	db := database.Connect()

	schoolService := services.NewSchoolService(db)
	absenceHandler := handlers.NewSchoolHandler(schoolService)

	router := mux.NewRouter()
	router.HandleFunc("/students/{id}/absences", absenceHandler.GetStudentAbsences).Methods("GET")

	handler := corsMiddleware(router)

	port := "8081"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("🚀 School service listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
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
