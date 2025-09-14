package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"tracker/app/helpers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xuri/excelize/v2"
)

func (h *AuthHandler) UploadTutorBigData(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

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

	oid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil || oid <= 0 {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}

	// Save the uploaded file locally (optional)
	tempFile, err := os.CreateTemp("", "uploaded-*.xlsx")
	if err != nil {
		http.Error(w, "Error creating temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	// Process Excel file
	excelFile, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		http.Error(w, "Error opening Excel file", http.StatusInternalServerError)
		return
	}
	defer excelFile.Close()
	rows, err := excelFile.GetRows("Sheet1")
	if err != nil {
		http.Error(w, "Error reading sheet", http.StatusInternalServerError)
		return
	}
	returnFile, response, err := h.authService.RegisterMultipleTutors(ctx, rows, &oid)
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
	w.Header().Set("Content-Disposition", "attachment; filename=data.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	str := strconv.Itoa(response.Count)
	// Include metadata in headers
	w.Header().Set("X-Response-Message", response.Status)
	w.Header().Set("X-Response-Status", str)
	// Write the Excel file to the response
	returnFile.WriteToBuffer()
	if err := returnFile.Write(w); err != nil {
		http.Error(w, "Failed to write Excel file", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) UploadStudentBigData(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

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

	oid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil || oid <= 0 {
		http.Error(w, "Unable to parse claims query", http.StatusBadRequest)
		return
	}

	location_id := r.FormValue("location_id")
	lid, err := strconv.ParseInt(location_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse lid", http.StatusInternalServerError)
		return
	}

	semester_id := r.FormValue("semester_id")
	sid, err := strconv.ParseInt(semester_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse sid", http.StatusInternalServerError)
		return
	}
	// Save the uploaded file locally (optional)
	tempFile, err := os.CreateTemp("", "uploaded-*.xlsx")
	if err != nil {
		http.Error(w, "Error creating temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error saving uploaded file", http.StatusInternalServerError)
		return
	}
	excelFile, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error opening Excel file", http.StatusInternalServerError)
		return
	}
	defer excelFile.Close()
	rows, err := excelFile.GetRows("Sheet1")
	if err != nil {
		http.Error(w, "Error reading sheet", http.StatusInternalServerError)
		return
	}
	response, err := h.authService.RegisterMultipleStudents(ctx, rows, &oid, sid, lid)
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
	json.NewEncoder(w).Encode(response)
}
