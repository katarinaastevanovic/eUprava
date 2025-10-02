package handlers

import (
	"encoding/json"
	"net/http"
	"school-service/services"
	"time"
)

type GradeHandler struct {
	gradeService *services.GradeService
}

func NewGradeHandler(gs *services.GradeService) *GradeHandler {
	return &GradeHandler{gradeService: gs}
}

type CreateGradeRequest struct {
	Value     int  `json:"value"`
	StudentID uint `json:"student_id"`
	SubjectID uint `json:"subject_id"`
	TeacherID uint `json:"teacher_id"`
}

type GradeResponse struct {
	ID        uint      `json:"id"`
	Value     int       `json:"value"`
	Date      time.Time `json:"date"`
	StudentID uint      `json:"student_id"`
	SubjectID uint      `json:"subject_id"`
	TeacherID uint      `json:"teacher_id"`
}

func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
	var req CreateGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	grade, err := h.gradeService.CreateGrade(req.Value, req.StudentID, req.SubjectID, req.TeacherID)
	if err != nil {
		http.Error(w, "Failed to create grade: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GradeResponse{
		ID:        grade.ID,
		Value:     grade.Value,
		Date:      grade.Date,
		StudentID: grade.StudentID,
		SubjectID: grade.SubjectID,
		TeacherID: grade.TeacherID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
