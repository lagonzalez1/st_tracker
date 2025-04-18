package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func (h *AuthHandler) UploadTutorBigData(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	organization_id := r.FormValue("organization_id")
	if err != nil {
		http.Error(w, "Error retrieving organization_id", http.StatusBadRequest)
		return
	}
	oid, err := strconv.ParseInt(organization_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
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
	returnFile, response, err := h.authService.RegisterMultipleTutors(rows, &oid)
	if err != nil {
		http.Error(w, "Error return file RegisterMultipleTutors", http.StatusInternalServerError)
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
	r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	organization_id := r.FormValue("organization_id")
	if err != nil {
		http.Error(w, "Error retrieving organization_id", http.StatusBadRequest)
		return
	}
	oid, err := strconv.ParseInt(organization_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse oid", http.StatusInternalServerError)
		return
	}

	location_id := r.FormValue("location_id")
	if err != nil {
		http.Error(w, "Error retrieving organization_id", http.StatusBadRequest)
		return
	}
	lid, err := strconv.ParseInt(location_id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse lid", http.StatusInternalServerError)
		return
	}

	semester_id := r.FormValue("semester_id")
	if err != nil {
		http.Error(w, "Error retrieving organization_id", http.StatusBadRequest)
		return
	}
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
	/*
		Dont think this is needed, just process the file and move along
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
	*/
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
	response, err := h.authService.RegisterMultipleStudents(rows, &oid, sid, lid)
	if err != nil {
		fmt.Printf("Error on RegisterMultipleStudents %v", err)
		http.Error(w, "Error RegisterMultipleTutors", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
