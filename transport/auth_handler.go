package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// District create, update, delete
func (h *AuthHandler) CreateDistrict(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "create:program")
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

func (h *AuthHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "update:program")
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
func (h *AuthHandler) CreateSemester(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterRequestTutor
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddTutor(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
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

	_, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	//fmt.Print(props["sub"].(string))
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

func (h *AuthHandler) GetLocationPrograms(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	org_id := query.Get("organization_id")
	loc_id := query.Get("location_id")

	_, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	//fmt.Print(props["sub"].(string))
	if email == "" || role == "" || id == "" || org_id == "" || loc_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	locId, err := strconv.ParseInt(loc_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	idd, err := strconv.ParseInt(org_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetProgramsByLocation(idd, locId, role)
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
	rows, err := h.authService.GetStudentsByID(idd, role, locationId)
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

func (h *AuthHandler) GetSemesters(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
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
	fmt.Printf("Request body %s\n", string(body))
	var models models.RegisterStudentSessionList
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.CreateStudentSessions(models)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
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

func (h *AuthHandler) GetSessionSearch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	fmt.Printf("Request body %s\n", string(body))
	var models models.SearchQuery
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	rows, err := h.authService.SessionSearch(models)
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
