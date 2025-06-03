package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
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
	var model models.RequestDownloadData
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

	file, err := h.authService.GetTutorFileData(ctx, model)
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

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=tutor_data.xlsx")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	// Write the Excel file to the response
	file.WriteToBuffer() // This can be an issue if the file is to large

	if err := file.Write(w); err != nil {
		http.Error(w, "Failed to write Excel file", http.StatusInternalServerError)
		return
	}
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
	var model models.RequestDownloadData
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

	file, err := h.authService.GetStudentFileData(ctx, model)
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

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=student_data.xlsx")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	// Write the Excel file to the response
	file.WriteToBuffer() // This can be an issue if the file is to large

	if err := file.Write(w); err != nil {
		http.Error(w, "Failed to write Excel file", http.StatusInternalServerError)
		return
	}
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
