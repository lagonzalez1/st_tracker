package services

import (
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddSubject(req models.RegisterSubject) (*models.ResponseRegisterSubject, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: title or org id")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Subjects(title, description, organization_id)
              VALUES ($1, $2, $3) RETURNING id;`

	err := s.db.QueryRow(query, req.Title, req.Description, *req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseRegisterSubject{
		Status:    "OK",
		SubjectID: newID,
	}, nil
}

func (s *AuthService) UpdateSubject(req models.RegisterSubject) (*models.ResponseUpdate, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: title, org id")
	}
	query := `UPDATE stu_tracker.Subjects SET title = $1, description = $2
              WHERE id = $3`

	_, err := s.db.Exec(query, req.Title, req.Description, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteSubject(req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `DELETE FROM stu_tracker.Subjects WHERE id = $1`

	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}

func (s *AuthService) AddSubjectLocation(req models.RegisterSubjectLocation) (*models.ResponseSubjectLocation, error) {
	// Input validation
	if req.SubjectID == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: title or org id")
	}
	query := `INSERT INTO stu_tracker.Location_subjects(subject_id, location_id)
              VALUES ($1, $2);`

	_, err := s.db.Exec(query, req.SubjectID, req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseSubjectLocation{
		Status: "OK",
		ID:     nil,
	}, nil
}

func (s *AuthService) DeleteSubjectLocation(req models.RemoveSubjectLocation) (*models.RemoveResponse, error) {
	// Input validation
	if req.SubjectID == nil || req.LocationID == nil {
		return nil, fmt.Errorf("missing required fields: title or org id")
	}
	query := `DELETE FROM stu_tracker.Location_subjects WHERE subject_id = $1 AND location_id = $2;`

	_, err := s.db.Exec(query, req.SubjectID, req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
