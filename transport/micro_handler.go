package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"tracker/app/helpers"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (h *AuthHandler) MicroEventStartStudentReport(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:student-data")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil || orgID <= 0 {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	var model models.RequestStudentReport
	if err := json.Unmarshal(body, &model); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding json %v", err)
		return
	}
	outputKey := uuid.New().String()
	model.S3OutputKey = &outputKey
	//id, err := h.authService.AddStudentReportQuery(ctx, model) Note replaced
	id, err := h.authService.ProcessStudentReport(ctx, model)
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
		http.Error(w, "Unable to start ", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	// Update the db, tag the payload
	payload, err := h.sqsHandler.TagPayloadStudentReport(ctx, "fetch_student_data", &model)
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
	json.NewEncoder(w).Encode(id)

}

func (h *AuthHandler) MicroGetStudentReport(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:student-data")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	inputKey := query.Get("input_key")
	if inputKey == "" {
		http.Error(w, "missing paramaters", http.StatusInternalServerError)
		return
	}
	status, outputKey, err := h.authService.GetStudentReportStatus(ctx, &inputKey)
	if err != nil {
		http.Error(w, "error on GetStudentReportStatus", http.StatusInternalServerError)
		return
	}
	var report *string
	if *status == "DONE" {
		var key = "student_reports/" + *outputKey
		report, err = h.authService.GetS3Object(ctx, key, "tracker-student-reports")
		if err != nil {
			http.Error(w, "Completed task, but unable to get s3 object", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]interface{}{"status": status, "report": report}
	json.NewEncoder(w).Encode(payload)
}

func (h *AuthHandler) MicroEventGenerate(w http.ResponseWriter, r *http.Request) {
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
	orgID, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}
	outputKey := uuid.New()
	var models models.RequestEventGeneration
	key := outputKey.String()
	if err := json.Unmarshal(body, &models); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}
	gt := strings.ToLower(strings.TrimSpace(*models.GenerationType))

	if models.GenerationType == nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v", err)
		return
	}

	switch gt {
	case "generate_questions":
		models.OrganizationID = &orgID
		models.RequestQuestions.S3OutputKey = &key
		id, err := h.authService.CreateAssessmentGenerationTask(ctx, &models)
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
			http.Error(w, "Unable to start CreateAssessmentGenerationTask ", http.StatusInternalServerError)
			fmt.Printf("service error: %v\n", err)
			return
		}
		payload, err := h.sqsHandler.TagPayloadAssessmentGenerator(ctx, "process_assessment_generation", &models)
		if err != nil {
			issue := fmt.Sprintf("unable to tag request: %v", err)
			http.Error(w, issue, http.StatusInternalServerError)
			return
		}
		sqs, err := h.sqsHandler.SendMessageToFIFOQueue(ctx, h.config.SQS.GenerateContentQueue, string(payload), orgID)
		if err != nil {
			fmt.Printf("Unable to send message to queue: %v\n", err)
			http.Error(w, "Unable to send message to queue ", http.StatusInternalServerError)
			return
		}
		fmt.Print(sqs.ResultMetadata)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(id)
	case "generate_materials":
		models.OrganizationID = &orgID
		models.RequestMaterials.S3OutputKey = &key

		id, err := h.authService.CreateMaterialsGenerationTask(ctx, &models)
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
			http.Error(w, "Unable to start CreateMaterialsGenerationTask ", http.StatusInternalServerError)
			fmt.Printf("service error: %v\n", err)
			return
		}
		payload, err := h.sqsHandler.TagPayloadAssessmentGenerator(ctx, "process_materials_generation", &models)
		if err != nil {
			issue := fmt.Sprintf("unable to tag request: %v", err)
			http.Error(w, issue, http.StatusInternalServerError)
			return
		}
		sqs, err := h.sqsHandler.SendMessageToFIFOQueue(ctx, h.config.SQS.GenerateContentQueue, string(payload), orgID)
		if err != nil {
			fmt.Printf("Unable to send message to queue: %v\n", err)
			http.Error(w, "Unable to send message to queue ", http.StatusInternalServerError)
			return
		}
		fmt.Print(sqs.ResultMetadata)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(id)
	default:
		http.Error(w, "unsupported generation_type", http.StatusBadRequest)
		return
	}

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
	valid, err := validateRequest(claims, "view:assessments")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	inputKey := query.Get("input_key")
	if inputKey == "" {
		http.Error(w, "missing paramaters", http.StatusInternalServerError)
		return
	}
	status, jsonOutput, err := h.authService.GetQuestionGenerationStatus(ctx, &inputKey)
	if err != nil {
		http.Error(w, "error on GetQuestionGenerationStatus", http.StatusInternalServerError)
		return
	}
	var response *string
	var assessment *models.Assessment
	if *status == "COMPLETE" || *status == "DONE" {
		err := json.Unmarshal(jsonOutput, &assessment)
		if err != nil {
			http.Error(w, "unable to parse json into assessment", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]interface{}{"status": status, "response": response, "assessment": assessment}
	json.NewEncoder(w).Encode(payload)
}

func (h *AuthHandler) GetGeneratedMaterials(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()

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
	inputKey := query.Get("input_key")
	if inputKey == "" {
		http.Error(w, "missing paramaters", http.StatusInternalServerError)
		return
	}
	status, outputKey, jsonOutput, err := h.authService.GetMaterialsGenerationStatus(ctx, &inputKey)
	if err != nil {
		http.Error(w, "error on GetMaterialsGenerationStatus", http.StatusInternalServerError)
		return
	}
	var response *string
	var materials *models.Materials
	var signedUrl *string
	if *status == "DONE" || *status == "COMPLETE" {
		err := json.Unmarshal(jsonOutput, &materials)
		if err != nil {
			http.Error(w, "unable to parse json into assessment", http.StatusInternalServerError)
			return
		}
		path := fmt.Sprintf("materials/%s.pdf", *outputKey)
		url, err := h.authService.GeneratePutPresignedUrlMaterials(ctx, &path, "application/pdf", 5)
		if err != nil {
			http.Error(w, "unable to create signed url for upload", http.StatusInternalServerError)
		}
		signedUrl = url
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]interface{}{"status": status, "response": response, "output_key": *outputKey, "materials": materials, "signed_url": signedUrl}
	json.NewEncoder(w).Encode(payload)
}

func (h *AuthHandler) MicroGetTutorFile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:tutor-data")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	inputKey := query.Get("input_key")
	if inputKey == "" {
		http.Error(w, "missing paramaters", http.StatusInternalServerError)
		return
	}
	status, url, err := h.authService.GetOrganizationReportStatus(ctx, &inputKey)
	if err != nil {
		http.Error(w, "error on GetTutorFileStatus", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]interface{}{"status": status, "url": url}
	json.NewEncoder(w).Encode(payload)
}

func (h *AuthHandler) MicroGetStudentFile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	query := r.URL.Query()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "view:student-data")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	inputKey := query.Get("input_key")
	if inputKey == "" {
		http.Error(w, "missing paramaters", http.StatusInternalServerError)
		return
	}
	status, url, err := h.authService.GetOrganizationReportStatus(ctx, &inputKey)
	if err != nil {
		http.Error(w, "error on GetTutorFileStatus", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]interface{}{"status": status, "url": url}
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
		http.Error(w, "Unable to delete generated assessment. ", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
