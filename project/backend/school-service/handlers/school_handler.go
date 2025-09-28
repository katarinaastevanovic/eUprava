package handlers

import (
	"encoding/json"
	"net/http"
	"school-service/models"
	"school-service/services"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type SchoolHandler struct {
	Service *services.SchoolService
}

func NewSchoolHandler(s *services.SchoolService) *SchoolHandler {
	return &SchoolHandler{Service: s}
}

type AbsenceResponse struct {
	ID      uint   `json:"id"`
	Type    string `json:"type"`
	Date    string `json:"date"`
	Subject string `json:"subject"`
}

type CreateAbsenceRequest struct {
	Date      time.Time `json:"date"`
	StudentID uint      `json:"studentId"`
	SubjectID uint      `json:"subjectId"`
}

type AbsenceResponse2 struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Date      time.Time `json:"date"`
	StudentID uint      `json:"student_id"`
	SubjectID uint      `json:"subject_id"`
}

type ClassDTO struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

type TeacherClassesResponse struct {
	SubjectName string     `json:"subject_name"`
	Classes     []ClassDTO `json:"classes"`
}

func (h *SchoolHandler) GetStudentAbsences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	absences, count, err := h.Service.GetAbsencesByUserID(uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseAbsences []AbsenceResponse
	for _, a := range absences {
		responseAbsences = append(responseAbsences, AbsenceResponse{
			ID:      a.ID,
			Type:    string(a.Type),
			Date:    a.Date.Format("2006-01-02 15:04"),
			Subject: a.Subject.Name,
		})
	}

	resp := map[string]interface{}{
		"user_id":  userID,
		"count":    count,
		"absences": responseAbsences,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

type UpdateAbsenceRequest struct {
	Type models.AbsenceType `json:"type"`
}

func (h *SchoolHandler) UpdateAbsenceType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	absenceID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid absence ID", http.StatusBadRequest)
		return
	}

	var req UpdateAbsenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Type != models.Excused && req.Type != models.Unexcused && req.Type != models.Pending {
		http.Error(w, "Invalid absence type", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateAbsenceType(uint(absenceID), req.Type); err != nil {
		http.Error(w, "Failed to update absence type", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Absence type updated successfully",
	})
}

func (h *SchoolHandler) CreateAbsence(w http.ResponseWriter, r *http.Request) {
	var req CreateAbsenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	absence := models.Absence{
		Type:      models.Pending,
		Date:      req.Date,
		StudentID: req.StudentID,
		SubjectID: req.SubjectID,
	}

	if absence.Date.IsZero() {
		absence.Date = time.Now()
	}

	if err := h.Service.CreateAbsence(&absence); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AbsenceResponse2{
		ID:        absence.ID,
		Type:      string(absence.Type),
		Date:      absence.Date,
		StudentID: absence.StudentID,
		SubjectID: absence.SubjectID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *SchoolHandler) GetClassesForTeacher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	teacherID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher id", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher
	if err := h.Service.DB.First(&teacher, teacherID).Error; err != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	var subject models.Subject
	if err := h.Service.DB.First(&subject, teacher.SubjectID).Error; err != nil {
		http.Error(w, "Subject not found", http.StatusNotFound)
		return
	}

	classes, err := h.Service.GetClassesByTeacherID(uint(teacherID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dtos := make([]ClassDTO, 0, len(classes))
	for _, c := range classes {
		dtos = append(dtos, ClassDTO{
			ID:    c.ID,
			Title: c.Title,
			Year:  c.Year,
		})
	}

	response := TeacherClassesResponse{
		SubjectName: subject.Name,
		Classes:     dtos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SchoolHandler) GetStudentsByClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classIDStr := vars["classID"]

	classID, err := strconv.Atoi(classIDStr)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	students, err := h.Service.GetStudentsByClassID(uint(classID))
	if err != nil {
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
