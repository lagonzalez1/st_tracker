package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
)

func (h *AuthHandler) CreateAssessmentSessions(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	var models models.RegisterStudentAssessmentSession
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.CreateAssessmentSession(ctx, models)
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
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetStudentAssessmentSessions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	session_id := query.Get("session_id")
	tutor_id := query.Get("tutor_id")
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

	idd, err := strconv.ParseInt(tutor_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse tutor_id", http.StatusInternalServerError)
		return
	}
	rows, err := h.authService.GetAssessmentBySessionId(ctx, session_id, &idd)
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

func (h *AuthHandler) DeleteAssessmentSessions(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	var models models.DeleteAssessmentSession
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteAssessmentSession(ctx, models)
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
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteStudentSession(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
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
	var models models.DeleteStudentSession
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.DeleteStudentSession(ctx, models)
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
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetAssessmentQuestionsExternal(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	assessmentID := query.Get("assessment_id")

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

func (h *AuthHandler) GetStudentAssessmentChoices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	assessment_id := query.Get("assessment_id")
	student_id := query.Get("student_id")
	session_id := query.Get("session_id")

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

	if assessment_id == "" || student_id == "" || session_id == "" {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	sid, err := strconv.ParseInt(student_id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid student id", http.StatusBadRequest)
		return
	}
	rows, err := h.authService.GetAssessmentChoicesByStudent(ctx, &sid, &session_id)
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

func (h *AuthHandler) CreateStudentAssessmentResponse(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	var models models.RegisterStudentAssessment
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
	}
	user, err := h.authService.CreateStudentAssessmentResponse(ctx, models)
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
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
