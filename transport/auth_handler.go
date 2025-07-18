package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tracker/app/models"
	"tracker/app/services"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService *services.AuthService
	s3Service   *s3.Client
}

func NewAuthHandler(authService *services.AuthService, s3client *s3.Client) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		s3Service:   s3client,
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	user, err := h.authService.LoginAction(ctx, models)
	if err != nil {
		http.Error(w, "failed to login with credentials", http.StatusInternalServerError)
		fmt.Printf("Error unable to login: %v", err)
		return
	}
	expiresAt := time.Now().Add(4380 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "_auth",
		Value:    *user.RefreshToken,
		Path:     "/",
		Expires:  expiresAt,
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

// Register student assessment external

// END

// District create, update, delete
func (h *AuthHandler) CreateDistrict(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	user, err := h.authService.AddDistrict(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to add district", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateDistrict(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		fmt.Printf("error found reading body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestDistrict
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateDistrict(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update district", http.StatusInternalServerError)
		fmt.Printf("Unable to update district: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteDistrict(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		fmt.Printf("error found reading body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteDistrict(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete district", http.StatusInternalServerError)
		fmt.Printf("Unable to delete district: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// End Create, update, delete

// Student Create, update, delete
func (h *AuthHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		fmt.Printf("error reading body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddStudent(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		issue := fmt.Sprintf("Unable to create student: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		fmt.Printf("error reading body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateStudent(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update student", http.StatusInternalServerError)
		fmt.Printf("Unable to update student: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		fmt.Printf("error reading body: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteStudent(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete student", http.StatusInternalServerError)
		fmt.Printf("Unable to delete student: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// END Student Create, update, delete

// Program create, update, delete AUTH
func (h *AuthHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	user, err := h.authService.AddProgram(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to register program", http.StatusInternalServerError)
		fmt.Printf("Unable to delete student: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (h *AuthHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	user, err := h.authService.UpdateProgram(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to update program", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Program create, update, delete AUTH
func (h *AuthHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	user, err := h.authService.CreatePermission(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.DeleteProgram(ctx, models)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteProgramLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.DeleteProgramLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to delete program location", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateProgramLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.CreateProgramLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to create program location", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	user, err := h.authService.AddSchedule(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// END Program create, update, delete

func (h *AuthHandler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
<<<<<<< Updated upstream
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:material")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RegisterRequestMaterials
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddMaterial(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to create material", http.StatusInternalServerError)
		fmt.Printf("Unable to insert material: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
=======
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// — JWT & permission checks omitted —

	// Parse up to 10 MB of multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Try to grab "file"; ErrMissingFile means "no file uploaded"
	file, _, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "error reading file", http.StatusBadRequest)
		return
	}
	// Only close if we actually got a file
	if file != nil {
		defer file.Close()
	}

	// Read the JSON payload from the "data" field
	jsonField := r.FormValue("data")
	if jsonField == "" {
		http.Error(w, "missing JSON data", http.StatusBadRequest)
		return
	}
	var payload models.RegisterRequestMaterials
	if err := json.Unmarshal([]byte(jsonField), &payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// If a file was provided, upload it to S3 and set the SReference
	if file != nil {
		keyPtr, err := h.authService.UploadMaterialFile(ctx, file, nil)
		if err != nil {
			http.Error(w, "unable to upload to s3", http.StatusInternalServerError)
			return
		}
		payload.SReference = keyPtr
	}

	// Persist whatever we’ve got (with or without S3 key)
	user, err := h.authService.AddMaterial(ctx, payload)
	if err != nil {
		http.Error(w, "unable to create material", http.StatusInternalServerError)
		return
	}

>>>>>>> Stashed changes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
<<<<<<< Updated upstream
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
=======
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	// Parse up to 10 MB of multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
>>>>>>> Stashed changes
		return
	}

	// Try to grab "file"; ErrMissingFile means "no file uploaded"
	file, _, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "error reading file", http.StatusBadRequest)
		return
	}
	// Only close if we actually got a file
	if file != nil {
		defer file.Close()
	}

<<<<<<< Updated upstream
	var models models.RegisterRequestMaterials
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
=======
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:material")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	jsonField := r.FormValue("data")
	if jsonField == "" {
		http.Error(w, "missing JSON data", http.StatusBadRequest)
		return
>>>>>>> Stashed changes
	}
	var models models.RegisterRequestMaterials
	if err := json.Unmarshal([]byte(jsonField), &models); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if models.SReferenceDelete {
		stringPtr, err := h.authService.DoesReferenceExist(ctx, models.ID)
		if err != nil {
			http.Error(w, "unable to find reference value", http.StatusInternalServerError)
			return
		}
		if stringPtr != nil {
			err := h.authService.DeleteObjectS3(ctx, *stringPtr)
			if err != nil {
				http.Error(w, "unable to delete from s3 ", http.StatusInternalServerError)
				return
			}
			models.SReference = nil
			r.MultipartForm.RemoveAll()
			file = nil
		}
	}

	if file != nil {
		stringPtr, err := h.authService.DoesReferenceExist(ctx, models.ID)
		if err != nil {
			http.Error(w, "unable to find reference value", http.StatusInternalServerError)
			return
		}

		keyPtr, err := h.authService.UploadMaterialFile(ctx, file, stringPtr)
		if err != nil {
			http.Error(w, "unable to upload to s3", http.StatusInternalServerError)
			return
		}
		models.SReference = keyPtr
	}

	user, err := h.authService.UpdateMaterial(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update material", http.StatusInternalServerError)
		fmt.Printf("Unable to update material: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:material")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteMaterial(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete material", http.StatusInternalServerError)
		fmt.Printf("Unable to delete material: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Location create, update, delete
func (h *AuthHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RegisterRequestLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to create location", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RegisterRequestLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update location", http.StatusInternalServerError)
		fmt.Printf("Unable to update location: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete location", http.StatusInternalServerError)
		fmt.Printf("Unable to delete location: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// END Location create, update, delete

// Semester create, update, delete
func (h *AuthHandler) CreateSemesterLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:semester-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RegisterRequestSemesterLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddSemesterLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to create semester location", http.StatusInternalServerError)
		fmt.Printf("Unable to insert location: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateSemesterLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:semester-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RegisterRequestSemesterLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateSemesterLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update semester location", http.StatusInternalServerError)
		fmt.Printf("Unable to update location: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSemesterLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:semester-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteSemesterLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete semester location", http.StatusInternalServerError)
		fmt.Printf("Unable to delete location: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Create organization
func (h *AuthHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.AddOrganization(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Semester create, update, delete
func (h *AuthHandler) CreateSemester(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.AddSemester(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateSemester(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.UpdateSemester(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSemester(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	user, err := h.authService.DeleteSemester(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// End create, update, delete

// Admin Create, update, delete
func (h *AuthHandler) CreateAdminStaff(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

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

	var models models.RegisterRequestAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddAdminStaff(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to create admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert admin staff: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateAdminStaff(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RegisterRequestAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateAdminStaff(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to update admin staff: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteAdminStaff(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	var models models.RemoveAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteAdminStaff(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to delete admin staff: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// END Admin Create, update, delete
// Tutor Create, update, delete
func (h *AuthHandler) CreateTutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddTutor(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		issue := fmt.Sprintf("Error found: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateTutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateTutor(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update tutor staff", http.StatusInternalServerError)
		fmt.Printf("Unable to update tutor staff: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteTutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteTutor(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete tutor staff", http.StatusInternalServerError)
		fmt.Printf("Unable to delete tutor staff: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteTutorLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	valid, err := validateRequest(claims, "delete:tutor-locations")
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
	user, err := h.authService.DeleteTutorLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateTutorLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	valid, err := validateRequest(claims, "write:tutor-locations")
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
	user, err := h.authService.AddTutorLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// End of tutor create update, delete

// Create subject schedule
func (h *AuthHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.DeleteSchedule(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

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

	var models models.RegisterSubject
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddSubject(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to create subject", http.StatusInternalServerError)
		fmt.Printf("Unable to create subject: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

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

	var models models.RegisterSubject
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateSubject(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update subject", http.StatusInternalServerError)
		fmt.Printf("Unable to update subject: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

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

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteSubject(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete subject", http.StatusInternalServerError)
		fmt.Printf("Unable to delete subject: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateSubjectLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.AddSubjectLocation(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteSubjectLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	user, err := h.authService.DeleteSubjectLocation(ctx, models)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// End of create, update, delete subject

// Announcements create, update, delete
func (h *AuthHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddAnnouncement(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to create announcement", http.StatusInternalServerError)
		fmt.Printf("Unable to create announcement: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateAnnouncement(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update announcement", http.StatusInternalServerError)
		fmt.Printf("Unable to update announcement: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteAnnouncement(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete announcement", http.StatusInternalServerError)
		fmt.Printf("Unable to delete announcement: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// End of announcments
func (h *AuthHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("unable to parse body \n")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:teacher")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var model models.RegisterTeacher
	if err := json.Unmarshal(body, &model); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	user, err := h.authService.AddTeacher(ctx, model)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("unable to parse body \n")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:teacher")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var model models.RegisterTeacher
	if err := json.Unmarshal(body, &model); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	user, err := h.authService.UpdateTeacher(ctx, model)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	fmt.Printf("Request body %s\n", string(body))

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:teacher")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RegisterTeacher
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteTeacher(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete assessment", http.StatusInternalServerError)
		fmt.Printf("Unable to delete assessment: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Assessments create, delete, update

func (h *AuthHandler) CreateAssessment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddAssessment(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to create assessment", http.StatusInternalServerError)
		fmt.Printf("Unable to create assessment: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateAssessment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.UpdateAssessment(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to update assessment", http.StatusInternalServerError)
		fmt.Printf("Unable to update assessment: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteAssessment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteAssessment(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to delete assessment", http.StatusInternalServerError)
		fmt.Printf("Unable to delete assessment: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Assessments create update delete end
func (h *AuthHandler) GetLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	valid, err := validateRequest(claims, "view:location")
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

	rows, err := h.authService.GetLocationsByID(ctx, idd, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to get locations", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionAccountability(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.GetTutorSessionsAccountability(ctx, model)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSubjectLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	rows, err := h.authService.GetSubjectByLocation(ctx, idd, ldd)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetLocationPrograms(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.GetProgramsByLocation(ctx, locId, org)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.GetStudentsByID(ctx, idd, role, locationId, tid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	session_id := query.Get("session_id")
	student_id := query.Get("student_id")
	stud_id, err := strconv.ParseInt(student_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse location id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetStudentDetails(ctx, &session_id, &stud_id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAdmins(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.GetAdminStaffById(ctx, idd, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetDistricts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	orgID := query.Get("organization_id")

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

	if email == "" || role == "" || id == "" || orgID == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		fmt.Printf("Error parsing organization_id: %v\n", err)
		return
	}

	rows, err := h.authService.GetDistrictsById(ctx, orgIDParsed, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to get districts", http.StatusInternalServerError)
		fmt.Printf("Error retrieving districts: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetPrograms(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	orgID := query.Get("organization_id")

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

	if email == "" || role == "" || id == "" || orgID == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		fmt.Printf("Error parsing organization_id: %v\n", err)
		return
	}

	rows, err := h.authService.GetProgramsId(ctx, orgIDParsed, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to get programs", http.StatusInternalServerError)
		fmt.Printf("Error retrieving programs: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSemesterLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	orgID := query.Get("organization_id")
	locationID := query.Get("location_id")

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:semester-location")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if email == "" || role == "" || id == "" || orgID == "" || locationID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	locationIDParsed, err := strconv.ParseInt(locationID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid location_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetSemesterLocationById(ctx, role, locationIDParsed, orgIDParsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to get semester locations", http.StatusInternalServerError)
		fmt.Printf("Error retrieving semester locations: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	orgID := query.Get("organization_id")

	if email == "" || role == "" || id == "" || orgID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetMaterialsById(ctx, orgIDParsed, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to get materials", http.StatusInternalServerError)
		fmt.Printf("Error retrieving materials: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
<<<<<<< Updated upstream
=======
	locationID := query.Get("location_id")
	orgID := query.Get("organization_id")

	if email == "" || role == "" || id == "" || orgID == "" || locationID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	locationIDParsed, err := strconv.ParseInt(locationID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid location_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetTutorsById(ctx, orgIDParsed, role, locationIDParsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to get tutors", http.StatusInternalServerError)
		fmt.Printf("Error retrieving tutors: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSignedUrlMaterials(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	uuid := query.Get("id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:material")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if uuid == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	url, err := h.authService.GenerateMaterialsPresignedUrl(ctx, uuid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to get semesters", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": url}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSemesters(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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

	rows, err := h.authService.GetSemestersById(ctx, idd, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to get semesters", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	tutorID := query.Get("tutor_id")
	orgID := query.Get("organization_id")

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

	if email == "" || role == "" || id == "" || orgID == "" || tutorID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	tutorIDParsed, err := strconv.ParseInt(tutorID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tutor_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetTutorSchedules(ctx, tutorIDParsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to retrieve schedules", http.StatusInternalServerError)
		fmt.Printf("Error retrieving schedules: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	orgID := query.Get("organization_id")

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

	if email == "" || role == "" || id == "" || orgID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetAssessmentsById(ctx, orgIDParsed, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to retrieve assessments", http.StatusInternalServerError)
		fmt.Printf("Error retrieving assessments: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentQuestions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	assessmentID := query.Get("assessment_id")

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

	if assessmentID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	assessmentIDParsed, err := strconv.ParseInt(assessmentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid assessment_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetAssessmentsQuestionsChoice(ctx, assessmentIDParsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to retrieve assessment questions", http.StatusInternalServerError)
		fmt.Printf("Error retrieving assessment questions: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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

	rows, err := h.authService.GetSubjectById(ctx, idd, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
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
	if query.Get("location_id") == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetTeachers(ctx, &idd)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	valid, err := validateRequest(claims, "view:tutor-locations")
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
	rows, err := h.authService.GetTutorLocations(ctx, tid, idd)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	rows, err := h.authService.SessionSearch(ctx, model)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentAssesssmentSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	sid := query.Get("student_assessment_id")
	studentAssessmentSearch, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	easyScore := query.Get("easy_score")
	easyScoreBool, err := strconv.ParseBool(easyScore)
	if err != nil {
		http.Error(w, "unable to parse boolean", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.StudentAssessmentSearch(ctx, &studentAssessmentSearch, easyScoreBool)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentSessionSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	if query.Get("search_term") != "" {
		model.SearchTerm = query.Get("search_term")
	}
	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse location", http.StatusInternalServerError)
			return
		}
		model.LocationId = &loc_id
	}
	if query.Get("program_id") != "" {
		prog_id, err := strconv.ParseInt(query.Get("program_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse program", http.StatusInternalServerError)
			return
		}
		model.ProgramId = &prog_id
	}
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse organization", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}
	if query.Get("subject_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("subject_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse subject", http.StatusInternalServerError)
			return
		}
		model.SubjectId = &sub_id
	}
	if query.Get("semester_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse semester", http.StatusInternalServerError)
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
	rows, err := h.authService.StudentSessionSearch(ctx, model)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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

	rows, err := h.authService.TutorSearch(ctx, model)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorsSessions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	semester_id := query.Get("semester_id")
>>>>>>> Stashed changes
	location_id := query.Get("location_id")
	org_id := query.Get("organization_id")

	if email == "" || role == "" || id == "" || orgID == "" || locationID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	locationIDParsed, err := strconv.ParseInt(locationID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid location_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetTutorsById(ctx, orgIDParsed, role, locationIDParsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to get tutors", http.StatusInternalServerError)
		fmt.Printf("Error retrieving tutors: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSemesters(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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

	rows, err := h.authService.GetSemestersById(ctx, idd, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "Unable to get semesters", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	tutorID := query.Get("tutor_id")
	orgID := query.Get("organization_id")

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

	if email == "" || role == "" || id == "" || orgID == "" || tutorID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	tutorIDParsed, err := strconv.ParseInt(tutorID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid tutor_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetTutorSchedules(ctx, tutorIDParsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to retrieve schedules", http.StatusInternalServerError)
		fmt.Printf("Error retrieving schedules: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	email := query.Get("email")
	role := query.Get("role")
	id := query.Get("id")
	orgID := query.Get("organization_id")

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

	if email == "" || role == "" || id == "" || orgID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	orgIDParsed, err := strconv.ParseInt(orgID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetAssessmentsById(ctx, orgIDParsed, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to retrieve assessments", http.StatusInternalServerError)
		fmt.Printf("Error retrieving assessments: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentQuestions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	assessmentID := query.Get("assessment_id")

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

	if assessmentID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	assessmentIDParsed, err := strconv.ParseInt(assessmentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid assessment_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetAssessmentsQuestionsChoice(ctx, assessmentIDParsed)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to retrieve assessment questions", http.StatusInternalServerError)
		fmt.Printf("Error retrieving assessment questions: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSubjects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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

	rows, err := h.authService.GetSubjectById(ctx, idd, role)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
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
	if query.Get("location_id") == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetTeachers(ctx, &idd)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	valid, err := validateRequest(claims, "view:tutor-locations")
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
	rows, err := h.authService.GetTutorLocations(ctx, tid, idd)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	rows, err := h.authService.SessionSearch(ctx, model)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentAssesssmentSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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
	sid := query.Get("student_assessment_id")
	studentAssessmentSearch, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	easyScore := query.Get("easy_score")
	easyScoreBool, err := strconv.ParseBool(easyScore)
	if err != nil {
		http.Error(w, "unable to parse boolean", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.StudentAssessmentSearch(ctx, &studentAssessmentSearch, easyScoreBool)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentSessionSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	if query.Get("search_term") != "" {
		model.SearchTerm = query.Get("search_term")
	}
	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse location", http.StatusInternalServerError)
			return
		}
		model.LocationId = &loc_id
	}
	if query.Get("program_id") != "" {
		prog_id, err := strconv.ParseInt(query.Get("program_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse program", http.StatusInternalServerError)
			return
		}
		model.ProgramId = &prog_id
	}
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse organization", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}
	if query.Get("subject_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("subject_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse subject", http.StatusInternalServerError)
			return
		}
		model.SubjectId = &sub_id
	}
	if query.Get("semester_id") != "" {
		sub_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse semester", http.StatusInternalServerError)
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
	rows, err := h.authService.StudentSessionSearch(ctx, model)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorSearch(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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

	rows, err := h.authService.TutorSearch(ctx, model)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorsSessions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.GetSessionsTutors(ctx, ss)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	a_rows, err := h.authService.AssessmentInfo(ctx, idd)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows, "assessment_data": a_rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.TrailSessions(ctx, idd)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	// Also return all assessments
	a_rows, err := h.authService.StudentAssessmentInfo(ctx, idd, organization_id)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows, "assessment_data": a_rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.GetPermissionsById(ctx, idd, role, aeid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetOrganizationPermissions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	rows, err := h.authService.GetOrganizationPermissions(ctx, idd)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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

	rows, err := h.authService.GetAnnouncements(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		// 6. You could also inspect SQL errors here if you like.
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
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
	w.WriteHeader(http.StatusCreated)
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
