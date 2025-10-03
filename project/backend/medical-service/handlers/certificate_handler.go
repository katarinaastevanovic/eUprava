package handlers

import (
	"encoding/json"
	"medical-service/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CertificateHandler struct {
	Service *services.CertificateService
}

func NewPatientHandler(service *services.CertificateService) *CertificateHandler {
	return &CertificateHandler{Service: service}
}

func (h *CertificateHandler) HasCertificateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr, ok := vars["userId"]
	if !ok {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	hasCert, err := h.Service.HasCertificate(uint(userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"userId":         userId,
		"hasCertificate": hasCert,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
