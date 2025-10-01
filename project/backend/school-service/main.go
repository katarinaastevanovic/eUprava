package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"school-service/database"
	"school-service/handlers"
	"school-service/models"
	"school-service/services"

	"github.com/gorilla/mux"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
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
	router.HandleFunc("/api/grades/student/{studentID}/subject/{subjectID}/teacher/{teacherID}", schoolHandler.GetGradesByStudentSubjectAndTeacherHandler).Methods("GET")
	router.HandleFunc("/api/grades/student/{studentID}", schoolHandler.GetAllGradesByStudentHandler).Methods("GET")
	router.HandleFunc("/students/{studentID}/subjects/{subjectID}/teachers/{teacherID}/average", schoolHandler.GetAverageByTeacherAndSubjectHandler).Methods("GET")
	router.HandleFunc("/students/{studentID}/average", schoolHandler.GetAverageByStudentHandler).Methods("GET")
	router.HandleFunc("/students/{studentID}/averages-per-subject", schoolHandler.GetAverageByStudentPerSubjectHandler).Methods("GET")
	router.HandleFunc("/api/teachers/user/{userId}", schoolHandler.GetTeacherByUserIDHandler).Methods("GET")
	router.HandleFunc("/api/classes/{classId}/students/search", schoolHandler.SearchStudentsHandler).Methods("GET")

	handler := corsMiddleware(router)

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "School service is running!")
	})

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
