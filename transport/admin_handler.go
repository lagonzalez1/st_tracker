package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"tracker/app/helpers"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
)

func (h *AuthHandler) GetSessionAnalytics(w http.ResponseWriter, r *http.Request) {
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
	if !query.Has("id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestSessionBChart
	org_id, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &org_id

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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
	user, err := h.authService.GetSessionAnalytics(model)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetSessionAnalyticsLocal(w http.ResponseWriter, r *http.Request) {
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
	if !query.Has("id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestSessionBChart
	org_id, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &org_id

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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
	user, err := h.authService.GetSessionsAnalyticsLocal(ctx, model)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetTutorSessionAnalytics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// Undefined variables like optional location_id
	// Check if exist before
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
	if !query.Has("email") || !query.Has("role") || !query.Has("id") || !query.Has("organization_id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestTutorSessions
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}

	user, err := h.authService.GetTutorSessions(model)
	if err != nil {
		http.Error(w, "Unable to create Admin staff", http.StatusInternalServerError)
		fmt.Printf("Unable to insert Admin staff: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"data": user}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionBChart(w http.ResponseWriter, r *http.Request) {
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
	var model models.RequestSessionBChart
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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

	rows, err := h.authService.GetSessionBChart(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentGroups(w http.ResponseWriter, r *http.Request) {
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
	// Undefined variables like optional location_id
	// Check if exist before
	if !query.Has("location_id") || !query.Has("tutor_id") || !query.Has("semester_id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	lid, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	tid, err := helpers.ExtractInt64Claim(claims, "id")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	sid, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	rows, err := h.authService.GetStudentGroups(ctx, &lid, &tid, &sid)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetCycleGrowth(w http.ResponseWriter, r *http.Request) {
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
	if !query.Has("email") || !query.Has("role") || !query.Has("id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestCycleGrowth
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to find orgid", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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

	rows, err := h.authService.GetCycleGrowth(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAbsentPresent(w http.ResponseWriter, r *http.Request) {
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

	var model models.RequestCycleGrowth
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to find orgid", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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

	rows, err := h.authService.GetAbsentPresent(ctx, model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentCompletion(w http.ResponseWriter, r *http.Request) {
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

	var model models.RequestCycleGrowth
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to find orgid", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}

	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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
	// Safe response
	if model.SemesterID == nil || query.Get("semester_id") == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Example response
		response := map[string]interface{}{"data": nil}
		json.NewEncoder(w).Encode(response)
		return
	}

	rows, err := h.authService.GetAssessmentCompletion(ctx, model)
	if err != nil {
		fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetCycleGrowthDelim(w http.ResponseWriter, r *http.Request) {
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
	if !query.Has("email") || !query.Has("role") || !query.Has("id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestCycleGrowth
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to find orgid", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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

	rows, err := h.authService.GetCycleGrowthDelim(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	query := r.URL.Query()
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
	// Undefined variables like optional location_id

	var keyPath string
	var signedUrl *string
	if query.Has("key") {
		keyPath = query.Get("key")
	}
	if keyPath != "" {
		url, err := h.authService.GeneratePresignedUrl(ctx, keyPath)
		if err != nil {
			http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionBChart %v", err)
			return
		}
		signedUrl = &url
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"url": signedUrl}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentBChart(w http.ResponseWriter, r *http.Request) {
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

	var model models.RequestSessionBChart
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "unable to get organization", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetAssessmentBChart %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetAssessmentBChart %v", err)
			return
		}
		model.EndDate = end_date
	}

	rows, err := h.authService.GetAssessmentBChart(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetAssessmentBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetProgramsBChart(w http.ResponseWriter, r *http.Request) {
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
	var model models.RequestSessionBChart
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetProgramsBChart %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetProgramsBChart %v", err)
			return
		}
		model.EndDate = end_date
	}

	rows, err := h.authService.GetProgramsBChart(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetProgramsBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorsBChart(w http.ResponseWriter, r *http.Request) {
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
	if !query.Has("organization_id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestSessionBChart
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.EndDate = end_date
	}

	rows, err := h.authService.GetTutorsBChart(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionVScore(w http.ResponseWriter, r *http.Request) {
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
	var model models.RequestSessionBChart
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.EndDate = end_date
	}

	rows, err := h.authService.GetSessionVScore(ctx, model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetStudentVAssessments(w http.ResponseWriter, r *http.Request) {
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
	if !query.Has("organization_id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestSessionBChart
	model.OrganizationID = &orgid

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.EndDate = end_date
	}

	rows, err := h.authService.GetStudentVAssessments(ctx, model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentTrendLine(w http.ResponseWriter, r *http.Request) {
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
	var model models.RequestSessionBChart
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetAssessmentTrendLine %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetAssessmentTrendLine %v", err)
			return
		}
		model.EndDate = end_date
	}

	rows, err := h.authService.GetAssessmentTrendLine(ctx, model)
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
		http.Error(w, "Unable to get assessments trend line", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSessionTrendLine(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:session")
	if err != nil || !valid {
		fmt.Println(err)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	// Undefined variables like optional location_id
	// Check if exist before
	if !query.Has("email") || !query.Has("role") || !query.Has("id") || !query.Has("organization_id") {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}
	var model models.RequestSessionBChart
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("start_date") != "" {
		start_time, err := time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionTrendLine %v", err)
			return
		}
		model.StartDate = start_time
	}
	if query.Get("end_date") != "" {
		end_date, err := time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetSessionTrendLine %v", err)
			return
		}
		model.EndDate = end_date
	}

	rows, err := h.authService.GetSessionTrendLine(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionTrendLine %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSemestersVAssessmentChart(w http.ResponseWriter, r *http.Request) {
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
	var model models.RequestSemestersVAssessmentChart
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}
	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester1_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester1_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.Semester1ID = &sem_id
	}
	if query.Get("semester2_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester2_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.Semester2ID = &sem_id
	}
	if query.Get("assessment1_id") != "" {
		assessment_id, err := strconv.ParseInt(query.Get("assessment1_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.Assessment1ID = &assessment_id
	}
	rows, err := h.authService.GetSemestersVAssessmentChart(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionTrendLine %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentGrowth(w http.ResponseWriter, r *http.Request) {
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
	var model models.RequestAssessmentGrowth
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}
	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester1_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester1_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.Semester1ID = &sem_id
	}
	if query.Get("assessment1_id") != "" {
		assessment_id, err := strconv.ParseInt(query.Get("assessment1_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.Assessment1ID = &assessment_id
	}
	if query.Get("assessment2_id") != "" {
		assessment_id, err := strconv.ParseInt(query.Get("assessment2_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.Assessment2ID = &assessment_id
	}

	rows, err := h.authService.GetAssessmentGrowth(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetSessionTrendLine %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Example response
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetTutorFile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	var model *models.RequestDownloadData
	model = &models.RequestDownloadData{}
	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("sort_key") != "" {
		model.SortKey = query.Get("sort_key")
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("program_id") != "" {
		program_id, err := strconv.ParseInt(query.Get("program_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.ProgramID = &program_id
	}
	if query.Get("date") != "" {
		DateStart, err := time.Parse("2006-01-02", query.Get("date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("unable to parse start time %v", err)
			return
		}
		model.DateStart = DateStart
	}
	if query.Get("end_date") != "" {
		DateEnd, err := time.Parse("2006-01-02", query.Get("date_end"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("unable to parse start timet %v", err)
			return
		}
		model.DateEnd = DateEnd
	}
	if query.Get("subject_id") != "" {
		subject_id, err := strconv.ParseInt(query.Get("subject_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SubjectID = &subject_id
	}

	var tutor = "tutor"
	model.Entity = &tutor
	inputKey, err := h.authService.ProcessDownloadEvent(ctx, model, &orgID)
	if err != nil {
		http.Error(w, "unable to add message to MQ", http.StatusInternalServerError)
		return
	}
	payload, err := h.sqsHandler.TagPayloadTutorDownload(ctx, "fetch_tutor_data", model)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to tag request ", http.StatusInternalServerError)
		return
	}

	sqs, err := h.sqsHandler.SendMessageToQueue(ctx, h.config.SQS.DataReportsQueue, string(payload))
	if err != nil {
		fmt.Printf("Unable to send message to queue: %v\n", err)
		http.Error(w, "Unable to send message to queue ", http.StatusInternalServerError)
		return
	}
	fmt.Print(sqs.ResultMetadata)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inputKey)

}

func (h *AuthHandler) GetStudentFile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	var model *models.RequestDownloadData
	model = &models.RequestDownloadData{}
	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}

	if query.Get("subject_id") != "" {
		subject_id, err := strconv.ParseInt(query.Get("subject_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SubjectID = &subject_id
	}
	if query.Get("sort_key") != "" {
		model.SortKey = query.Get("sort_key")
	}

	if query.Get("program_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("program_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
	if query.Get("date") != "" {
		DateStart, err := time.Parse("2006-01-02", query.Get("date"))
		if err != nil {
			http.Error(w, "unable to parse start time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.DateStart = DateStart
	}
	if query.Get("end_date") != "" {
		DateEnd, err := time.Parse("2006-01-02", query.Get("date_end"))
		if err != nil {
			http.Error(w, "unable to parse end time", http.StatusInternalServerError)
			fmt.Printf("Unable to get rows in GetTutorsBChart %v", err)
			return
		}
		model.DateEnd = DateEnd
	}
	if query.Get("data_type") != "" {
		stringPtr := query.Get("data_type")
		model.DataType = &stringPtr
	}
	var student = "student"
	model.Entity = &student
	inputKey, err := h.authService.ProcessDownloadEvent(ctx, model, &orgID)
	if err != nil {
		http.Error(w, "unable to add message to MQ", http.StatusInternalServerError)
		return
	}
	payload, err := h.sqsHandler.TagPayloadTutorDownload(ctx, "fetch_student_data", model)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to tag request ", http.StatusInternalServerError)
		return
	}

	sqs, err := h.sqsHandler.SendMessageToQueue(ctx, h.config.SQS.DataReportsQueue, string(payload))
	if err != nil {
		fmt.Printf("Unable to send message to queue: %v\n", err)
		http.Error(w, "Unable to send message to queue ", http.StatusInternalServerError)
		return
	}
	fmt.Print(sqs.ResultMetadata)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inputKey)

}

func (h *AuthHandler) GetTutorLowPerformance(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	var model models.RequestSessionBChart
	if query.Get("organization_id") != "" {
		org_id, err := strconv.ParseInt(query.Get("organization_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.OrganizationID = &org_id
	}

	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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
	user, err := h.authService.GetTutorLowPerformance(ctx, model)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Unable to get GetTutorLowPerformance", http.StatusInternalServerError)
		fmt.Printf("Unable to get GetTutorLowPerformance: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetSentiment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:sessions")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var model models.RequestSentiment
	if query.Get("location_id") != "" {
		loc_id, err := strconv.ParseInt(query.Get("location_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.LocationID = &loc_id
	}
	if query.Get("semester_id") != "" {
		sem_id, err := strconv.ParseInt(query.Get("semester_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.SemesterID = &sem_id
	}
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

	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	model.OrganizationID = &orgid

	rows, err := h.authService.GetSentimentAnalysis(ctx, model)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSentimentByTutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:sessions")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var tid *int64
	if query.Get("tutor_id") != "" {
		t_id, err := strconv.ParseInt(query.Get("tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		tid = &t_id
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetSentimentAnalysisByTutor(ctx, tid, &orgid)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentByTutor(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:sessions")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var tid *int64
	if query.Get("tutor_id") != "" {
		t_id, err := strconv.ParseInt(query.Get("tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		tid = &t_id
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetAssessmentsByTutor(ctx, tid, &orgid)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAssessmentContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:sessions")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var sid *int64
	if query.Get("session_id") != "" {
		s_id, err := strconv.ParseInt(query.Get("session_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		sid = &s_id
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetAssessmentPivotTable(ctx, sid, &orgid)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:*")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetSubscriptions(ctx, &orgid)
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

	ent, err := h.authService.GetSubscriptionsEntitlements(ctx, &orgid)
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

	user, err := h.authService.GetSubscriptionsByOrganization(ctx, &orgid)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"sub": rows, "ent": ent, "user": user}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetAdminLocations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:*")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	if query.Get("admin_id") == "" {
		http.Error(w, "Invalid request, missing params", http.StatusInternalServerError)
		return
	}
	admin_id := query.Get("admin_id")
	aid, err := strconv.ParseInt(admin_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	rows, err := h.authService.GetAdminLocations(ctx, &aid, &orgid)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": rows}
	json.NewEncoder(w).Encode(response)
}
