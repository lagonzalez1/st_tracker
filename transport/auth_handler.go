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
	"tracker/app/cache"
	"tracker/app/helpers"
	"tracker/app/models"
	"tracker/app/services"
	"tracker/app/sqs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Stripe account details
// account_id: acct_1OPIeqFBv6mpTbIZ
// secret: wins-finely-quiet-cool

type AuthHandler struct {
	authService *services.AuthService
	cacheHander *cache.CacheHandler
	sqsHandler  *sqs.SqsHandler
	config      *models.Config
}

func NewAuthHandler(authService *services.AuthService, cacheHandler *cache.CacheHandler, sqsHandler *sqs.SqsHandler, config *models.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cacheHander: cacheHandler,
		sqsHandler:  sqsHandler,
		config:      config,
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

	var models models.LoginRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.LoginAction(ctx, models)
	if err != nil {
		fmt.Print(err)
		fmt.Printf("Error unable to login: %v", err)
		http.Error(w, "failed to login with credentials", http.StatusInternalServerError)
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	var models models.RegisterRequestDistrict
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	// Check if claims is the same as input data
	models.OrganizationId = &orgID
	cacheKey := fmt.Sprintf("get:districts:%d", orgID)
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
	// Clear cache
	h.cacheHander.ClearCache(ctx, cacheKey)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateDistrict(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	orgID, err := helpers.ExtractFloat64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
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
	var models models.RegisterRequestDistrict
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	// Check if claims is the same as input data
	if int64(orgID) != *models.OrganizationId {
		http.Error(w, "Invalid claims and input missmatch", http.StatusBadRequest)
		fmt.Printf("Claims is different from input data: %v", err)
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
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:district")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
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

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationId = orgID
	key := fmt.Sprintf("get:districts:%d", orgID)
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
	h.cacheHander.ClearCache(ctx, key)
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	var models models.RegisterRequestProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	models.OrganizationId = &orgID
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
	orgID, err := helpers.ExtractFloat64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}

	var models models.RegisterRequestProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	if int64(orgID) != *models.OrganizationId {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
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
	orgID, err := helpers.ExtractFloat64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	models.OrganizationId = int64(orgID)
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
	orgID, err := helpers.ExtractFloat64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	var models models.RegisterPermissionRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	if *models.OrganizationId != int64(orgID) {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}

	var models models.RemoveLocationProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	models.OrganizationID = &orgID
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	var models models.RegisterLocationProgram
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	models.OrganizationID = &orgID

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
	var models models.RegisterSchedule
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

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

func (h *AuthHandler) CreateScheduleV3(w http.ResponseWriter, r *http.Request) {
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
	var models models.RegisterScheduleLink
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddScheduleLink(ctx, models)
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

func (h *AuthHandler) DeleteEntitySchedule(w http.ResponseWriter, r *http.Request) {
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
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteScheduleLink(ctx, models)
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

func (h *AuthHandler) CreateScheduleV2(w http.ResponseWriter, r *http.Request) {
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
	var models models.RegisterScheduleV2
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	user, err := h.authService.AddScheduleV2(ctx, models)
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
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	// Parse up to 10 MB of multipart form data
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}

	var payload models.RegisterRequestMaterials
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if *payload.OrganizationId != int64(orgID) {
		http.Error(w, "Invalid request", http.StatusInternalServerError)
		return
	}
	key := fmt.Sprintf("get:materials:%d", orgID)
	// If a file was provided, upload it to S3 and set the SReference
	// If a file exist then get the presigned url with key
	var presigned_url *string
	if payload.File && payload.FileType != nil {
		url, key, err := h.authService.GeneratePutPresignedUrl(ctx, payload.FileType, 5)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "unable to upload to s3", http.StatusInternalServerError)
			return
		}
		presigned_url = url
		payload.SReference = key
	}
	// Persist whatever we’ve got (with or without S3 key)
	user, err := h.authService.AddMaterial(ctx, payload)
	if err != nil {
		http.Error(w, "unable to create material", http.StatusInternalServerError)
		return
	}
	h.cacheHander.ClearCache(ctx, key)
	user.UploadUrl = presigned_url
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to parse claims orgid", http.StatusBadRequest)
		return
	}
	var models models.RegisterRequestMaterials
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("get:materials:%d", orgid)
	// If there is a delete or update request delete the file associated with key
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
		}
	}
	var presigned_url *string
	// If there si another file attched then add
	if models.File && models.FileType != nil {
		url, key, err := h.authService.GeneratePutPresignedUrl(ctx, models.FileType, 5)
		if err != nil {
			http.Error(w, "unable to upload to s3", http.StatusInternalServerError)
			return
		}
		presigned_url = url
		models.SReference = key

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
	user.UploadUrl = presigned_url
	h.cacheHander.ClearCache(ctx, key)
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

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	key := fmt.Sprintf("get:materials%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to parse claims orgid", http.StatusBadRequest)
	}
	cacheKey := fmt.Sprintf("get:locations:%d", orgid)
	var models models.RegisterRequestLocation
	models.OrganizationId = &orgid

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

	h.cacheHander.ClearCache(ctx, cacheKey)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	key := fmt.Sprintf("get:locations:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
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

func (h *AuthHandler) CreateAdminLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
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

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return

	}

	var models models.RegisterAdminLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.AddAdminLocation(ctx, models, &orgid)
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

func (h *AuthHandler) DeleteAdminLocation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "delete:admin")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RemoveAdminLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	user, err := h.authService.DeleteAdminLocation(ctx, models)
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
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

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
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("get:semesters:%d", orgid)
	var models models.RegisterRequestSemester
	models.OrganizationId = &orgid
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
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
	h.cacheHander.ClearCache(ctx, key)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateSemester(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	var models models.RegisterRequestSemester
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	key := fmt.Sprintf("get:semesters:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
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
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	key := fmt.Sprintf("get:semesters:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateSemesterDates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	_, err = helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	var models models.RegisterSemesterDates
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddSemesterDates(ctx, models)
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

func (h *AuthHandler) CreateStudentGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	_, err = helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	var models models.RegisterStudentGroup
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddStudentGroup(ctx, models)
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

func (h *AuthHandler) UpdateStudentGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	_, err = helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	var models models.RegisterStudentGroup
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.AddStudentGroup(ctx, models)
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
func (h *AuthHandler) DeleteStudentGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	_, err = helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal([]byte(body), &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteStudentGroup(ctx, models)
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

func (h *AuthHandler) DeleteSemesterDates(w http.ResponseWriter, r *http.Request) {
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to semester", http.StatusBadRequest)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	key := fmt.Sprintf("get:semesters_dates:%d:%d", orgid, *models.ID)
	// Need to handle 3 cases of logins for different permissions
	user, err := h.authService.DeleteSemesterDates(ctx, models)
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
	h.cacheHander.ClearCache(ctx, key)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetSemesterDates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:semesters")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if query.Get("semester_id") == "" {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	sid, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("get:semesters_dates:%d:%d", orgid, sid)
	lkey := h.cacheHander.LockKey(key)
	cdata, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var rows []models.ResponseSemesterDates
		if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetSemesterDates(ctx, &sid)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "request timeout", http.StatusGatewayTimeout)
				return
			}
			if errors.Is(err, context.Canceled) {
				http.Error(w, "request canceled", http.StatusRequestTimeout)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			fmt.Printf("service error: %v\n", err)
			return
		}
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		cdata, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseSemesterDates
			if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
				http.Error(w, "Unable to parse id", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RegisterRequestAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationId = &orgid
	key := fmt.Sprintf("get:admins:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return
}

func (h *AuthHandler) UpdateAdminStaff(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}

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
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RemoveAdmin
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	key := fmt.Sprintf("get:admins:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	var models models.RegisterRequestTutor
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	key := fmt.Sprintf("get:tutors:%d:%d", orgid, *models.LocationId)
	keyall := fmt.Sprintf("get:tutors:%d:%d", orgid, -1)
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
	h.cacheHander.ClearCache(ctx, key)
	h.cacheHander.ClearCache(ctx, keyall)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to get orgid", http.StatusInternalServerError)
		return
	}

	var models models.RegisterRequestTutor
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	key := fmt.Sprintf("get:tutors:%d:%d", orgid, -1)
	key2 := fmt.Sprintf("get:tutors:%d:%d", orgid, *models.LocationId)
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
	h.cacheHander.ClearCache(ctx, key)
	h.cacheHander.ClearCache(ctx, key2)

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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RemoveTutorLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	models.OrganizationID = &orgid

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

func (h *AuthHandler) DeleteSurvey(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "delete:survey")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteSurvey(ctx, models)
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

func (h *AuthHandler) DeleteSession(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "delete:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	models.OrganizationId = orgid
	user, err := h.authService.DeleteSession(ctx, models)
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

func (h *AuthHandler) DeleteProgramSurvey(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "delete:survey-program")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	models.OrganizationId = orgid

	user, err := h.authService.DeleteProgramSurvey(ctx, models)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetProgramSurveys(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	program_id := query.Get("program_id")
	pid, err := strconv.ParseInt(program_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:surveys-program")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetProgramSurveysById(ctx, orgid, pid)
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

func (h *AuthHandler) CreateProgramSurvey(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "write:survey-program")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var models models.RegisterProgramSurvey
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	models.OrganizationID = &orgid

	user, err := h.authService.CreateProgramSurvey(ctx, models)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("get:subjects:%d", orgid)
	var models models.RegisterSubject
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationId = &orgid
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

	h.cacheHander.ClearCache(ctx, key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateGlobalSchedule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:global-schedule")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RegisterGlobalSchedule
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationID = int(orgid)
	user, err := h.authService.AddScheduleGlobal(ctx, models)
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

func (h *AuthHandler) UpdateGlobalSchedule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:global-schedule")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RegisterGlobalSchedule
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationID = int(orgid)
	user, err := h.authService.UpdateScheduleGlobal(ctx, models)
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

func (h *AuthHandler) DeleteGlobalSchedule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:global-schedule")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	user, err := h.authService.HardDeleteScheduleGlobal(ctx, models.ID, &orgid)
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

func (h *AuthHandler) CreateStudentGroupAttendies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}

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
	//key := fmt.Sprintf("get:student-group:%d", orgid)
	var models models.RegisterStudentGroupAttendies
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	user, err := h.authService.AddStudentGroupAttendies(ctx, models)
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

	//h.cacheHander.ClearCache(ctx, key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateLocationContact(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RegisterLocationContact
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	key := fmt.Sprintf("get:location-contact:%d:%d", orgid, *models.LocationID)
	models.OrganizationId = &orgid
	user, err := h.authService.CreateLocationContact(ctx, models)
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

	h.cacheHander.ClearCache(ctx, key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateLocationContact(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("get:location-contact:%d", orgid)
	var models models.RegisterLocationContact
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationId = &orgid
	user, err := h.authService.UpdateLocationContact(ctx, models)
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

	h.cacheHander.ClearCache(ctx, key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteLocationContact(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var models models.RemoveRequestLocation
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	key := fmt.Sprintf("get:location-contact:%d:%d", orgid, *models.LocationID)
	models.OrganizationId = orgid
	user, err := h.authService.DeleteLocationContact(ctx, models)
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

	h.cacheHander.ClearCache(ctx, key)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	key := fmt.Sprintf("get:subjects:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
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
	orgId, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing organization id", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationID = &orgId

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

func (h *AuthHandler) CreateS3Object(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	path := "assessment_images/"
	id, err := h.authService.CreateObjectS3(ctx, file, nil, &path)
	if err != nil {
		http.Error(w, "unable to upload to s3", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"id": id}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) DeleteS3Object(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body\n")
		return
	}
	var model models.DeleteImageRequest
	if err := json.Unmarshal([]byte(body), &model); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	err = h.authService.DeleteObjectS3(ctx, model.ID)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"id": model.ID}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

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
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to AddTeacher", http.StatusInternalServerError)
		fmt.Printf("Unable to AddTeacher: %v\n", err)
		return
	}
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
	if err != nil {
		http.Error(w, "Unable to update teacher", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateAckAnnouncements(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "write:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	tid, err := helpers.ExtractInt64Claim(claims, "id")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	oid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}
	var model models.RegisterAnnouncementsAck
	if err := json.Unmarshal(body, &model); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	user, err := h.authService.AnnouncementAck(ctx, &oid, &tid, model.AnnouncmentID)
	if err != nil {
		http.Error(w, "Unable to AnnouncementAck", http.StatusBadRequest)
		return
	}
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	var models models.RegisterAssessment
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	models.OrganizationID = &orgid
	key := fmt.Sprintf("get:assessments:%d", orgid)

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
	h.cacheHander.ClearCache(ctx, key)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
	}

	var models models.RegisterAssessment
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	key := fmt.Sprintf("get:assessments:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)

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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to prase orgid", http.StatusForbidden)
		return
	}
	var models models.RemoveRequest
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	key := fmt.Sprintf("get:assessments:%d", orgid)
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
	h.cacheHander.ClearCache(ctx, key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Assessments create update delete end
func (h *AuthHandler) GetLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	key := fmt.Sprintf("get:locations:%d", orgid)
	lkey := h.cacheHander.LockKey(key)

	data, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var res []models.ResponseRequestLocations
		if err := json.Unmarshal([]byte(data), &res); err != nil {
			http.Error(w, "Unable to unmarshal cache data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Example response
		response := map[string]interface{}{"data": res}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetLocationsByID(ctx, orgid, role)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "request timeout", http.StatusGatewayTimeout)
				return
			}
			if errors.Is(err, context.Canceled) {
				http.Error(w, "request canceled", http.StatusRequestTimeout)
				return
			}
			http.Error(w, "Unable to get locations", http.StatusInternalServerError)
			fmt.Printf("service error: %v\n", err)
			return
		}
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		res, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseRequestLocations
			err = json.Unmarshal([]byte(res), &rows)
			if err != nil {
				fmt.Printf("unable to parse byte stream %v", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func (h *AuthHandler) GetLocationContact(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	if query.Get("location_id") == "" {
		http.Error(w, "no location specified", http.StatusInternalServerError)
		return
	}
	lid := query.Get("location_id")
	locid, err := strconv.ParseInt(lid, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse location id", http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("get:location-contact:%d:%d", orgid, locid)
	lkey := h.cacheHander.LockKey(key)

	data, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var res []models.ResponseLocationContact
		if err := json.Unmarshal([]byte(data), &res); err != nil {
			http.Error(w, "Unable to unmarshal cache data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Example response
		response := map[string]interface{}{"data": res}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetLocationContactByID(ctx, &orgid, &locid)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "request timeout", http.StatusGatewayTimeout)
				return
			}
			if errors.Is(err, context.Canceled) {
				http.Error(w, "request canceled", http.StatusRequestTimeout)
				return
			}
			http.Error(w, "Unable to get locations", http.StatusInternalServerError)
			fmt.Printf("service error: %v\n", err)
			return
		}
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		res, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseLocationContact
			err = json.Unmarshal([]byte(res), &rows)
			if err != nil {
				fmt.Printf("unable to parse byte stream %v", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func (h *AuthHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:*")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	org, err := h.authService.GetOrganizationById(ctx, &orgid)
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
		http.Error(w, "Unable to GetOrganizationById", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"organization": org}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetGenerationUsage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:*")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	materials, questions, err := h.authService.GetGenerationUsage(ctx, &orgid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to GetOrganizationById", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"materials_usage": materials, "questions_usage": questions}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetEntitySchedule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	tutor_id := query.Get("tutor_id")
	if tutor_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	valid, err := validateRequest(claims, "view:tutors")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	tid, err := strconv.ParseInt(tutor_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusBadRequest)
		return
	}

	data, err := h.authService.GetEntitySchedule(ctx, &tid, &orgid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to GetOrganizationById", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": data}
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

func (h *AuthHandler) GetEntityScheduleShift(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	uid, err := helpers.ExtractInt64Claim(claims, "id")
	if err != nil {
		http.Error(w, "missing id value", http.StatusBadRequest)
		return
	}
	// Need to cache this
	/**
		// Level 1: Individual schedule data (rarely changes)
		schedule:data:{id} = { schedule details }

		// Level 2: Tutor's combined schedule (depends on subscriptions)
		tutor:schedule:{tutorId} = { combined view }

		// Level 3: Schedule-to-tutor mapping
		schedule:subscribers:{scheduleId} = Set[tutor1, tutor2, tutor3]


		1. First set the schedule:data:{id} => Schedule details

		2. Then set the tutors to subscribe to each schedule => tutor:schedule:{schedule_id} = { schedules }
		** This is the mapping of Table(Tutor_Schedule_Assignment)

		3. Then the mappings to each subscriber -> subscriber:schedule:{schedule_id} = SET[tutor_id, tutor_id]

	**/
	rows, err := h.authService.GetEntityScheduleList(ctx, &uid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to GetTutorSchedule", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionScheduled(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	program_ids := query["program_ids[]"]
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	uid, err := helpers.ExtractInt64Claim(claims, "id")
	if err != nil {
		http.Error(w, "missing id value", http.StatusBadRequest)
		return
	}
	if len(program_ids) == 0 {
		http.Error(w, "missing program ids", http.StatusBadRequest)
		return
	}
	pids := make([]int64, len(program_ids))
	for i, v := range program_ids {
		pids[i], _ = strconv.ParseInt(v, 10, 64)
	}
	// Need to cache this
	// Cache on tutor_id:semester_id
	// Cache invalidate on -> Tutor submits a new session
	// Cache invalidate on -> Admin submits a new schedule
	rows, err := h.authService.GetTutorSchedule(ctx, &uid, pids)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "Unable to GetTutorSchedule", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSubjectLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	loc_id := query.Get("location_id")
	if loc_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	ldd, err := strconv.ParseInt(loc_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetSubjectByLocation(ctx, orgid, ldd)
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
	loc_id := query.Get("location_id")
	// Implement cache programs:org_id:loc_id on this endpoint using elasticache
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
	lid, err := strconv.ParseInt(loc_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	org, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetProgramsByLocation(ctx, lid, org)
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
	lid := query.Get("locationId")
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
	if lid == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	locationId, err := strconv.ParseInt(lid, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse location id", http.StatusInternalServerError)
		return
	}
	tid, err := helpers.ExtractInt64Claim(claims, "id")
	if err != nil {
		http.Error(w, "Unable to parse location id", http.StatusInternalServerError)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Unable to parse location id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetStudentsByID(ctx, orgid, role, locationId, tid)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	key := fmt.Sprintf("get:admins:%d", orgid)
	lkey := h.cacheHander.LockKey(key)

	cdata, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var rows []models.ResponseRequestAdminList
		if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetAdminStaffById(ctx, orgid, role)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "request timeout", http.StatusGatewayTimeout)
				return
			}
			if errors.Is(err, context.Canceled) {
				http.Error(w, "request canceled", http.StatusRequestTimeout)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			fmt.Printf("service error: %v\n", err)
			return
		}
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		cdata, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseRequestAdminList
			if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
				http.Error(w, "Unable to parse id", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

}

func (h *AuthHandler) GetDistricts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		fmt.Printf("Error parsing organization_id: %v\n", err)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		fmt.Printf("Error parsing role: %v\n", err)
		return
	}
	key := fmt.Sprintf("get:districts:%d", orgid)
	lkey := h.cacheHander.LockKey(key)
	res, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var rows []models.ResponseRequestDistrictList
		err = json.Unmarshal([]byte(res), &rows)
		if err != nil {
			fmt.Printf("unable to parse byte stream %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": &rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetDistrictsById(ctx, orgid, role)
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
		_, err = h.cacheHander.SetCache(ctx, key, rows)
		if err != nil {
			fmt.Printf("SetCache error %v", err)
			return
		}
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		res, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseRequestDistrictList
			err = json.Unmarshal([]byte(res), &rows)
			if err != nil {
				fmt.Printf("unable to parse byte stream %v", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func (h *AuthHandler) GetPrograms(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

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

	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}

	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetProgramsId(ctx, orgID, role)
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
	if locationID == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	locid, err := strconv.ParseInt(locationID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid location_id", http.StatusBadRequest)
		return
	}

	rows, err := h.authService.GetSemesterLocationById(ctx, locid, orgid)
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

// Can this result in issues since we are caching the s3 reference link?
func (h *AuthHandler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}

	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("get:materials:%d", orgid)
	lkey := h.cacheHander.LockKey(key)

	cdata, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var rows []models.ResponseRequestMaterialsList
		if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
			http.Error(w, "Unable to parse cache", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetMaterialsById(ctx, orgid, role)
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
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		cdata, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseRequestMaterialsList
			if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
				http.Error(w, "Unable to parse cache", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

}

func (h *AuthHandler) GetTutors(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	location_id := query.Get("location_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}
	locid, err := strconv.ParseInt(location_id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid location_id", http.StatusBadRequest)
		return
	}
	var key *string
	var lkey *string
	if locid <= 0 {
		l := fmt.Sprintf("get:tutors:%d:%d", orgid, locid)
		key = &l
		k := h.cacheHander.LockKey(*key)
		lkey = &k
	} else {
		l := fmt.Sprintf("get:tutors:%d:%d", orgid, locid)
		key = &l
		k := h.cacheHander.LockKey(*key)
		lkey = &k
	}

	cdata, isHit := h.cacheHander.CheckCache(ctx, *key)
	if isHit {
		var rows []models.ResponseRequestTutorsList
		if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
			http.Error(w, "Invalid location_id", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return

	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, *lkey, token)
	if lock {
		rows, err := h.authService.GetTutorsById(ctx, orgid, role, locid)
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
		h.cacheHander.SetCache(ctx, *key, rows)
		h.cacheHander.SafeUnlock(ctx, *lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, *key)
	if status {
		cdata, isHit := h.cacheHander.CheckCache(ctx, *key)
		if isHit {
			var rows []models.ResponseRequestTutorsList
			if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
				http.Error(w, "Invalid location_id", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

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
	var objectKey = "materials/" + uuid
	url, err := h.authService.GenerateMaterialsPresignedUrl(ctx, objectKey)
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("get:semesters:%d", orgid)
	lkey := h.cacheHander.LockKey(key)
	data, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var rows []models.ResponseRequestSemesterList
		if err := json.Unmarshal([]byte(data), &rows); err != nil {
			http.Error(w, "Unable to parse from cach ", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Example response
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetSemestersById(ctx, orgid, role)
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
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		data, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseRequestSemesterList
			if err := json.Unmarshal([]byte(data), &rows); err != nil {
				http.Error(w, "Unable to parse from cach ", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// Example response
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

}

func (h *AuthHandler) GetSchedules(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	tutorID := query.Get("tutor_id")
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

	if tutorID == "" {
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Invalid organization_id", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("get:assessments:%d", orgid)
	lkey := h.cacheHander.LockKey(key)
	cdata, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var rows []models.ResponseAssessmentList
		if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
			http.Error(w, "Invalid organization_id", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetAssessmentsById(ctx, orgid, role)
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
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		cdata, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.ResponseAssessmentList
			if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
				http.Error(w, "Invalid organization_id", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	key := fmt.Sprintf("get:subjects:%d", orgid)
	lkey := h.cacheHander.LockKey(key)
	cdata, isHit := h.cacheHander.CheckCache(ctx, key)
	if isHit {
		var rows []models.SubjectList
		if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
			http.Error(w, "Unable to parse cache", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	token := uuid.NewString()
	lock, _ := h.cacheHander.TryAcquireLock(ctx, lkey, token)
	if lock {
		rows, err := h.authService.GetSubjectById(ctx, orgid, role)
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
		h.cacheHander.SetCache(ctx, key, rows)
		h.cacheHander.SafeUnlock(ctx, lkey, token)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": rows}
		json.NewEncoder(w).Encode(response)
		return
	}
	status := h.cacheHander.WaitForCacheUpdate(ctx, key)
	if status {
		cdata, isHit := h.cacheHander.CheckCache(ctx, key)
		if isHit {
			var rows []models.SubjectList
			if err := json.Unmarshal([]byte(cdata), &rows); err != nil {
				http.Error(w, "Unable to parse cache", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := map[string]interface{}{"data": rows}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func (h *AuthHandler) GetSurveys(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:survey")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	org_id, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetSurveysById(ctx, org_id)
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
	var questions []models.SurveyQuestions
	if len(rows) > 0 {
		var ids []int64
		for i := 0; i < len(rows); i++ {
			ids = append(ids, *rows[i].ID)
		}
		q, err := h.authService.GetSurveyQuestions(ctx, ids)
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
		questions = q
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows, "questions": questions}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSurveysByProgram(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	program_id := query.Get("program_id")

	pid, err := strconv.ParseInt(program_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
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
	org_id, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetSurveyProgramsById(ctx, org_id, pid)
	print(rows)
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

func (h *AuthHandler) GetGroupAttendies(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
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
	if query.Get("group_id") == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(query.Get("group_id"), 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetGroupAttendies(ctx, &idd)
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
	query := r.URL.Query()
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	tid, err := helpers.ExtractInt64Claim(claims, "id")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	if query.Get("tutor_id") != "" {
		tutor_id, err := strconv.ParseInt(query.Get("tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		tid = tutor_id
	}
	valid, err := validateRequest(claims, "view:tutor-locations")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	oid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetTutorLocations(ctx, tid, oid)

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
	var model models.SearchQuery
	org_id, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &org_id

	if query.Get("search_term") != "" {
		model.SearchTerm = query.Get("search_term")
	}
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
	org_id, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse organization", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &org_id
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
	if !query.Has("search_term") {
		http.Error(w, "Missing search term", http.StatusBadRequest)
		return
	}
	var model models.SearchQueryTutor

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &orgid
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
	if semester_id == "" {
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
	tid, err := helpers.ExtractInt64Claim(claims, "id")
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

func (h *AuthHandler) GetRecentSessionsTutors(w http.ResponseWriter, r *http.Request) {
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	tutor_id := query.Get("tutor_id")
	tid, err := strconv.ParseInt(tutor_id, 10, 64)
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	rows, err := h.authService.GetRecentSessionsById(ctx, &orgid, &tid)
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

func (h *AuthHandler) GetRecentLocationSessions(w http.ResponseWriter, r *http.Request) {
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	tutor_id := query.Get("location_id")
	lid, err := strconv.ParseInt(tutor_id, 10, 64)
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	semester_id := query.Get("semester_id")
	sid, err := strconv.ParseInt(semester_id, 10, 64)
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	rows, err := h.authService.GetRecentLocationSessions(ctx, &orgid, &lid, &sid)
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

func (h *AuthHandler) GetLocationSessionAverage(w http.ResponseWriter, r *http.Request) {
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
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	tutor_id := query.Get("location_id")
	lid, err := strconv.ParseInt(tutor_id, 10, 64)
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	semester_id := query.Get("semester_id")
	sid, err := strconv.ParseInt(semester_id, 10, 64)
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	rows, err := h.authService.GetLocationSessionAverage(ctx, &orgid, &lid, &sid)
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
	if session_id == "" {
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

func (h *AuthHandler) GetGlobalSchedule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:global-schedule")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	// I want to return all the session and their information
	data, err := h.authService.GetGlobalSchedule(ctx, orgid)
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
	response := map[string]interface{}{"data": data}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	student_id := query.Get("student_id")
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	organization_id, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}

	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if student_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	idd, err := strconv.ParseInt(student_id, 10, 64)
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
	oid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetOrganizationPermissions(ctx, oid)
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
	location_ids := query.Get("location_ids")
	program_ids := query.Get("program_ids")
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

	oid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	tid, err := helpers.ExtractInt64Claim(claims, "id")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	role, err := helpers.ExtractStringClaims(claims, "type")
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
	models.ID = tid
	models.Role = role
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

func (h *AuthHandler) GetAnnouncementsAck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	// This needs to be handled in a PUT request.
	query := r.URL.Query()
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
	if query.Get("announcement_id") == "" {
		http.Error(w, "Missing getter", http.StatusForbidden)
		return
	}
	announcement_id := query.Get("announcement_id")
	if err != nil {
		http.Error(w, "request timeout", http.StatusInternalServerError)
		return
	}
	aid, err := strconv.ParseInt(announcement_id, 10, 64)
	if err != nil {
		http.Error(w, "request timeout", http.StatusInternalServerError)
		return
	}
	oid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "request timeout", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetAnnouncementsAck(ctx, &oid, &aid)
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

// to do
func (h *AuthHandler) CreateStudentSession(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "write:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to get org id", http.StatusBadRequest)
		return
	}
	var req models.RequestAssessmentGrader
	var models models.RegisterStudentSessionList
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	// Send to MQ
	user, err := h.authService.CreateStudentSessions(models)
	if err != nil {
		issue := fmt.Sprintf("Unable to create session: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return
	}
	if models.SessionToken != nil && user.SessionID != nil {
		req.SessionID = user.SessionID
		req.OrganizationID = &orgid
		req.SessionToken = models.SessionToken
		req.TutorID = models.Session.TutorId
		req.SemesterID = models.Session.SemesterId
		payload, err := h.sqsHandler.TagPayloadAssessmentGrader(ctx, "process_assessment_grader", &req)
		if err != nil {
			issue := fmt.Sprintf("unable to add event to mq: %v", err)
			http.Error(w, issue, http.StatusInternalServerError)
			return
		}
		sqs, err := h.sqsHandler.SendMessageToQueue(ctx, h.config.SQS.AssessmentGraderQueue, string(payload))
		if err != nil {
			fmt.Printf("Unable to send message to queue: %v\n", err)
			http.Error(w, "Unable to send message to queue ", http.StatusInternalServerError)
			return
		}
		fmt.Print(sqs.ResultMetadata)

	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) CreateSurveyResponse(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "write:session")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	var models models.RegisterSurvey
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.AddSurveyResponse(ctx, models, &orgid)
	if err != nil {
		issue := fmt.Sprintf("Unable to create session: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return
	}
	ok, err = h.authService.AddSentimentWorker(ctx, models.SessionID, &orgid)
	if err != nil {
		issue := fmt.Sprintf("Unable to create session: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{"status": user, "sentiment_worker": ok}
	json.NewEncoder(w).Encode(resp)
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

func (h *AuthHandler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "write:survey")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	var models models.RegisterRequestSurvey
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	res, err := h.authService.AddSurvey(ctx, models, &orgid)
	if err != nil {
		issue := fmt.Sprintf("Unable to create session: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) UpdateSurvey(w http.ResponseWriter, r *http.Request) {
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
	valid, err := validateRequest(claims, "write:survey")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}

	var models models.RegisterRequestSurvey
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	res, err := h.authService.UpdateSurvey(ctx, models, &orgid)
	if err != nil {
		issue := fmt.Sprintf("Unable to create session: %v", err)
		http.Error(w, issue, http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
