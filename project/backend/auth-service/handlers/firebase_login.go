package handlers

import (
	"auth-service/database"
	"auth-service/firebase"
	"auth-service/models"
	"auth-service/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type FirebaseLoginRequest struct {
	IDToken string `json:"idToken"`
}

type FirebaseLoginResponse struct {
	Token string `json:"token"`
}

func FirebaseLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req FirebaseLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("[FirebaseLoginHandler] Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.IDToken == "" {
		log.Println("[FirebaseLoginHandler] Missing idToken")
		http.Error(w, "Missing idToken", http.StatusBadRequest)
		return
	}

	log.Println("[FirebaseLoginHandler] Initializing Firebase App")
	client := firebase.InitFirebase()
	authClient, err := client.Auth(context.Background())
	if err != nil {
		log.Printf("[FirebaseLoginHandler] Firebase auth init error: %v", err)
		http.Error(w, "Firebase auth init error", http.StatusInternalServerError)
		return
	}

	log.Println("[FirebaseLoginHandler] Verifying ID token")
	token, err := authClient.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.Printf("[FirebaseLoginHandler] Invalid Firebase token: %v", err)
		http.Error(w, "Invalid Firebase token", http.StatusUnauthorized)
		return
	}

	uid := token.UID
	email, _ := token.Claims["email"].(string)
	log.Printf("[FirebaseLoginHandler] Token verified. UID: %s, Email: %s", uid, email)

	user, err := services.FindUserByFirebaseUID(uid)
	if err == nil {
		jwtToken, genErr := services.GenerateJWT(user.ID, user.Username, string(user.Role))
		if genErr != nil {
			log.Printf("[FirebaseLoginHandler] JWT generation error: %v", genErr)
			http.Error(w, "JWT generation error", http.StatusInternalServerError)
			return
		}
		log.Printf("[FirebaseLoginHandler] User found, returning JWT for UID: %s", uid)
		json.NewEncoder(w).Encode(FirebaseLoginResponse{Token: jwtToken})
		return
	}

	log.Printf("[FirebaseLoginHandler] User not found, profile incomplete. UID: %s, Email: %s", uid, email)

	w.WriteHeader(http.StatusPreconditionRequired)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "profile incomplete",
		"email":   email,
		"uid":     uid,
	})
}

func CompleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UID      string `json:"uid"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Name     string `json:"name"`
		LastName string `json:"lastName"`
		UMCN     string `json:"umcn"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	user := models.Member{
		FirebaseUID: req.UID,
		Email:       req.Email,
		Username:    req.Username,
		Name:        req.Name,
		LastName:    req.LastName,
		UMCN:        req.UMCN,
		Role:        "user",
	}

	database.DB.Create(&user)

	jwtToken, _ := services.GenerateJWT(user.ID, user.Username, string(user.Role))
	json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}
