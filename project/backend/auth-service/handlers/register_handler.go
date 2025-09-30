package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"auth-service/services"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := services.RegisterUser(req)
	if err != nil {
		switch err {
		case services.ErrUsernameExists:
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		case services.ErrInvalidUsername:
			http.Error(w, "Username must be at least 4 characters long", http.StatusBadRequest)
			return
		case services.ErrUMCNExists:
			http.Error(w, "UMCN already exists", http.StatusConflict)
			return
		case services.ErrEmailExists:
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		case services.ErrInvalidUMCN:
			http.Error(w, "UMCN must be exactly 13 characters", http.StatusBadRequest)
			return
		case services.ErrInvalidPassword:
			http.Error(w, "Password must be at least 8 characters long, contain one uppercase, one lowercase, and one special character", http.StatusBadRequest)
			return
		default:
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
	}

	jwtToken, err := services.GenerateJWT(user.ID, user.Username, string(user.Role))
	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}

	switch user.Role {
	case "STUDENT":
		patientBody, _ := json.Marshal(map[string]interface{}{
			"userId":   user.ID,
			"name":     user.Name,
			"lastName": user.LastName,
		})
		patientReq, _ := http.NewRequest("POST", "http://medical-service:8082/patients", bytes.NewBuffer(patientBody))
		patientReq.Header.Set("Content-Type", "application/json")
		patientReq.Header.Set("Authorization", "Bearer "+jwtToken)

		resp, err := client.Do(patientReq)
		if err != nil || resp.StatusCode != http.StatusCreated {
			http.Error(w, "User registered but failed to create patient", http.StatusAccepted)
			return
		}
		defer resp.Body.Close()

		var createdPatient struct {
			ID uint `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&createdPatient); err != nil {
			http.Error(w, "Failed to decode patient response", http.StatusInternalServerError)
			return
		}

		recordBody, _ := json.Marshal(map[string]interface{}{
			"patientId": createdPatient.ID,
		})
		recordReq, _ := http.NewRequest("POST", "http://medical-service:8082/medical-records", bytes.NewBuffer(recordBody))
		recordReq.Header.Set("Content-Type", "application/json")
		recordReq.Header.Set("Authorization", "Bearer "+jwtToken)

		resp, err = client.Do(recordReq)
		if err != nil || resp.StatusCode != http.StatusCreated {
			http.Error(w, "User registered but failed to create medical record", http.StatusAccepted)
			return
		}
		defer resp.Body.Close()

	case "DOCTOR":
		doctorBody, _ := json.Marshal(map[string]interface{}{
			"userId":   user.ID,
			"name":     user.Name,
			"lastName": user.LastName,
		})
		doctorReq, _ := http.NewRequest("POST", "http://medical-service:8082/doctors", bytes.NewBuffer(doctorBody))
		doctorReq.Header.Set("Content-Type", "application/json")
		doctorReq.Header.Set("Authorization", "Bearer "+jwtToken)

		resp, err := client.Do(doctorReq)
		if err != nil || resp.StatusCode != http.StatusCreated {
			http.Error(w, "User registered but failed to create doctor", http.StatusAccepted)
			return
		}
		defer resp.Body.Close()
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"token":   jwtToken,
	})
}

func CheckUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	available, err := services.CheckUsernameAvailable(username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !available {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
