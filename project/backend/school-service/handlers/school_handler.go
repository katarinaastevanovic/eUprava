package handlers

import (
	"encoding/json"
	"net/http"
	"school-service/services"
	"strconv"

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
