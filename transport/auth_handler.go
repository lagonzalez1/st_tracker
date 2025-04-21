package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tracker/app/models"
	"tracker/app/services"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Example response
}

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
	user, err := h.authService.LoginAction(models)
	if err != nil {
		http.Error(w, "failed to login with credentials", http.StatusInternalServerError)
		fmt.Printf("Error unable to login: %v", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "_auth",
		Value:    *user.RefreshToken,
		Path:     "/",
		MaxAge:   3600 * 5,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	user.RefreshToken = nil
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Print(w)
	json.NewEncoder(w).Encode(user)
}

// District create, update, delete
func (h *AuthHandler) CreateDistrict(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:district")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestDistrict
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddDistrict(models)
	if err != nil {
		http.Error(w, "Unable to create district", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateDistrict(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestDistrict
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.UpdateDistrict(models)
	if err != nil {
		http.Error(w, "Unable to create district", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteDistrict(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteDistrict(models)
	if err != nil {
		http.Error(w, "Unable to create district", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// End Create, update, delete

// Student Create, update, delete
func (h *AuthHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:students")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterRequestStudents
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddStudent(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:students")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterRequestStudents
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.UpdateStudent(models)
	if err != nil {
		http.Error(w, "Unable to update student", http.StatusInternalServerError)
		fmt.Printf("Unable to update student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:students")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteStudent(models)
	if err != nil {
		http.Error(w, "Unable to delete student", http.StatusInternalServerError)
		fmt.Printf("Unable to delete student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// END Student Create, update, delete

// Program create, update, delete AUTH
func (h *AuthHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:program")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddProgram(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Program create, update, delete AUTH
func (h *AuthHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:permissions")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterPermissionRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.CreatePermission(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:tutor")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterSchedule
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddSchedule(models)
	if err != nil {
		http.Error(w, "Unable to create AddSchedule", http.StatusInternalServerError)
		fmt.Printf("Unable to insert AddSchedule: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:program")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.UpdateProgram(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:program")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteProgram(models)
	if err != nil {
		http.Error(w, "Unable to Unable to delete program", http.StatusInternalServerError)
		fmt.Printf("Unable to delete program: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteProgramLocation(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:program-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveLocationProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteProgramLocation(models)
	if err != nil {
		http.Error(w, "Unable to Unable to delete program", http.StatusInternalServerError)
		fmt.Printf("Unable to delete program: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// END Program create, update, delete

func (h *AuthHandler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestMaterials
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddMaterial(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestMaterials
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.UpdateMaterial(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteMaterial(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Location create, update, delete
func (h *AuthHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddLocation(models)
	if err != nil {
		http.Error(w, "Unable to create location", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (h *AuthHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.UpdateLocation(models)
	if err != nil {
		http.Error(w, "Unable to create location", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (h *AuthHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteLocation(models)
	if err != nil {
		http.Error(w, "Unable to create location", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// END Location create, update, delete

// Semester create, update, delete
func (h *AuthHandler) CreateSemesterLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestSemesterLocation
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddSemesterLocation(models)
	if err != nil {
		http.Error(w, "Unable to create Semester", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateSemesterLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestSemesterLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.UpdateSemesterLocation(models)
	if err != nil {
		http.Error(w, "Unable to create Semester", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSemesterLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteSemesterLocation(models)
	if err != nil {
		http.Error(w, "Unable to delete Semester", http.StatusInternalServerError)
		fmt.Printf("Unable to delete location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Create organization
func (h *AuthHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterOrganization
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddOrganization(models)
	if err != nil {
		http.Error(w, "Unable to add organization", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Semester create, update, delete
func (h *AuthHandler) CreateSemester(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestSemester
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddSemester(models)
	if err != nil {
		http.Error(w, "Unable to create Semester", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateSemester(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestSemester
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.UpdateSemester(models)
	if err != nil {
		http.Error(w, "Unable to create Semester", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSemester(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteSemester(models)
	if err != nil {
		http.Error(w, "Unable to delete Semester", http.StatusInternalServerError)
		fmt.Printf("Unable to delete location: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// End create, update, delete
// Admin Create, update, delete
func (h *AuthHandler) CreateAdminStaff(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:admin")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterRequestAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddAdminStaff(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateAdminStaff(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterRequestAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.UpdateAdminStaff(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteAdminStaff(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RemoveAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteAdminStaff(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

//END Admin Create, update, delete

func (h *AuthHandler) CreateTutor(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterRequestTutor
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddTutor(models)
	if err != nil {
		http.Error(w, "Unable to create tutor staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert tutor staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateTutor(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterRequestTutor
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.UpdateTutor(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteTutor(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteTutor(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:locations")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetLocationsByID(idd, role)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionAccountability(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	tid := query.Get("tutor_id")
	start_date := query.Get("start_date")
	end_Date := query.Get("end_date")

	if email == "" || role == "" || id == "" || org_id == "" || start_date == "" || end_Date == "" || tid == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	tutor_id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	var model models.RequestTutorAccountability
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
			return
		}
		model.EndDate = end_date
	}
	model.TutorID = &tutor_id
	model.OrganizationID = &idd
	rows, err := h.authService.GetTutorSessionsAccountability(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSubjectLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	loc_id := query.Get("location_id")

	if email == "" || role == "" || id == "" || org_id == "" || loc_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	ldd, err := strconv.ParseInt(loc_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetSubjectByLocation(idd, ldd)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetLocationPrograms(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	loc_id := query.Get("location_id")

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:program-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" || loc_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	locId, err := strconv.ParseInt(loc_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	org, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetProgramsByLocation(locId, org)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) CreateProgramLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:program-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterLocationProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.CreateProgramLocation(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	lid := query.Get("locationId")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:students")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || lid == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	locationId, err := strconv.ParseInt(lid, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse location id", http.StatusInternalServerError)
		return
	}
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse location id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetStudentsByID(idd, role, locationId, tid)
	if err != nil {
		fmt.Printf("Error unable to login: %v", err)
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAdmins(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:admin")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetAdminStaffById(idd, role)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetDistricts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:district")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetDistrictsById(idd, role)
	if err != nil {
		http.Error(w, "Unable to get districts by id", http.StatusInternalServerError)
		fmt.Printf("error scanning row: %v", err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetPrograms(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:program")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetProgramsId(idd, role)
	if err != nil {
		http.Error(w, "Unable to get programs by id", http.StatusInternalServerError)
		fmt.Printf("error scanning row: %v", err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSemesterLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	location_id := query.Get("location_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:semester-locations")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	lod, err := strconv.ParseInt(location_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetSemesterLocationById(role, lod, idd)
	if err != nil {
		http.Error(w, "Unable to get programs by id", http.StatusInternalServerError)
		fmt.Printf("error scanning row: %v", err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")

	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetMaterialsById(idd, role)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in materials %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	location_id := query.Get("location_id")
	org_id := query.Get("organization_id")

	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	locid, err := strconv.ParseInt(location_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetTutorsById(idd, role, locid)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in materials %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	tutor_id := query.Get("tutor_id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	tid, err := strconv.ParseInt(tutor_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetTutorSchedules(tid)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in materials %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSemesters(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:semester")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetSemestersById(idd, role)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:assessments")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetAssessmentsById(idd, role)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:subject")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetSubjectById(idd, role)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) CreateTutorLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterTutorLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddTutorLocation(models)
	if err != nil {
		http.Error(w, "Unable to create RegisterTutorLocation staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert RegisterTutorLocation staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetTutorLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetTutorLocations(tid, idd)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) DeleteTutorLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RemoveTutorLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteTutorLocation(models)
	if err != nil {
		http.Error(w, "Unable to delete DeleteTutorLocation staff", http.StatusInternalServerError)
		fmt.Printf("Unable to delete DeleteTutorLocation staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RemoveSchedule
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteSchedule(models)
	if err != nil {
		http.Error(w, "Unable to delete DeleteTutorLocation staff", http.StatusInternalServerError)
		fmt.Printf("Unable to delete DeleteTutorLocation staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:subject")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterSubject
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddSubject(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:subject")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterSubject
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.UpdateSubject(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:subject")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteSubject(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateStudentSession(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterStudentSessionList
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.CreateStudentSessions(models)
	if err != nil {
		issue := fmt.Sprintf("Unable to create session: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateAssessment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:assessments")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterAssessment
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddAssessment(models)
	if err != nil {
		http.Error(w, "Unable to create assessment", http.StatusInternalServerError)
		fmt.Printf("Unable to insert assessment: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateAssessment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:assessments")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterAssessment
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.UpdateAssessment(models)
	if err != nil {
		http.Error(w, "Unable to create UpdateAnnouncement", http.StatusInternalServerError)
		fmt.Printf("Unable to insert UpdateAnnouncement: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (h *AuthHandler) DeleteAssessment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:assessments")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteAssessment(models)
	if err != nil {
		http.Error(w, "Unable to create UpdateAnnouncement", http.StatusInternalServerError)
		fmt.Printf("Unable to insert UpdateAnnouncement: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:announcements")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterAnnouncements
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddAnnouncement(models)
	if err != nil {
		http.Error(w, "Unable to AddAnnouncement", http.StatusInternalServerError)
		fmt.Printf("Unable to insert AddAnnouncement: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:announcements")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteAnnouncement(models)
	if err != nil {
		http.Error(w, "Unable to create DeleteAnnouncement", http.StatusInternalServerError)
		fmt.Printf("Unable to insert DeleteAnnouncement: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:announcements")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterUpdateAnnouncements
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.UpdateAnnouncement(models)
	if err != nil {
		http.Error(w, "Unable to create UpdateAnnouncement", http.StatusInternalServerError)
		fmt.Printf("Unable to insert UpdateAnnouncement: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetSessionSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	// Undefined variables like optional location_id
	// Check if exist before
	if !query.Has("email") || !query.Has("role") || !query.Has("id") || !query.Has("organization_id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.SearchQuery

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationId = &loc_id
	}
	if query.Get("program_id") != "" {
		prog_id, err := strconv.ParseInt(query.Get("program_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.ProgramId = &prog_id
	}
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}
	if query.Get("subject_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("subject_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SubjectId = &sub_id
	}
	if query.Get("semester_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sub_id
	}

	if query.Get("date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
			return
		}
		model.DateStart = start_time
	}
	if query.Get("date_end") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("date_end"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
			return
		}
		model.DateEnd = end_date
	}
	rows, err := h.authService.SessionSearch(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentSessionSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	// Undefined variables like optional location_id
	// Check if exist before
	if !query.Has("email") || !query.Has("role") || !query.Has("id") || !query.Has("organization_id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.SearchQuery

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationId = &loc_id
	}
	if query.Get("program_id") != "" {
		prog_id, err := strconv.ParseInt(query.Get("program_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.ProgramId = &prog_id
	}
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}
	if query.Get("subject_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("subject_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SubjectId = &sub_id
	}
	if query.Get("semester_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sub_id
	}

	if query.Get("date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
			return
		}
		model.DateStart = start_time
	}
	if query.Get("date_end") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("date_end"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
			return
		}
		model.DateEnd = end_date
	}
	rows, err := h.authService.StudentSessionSearch(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if !query.Has("email") || !query.Has("role") || !query.Has("id") || !query.Has("organization_id") || !query.Has("search_term") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.SearchQueryTutor
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}

	if query.Get("search_term") != "" {
		model.SearchTerm = query.Get("search_term")
	}

	rows, err := h.authService.TutorSearch(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorsSessions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	semester_id := query.Get("semester_id")
	location_id := query.Get("location_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var ss models.RequestTutorsSessions
	if email == "" || role == "" || id == "" || semester_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(semester_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	ss.SemesterID = &idd
	locid, err := strconv.ParseInt(location_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	ss.LocationID = &locid
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	ss.ID = &tid
	rows, err := h.authService.GetSessionsTutors(ss)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	session_id := query.Get("session_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" || session_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(session_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	// THIS SHOULD RETURN SESSION BASED ON TUTOR_ID
	rows, err := h.authService.SessionInfo(idd)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	// THIS SHOULD RETURN SESSIONS BASED ON SESSIONS ABOVE
	a_rows, err := h.authService.AssessmentInfo(idd)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows, "assessment_data": a_rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	student_id := query.Get("student_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || id == "" || org_id == "" || student_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(student_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	organization_id, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	// I want to return all the session and their information
	rows, err := h.authService.TrailSessions(idd)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	// Also return all assessments
	a_rows, err := h.authService.StudentAssessmentInfo(idd, organization_id)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows, "assessment_data": a_rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	org_id := query.Get("organization_id")
	id := query.Get("id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:permissions")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	aeid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetPermissionsById(idd, role, aeid)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetPermissionsById %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetOrganizationPermissions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:permissions")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if email == "" || role == "" || org_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetOrganizationPermissions(idd)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetOrganizationPermissions %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	// This needs to be handled in a PUT request.
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	location_ids := query.Get("location_ids")
	program_ids := query.Get("program_ids")
	org_id := query.Get("organization_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:announcements")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	oid, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	idd, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	locations, err := ParseArrayParam(location_ids)
	if err != nil {
		http.Error(w, "Unable to ParseArrayParam", http.StatusInternalServerError)
		return
	}
	programs, err := ParseArrayParam(program_ids)
	if err != nil {
		http.Error(w, "Unable to ParseArrayParam", http.StatusInternalServerError)
		return
	}
	var models models.AnnouncementRequest
	models.OrganizationID = oid
	models.ID = idd
	models.Role = role
	models.Email = email
	models.LocationIDs = locations
	models.ProgramID = programs

	rows, err := h.authService.GetAnnouncements(models)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in Semesters %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) CreateSubjectLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:subject-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterSubjectLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddSubjectLocation(models)
	if err != nil {
		http.Error(w, "Unable to create student", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSubjectLocation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:subject-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RemoveSubjectLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteSubjectLocation(models)
	if err != nil {
		http.Error(w, "Unable to create district", http.StatusInternalServerError)
		fmt.Printf("Unable to insert student: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func ParseArrayParam(param string) ([]int64, error) {
	if param == "" {
		return []int64{}, nil
	}

	// Split by comma
	stringValues := strings.Split(param, ",")

	// Convert to integers
	var intValues []int64
	for _, strVal := range stringValues {
		num, err := strconv.ParseInt(strings.TrimSpace(strVal), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number in array: %v", err)
		}
		intValues = append(intValues, num)
	}

	return intValues, nil
}
