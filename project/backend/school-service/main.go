package main

import (
	"log"
	"net/http"
	"os"
	"school-service/database"
	"school-service/handlers"
	"school-service/models"
	"school-service/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Nije pronađen .env fajl, koristim sistemske varijable")
	}

	db := database.Connect()

	if err := db.AutoMigrate(&models.Class{}, &models.Teacher{}); err != nil {
		log.Fatal(err)
	}
	log.Println("Migracija uspešna!")

	schoolService := services.NewSchoolService(db)
	schoolHandler := handlers.NewSchoolHandler(schoolService)

	router := mux.NewRouter()
	router.HandleFunc("/students/{id}/absences", schoolHandler.GetStudentAbsences).Methods("GET")
	router.HandleFunc("/absences/{id}/type", schoolHandler.UpdateAbsenceType).Methods("PUT")
	router.HandleFunc("/absences", schoolHandler.CreateAbsences).Methods("POST", "OPTIONS")
	router.HandleFunc("/teachers/{id}/classes", schoolHandler.GetClassesForTeacher).Methods("GET")
	router.HandleFunc("/api/classes/{classID}/students", schoolHandler.GetStudentsByClass).Methods("GET")
	router.HandleFunc("/students/{studentID}/subjects/{subjectID}/absences/count", schoolHandler.GetAbsenceCountForSubject).Methods("GET")
	router.HandleFunc("/students/by-user/{userId}", schoolHandler.GetStudentByUserID).Methods("GET")
	router.HandleFunc("/students/by-user/{userId}/profile", schoolHandler.GetStudentFullProfile).Methods("GET")

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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
