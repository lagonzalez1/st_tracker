package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tracker/app/models"
	"tracker/app/services"
)

/*
	// Decode HTTP request
    // Validate request data
    // Call authService.CreateUser()
    // Send HTTP response
*/

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Decode HTTP request
// Validate request data
// Call authService.CreateUser()
// Send HTTP response
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error found reading body")
		return
	}
	// Print the formated string body
	// fmt.Printf("Request body %s\n", string(body))

	// Parse values to json model
	var models models.RegisterRequestAdminRoot
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
	}
	// Create the user using authService
	user, err := h.authService.RegisterRootUser(models)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Printf("Error creating user: %v\n", err)
		return
	}
	// Return statment
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created successfully", user.Email)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.LoginRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.Login(models)
	if err != nil {
		http.Error(w, "failed to login with credentials", http.StatusInternalServerError)
		fmt.Printf("Error unable to login: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
