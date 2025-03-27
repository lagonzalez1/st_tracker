package services

import (
	"fmt"
	"log"
	"tracker/app/models"
)

/**
 Admin types:

 For Placment info regarding substitutes and logistics
 View: locations, materail, district, programs, students, subjects, tutors
 Delete: Tutor-locations

**/

func (s *AuthService) AddAnnouncement(req models.RegisterAnnouncements) (*models.ResponseAnnouncement, error) {
	// Input validation
	if req.Title == "" || req.Body == "" || req.OrganizationID == 0 {
		return nil, fmt.Errorf("missing required fields: title, body, admin_id, or organization_id")
	}
	if req.AdminID == nil && req.StaffID == nil {
		return nil, fmt.Errorf("staff and admin id null")

	}
	// Set default severity if not provided
	if req.Severity == "" {
		req.Severity = "info"
	}
	// Insert the announcement into the database
	values := []interface{}{}
	query := `INSERT INTO stu_tracker.Announcements (
			title, body, location_id, severity, organization_id, program_id, admin_id, staff_id) VALUES `
	placeHolderIdx := 1

	if len(req.LocationID) > 0 {
		for i := 0; i < len(req.LocationID); i++ {
			if i > 0 {
				query += `,`
			}
			query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				placeHolderIdx, placeHolderIdx+1, placeHolderIdx+2, placeHolderIdx+3,
				placeHolderIdx+4, placeHolderIdx+5, placeHolderIdx+6, placeHolderIdx+7)
			values = append(values, req.Title, req.Body, req.LocationID[i], req.Severity, req.OrganizationID, req.ProgramID, req.AdminID, req.StaffID)
			placeHolderIdx += 8
		}
	} else {
		query += `($1, $2, $3, $4, $5, $6, $7, $8)`
		values = append(values, req.Title, req.Body, nil, req.Severity, req.OrganizationID, req.ProgramID, req.AdminID, req.StaffID)
	}

	result, err := s.db.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to session students query: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Rows affected error: %v", err)
	}

	// Return the response
	return &models.ResponseAnnouncement{
		Status: "OK",
		ID:     rowsAffected,
	}, nil
}

func (s *AuthService) UpdateAnnouncement(req models.RegisterUpdateAnnouncements) (*models.ResponseAnnouncement, error) {
	if req.ID == nil || req.Title == "" || req.Body == "" || req.OrganizationID == 0 {
		return nil, fmt.Errorf("missing required fields: id, title, body, admin_id, or organization_id")
	}
	if req.AdminID == nil && req.StaffID == nil {
		return nil, fmt.Errorf("admin id and staff id null")
	}
	// Set default severity if not provided
	if req.Severity == "" {
		req.Severity = "info"
	}
	var query string
	var rowsAffected int64
	query += `
		UPDATE stu_tracker.Announcements
		SET 
			title = $1,
			body = $2,
			location_id = $3,
			severity = $4,
			organization_id = $5,
			program_id = $6,
			admin_id = $7,
			staff_id = $8
		WHERE id = $9;`

	result, err := s.db.Exec(
		query,
		req.Title,
		req.Body,
		req.LocationID,
		req.Severity,
		req.OrganizationID,
		req.ProgramID,
		req.AdminID,
		req.StaffID,
		req.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update announcement: %w", err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		log.Fatalf("Rows affected error: %v", err)
	}

	// Return the response
	return &models.ResponseAnnouncement{
		Status: "Updated",
		ID:     rowsAffected,
	}, nil
}

func (s *AuthService) DeleteAnnouncement(req models.RemoveRequest) (*models.RemoveResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `DELETE FROM stu_tracker.Announcements WHERE id = $1;`
	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete staff: %w", err)

	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
