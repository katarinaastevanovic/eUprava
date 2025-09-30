package main

import (
	"log"
	"medical-service/database"
	"medical-service/handlers"
	"medical-service/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.Connect()

	r := mux.NewRouter()

	studentRouter := r.PathPrefix("/").Subrouter()
	studentRouter.Use(middleware.JWTAuth)
	studentRouter.Use(middleware.RBAC("STUDENT"))
	studentRouter.HandleFunc("/medical-records", handlers.CreateMedicalRecord).Methods("POST")
	studentRouter.HandleFunc("/requests", handlers.CreateRequest).Methods("POST")
	studentRouter.HandleFunc("/requests/patient/{id}", handlers.GetRequestsByPatient).Methods("GET")
	studentRouter.HandleFunc("/medical-records/{medicalRecordId}/examinations", handlers.GetExaminationsByMedicalRecord).Methods("GET")

	stDocRouter := r.PathPrefix("/").Subrouter()
	stDocRouter.Use(middleware.JWTAuth)
	stDocRouter.Use(middleware.RBAC("STUDENT", "DOCTOR"))
	stDocRouter.HandleFunc("/medical-record/{userId}", handlers.GetMedicalRecord).Methods("GET")
	stDocRouter.HandleFunc("/medical-record/{userId}", handlers.UpdateMedicalRecord).Methods("PUT")
	stDocRouter.HandleFunc("/medical-record/full/{userId}", handlers.GetFullMedicalRecord).Methods("GET")

	doctorRouter := r.PathPrefix("/").Subrouter()
	doctorRouter.Use(middleware.JWTAuth)
	doctorRouter.Use(middleware.RBAC("DOCTOR"))
	doctorRouter.HandleFunc("/requests/doctor/{id}", handlers.GetRequestsByDoctor).Methods("GET")
	doctorRouter.HandleFunc("/requests/{id}/approve", handlers.ApproveRequest).Methods("PATCH")
	doctorRouter.HandleFunc("/requests/{id}/reject", handlers.RejectRequest).Methods("PATCH")
	doctorRouter.HandleFunc("/requests/doctor/{id}/approved", handlers.GetApprovedRequestsByDoctor).Methods("GET")
	doctorRouter.HandleFunc("/examinations", handlers.CreateExamination).Methods("POST")
	doctorRouter.HandleFunc("/examinations/{requestId}", handlers.GetExaminationByRequest).Methods("GET")
	doctorRouter.HandleFunc("/requests/{requestId}/medical-record-id", handlers.GetMedicalRecordIdByRequest).Methods("GET")
	doctorRouter.HandleFunc("/requests/{requestId}", handlers.GetRequestById).Methods("GET")
	doctorRouter.HandleFunc("/certificates", handlers.CreateMedicalCertificateHandler).Methods("POST")
	doctorRouter.HandleFunc("/medical-records/{medicalRecordId}", handlers.GetFullMedicalRecordById).Methods("GET")

	serviceRouter := r.PathPrefix("/").Subrouter()
	serviceRouter.Use(middleware.JWTAuth)
	serviceRouter.HandleFunc("/patients", handlers.CreatePatientHandler).Methods("POST")
	serviceRouter.HandleFunc("/doctors", handlers.CreateDoctorHandler).Methods("POST")

	handler := corsMiddleware(r)

	log.Println("Medical service listening on port 8082")
	if err := http.ListenAndServe(":8082", handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
