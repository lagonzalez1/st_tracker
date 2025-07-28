package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"tracker/app/helpers"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (h *AuthHandler) MicroEventGenQuestions(w http.ResponseWriter, r *http.Request) {
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	orgID, err := helpers.ExtractFloat64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	outputKey := uuid.New()
	var models models.RequestQuestions
	key := outputKey.String()
	models.S3OutputKey = &key

	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	// Check if claims is the same as input data
	if int64(orgID) != *models.OrganizationID {
		http.Error(w, "Invalid claims and input missmatch", http.StatusBadRequest)
		fmt.Printf("Claims is different from input data: %v", err)
	}
	// Need to handle 3 cases of logins for different permissions
	id, err := h.authService.AddQueueQuestionEvent(ctx, models)
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
		http.Error(w, "Unable to start AddQueueQuestionEvent ", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func (h *AuthHandler) GetGeneratedQuestion(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()

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
	inputKey := query.Get("input_key")
	if inputKey == "" {
		http.Error(w, "missing paramaters", http.StatusInternalServerError)
		return
	}
	status, outputKey, err := h.authService.GetQuestionGenerationStatus(ctx, &inputKey)
	if err != nil {
		http.Error(w, "error on GetQuestionGenerationStatus", http.StatusInternalServerError)
		return
	}
	var response *string
	if *status == "COMPLETE" {
		response, err = h.authService.GetS3Object(ctx, *outputKey)
		if err != nil {
			http.Error(w, "unable to get s3 object", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]interface{}{"status": status, "response": response}
	json.NewEncoder(w).Encode(payload)
}

func (h *AuthHandler) MicroEventDeleteQuestions(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("error found reading body")
		return
	}
	orgID, err := helpers.ExtractFloat64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	var model models.RemoveGeneratedQuestion
	if err := json.Unmarshal(body, &model); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	model.OrganizationID = int64(orgID)

	status, err := h.authService.DeleteGeneratedAssessment(ctx, model)
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
		http.Error(w, "Unable to start AddQueueQuestionEvent ", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
