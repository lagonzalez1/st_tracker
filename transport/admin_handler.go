package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"tracker/app/models"
)

func (h *AuthHandler) GetSessionAnalytics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

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
	query := r.URL.Query()

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

	rows, err := h.authService.GetAssessmentTrendLine(model)
	if err != nil {
		http.Error(w, "Unable to retrive rows given id", http.StatusInternalServerError)
		fmt.Printf("Unable to get rows in GetAssessmentTrendLine %v", err)
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
	if query.Get("program_id") != "" {
		program_id, err := strconv.ParseInt(query.Get("program_id"), 10, 64)
		if err != nil {
			http.Error(w, "Unable to parse id", http.StatusInternalServerError)
			return
		}
		model.ProgramID = &program_id
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
