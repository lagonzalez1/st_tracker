package services

import (
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddSemester(req models.RegisterRequestSemester) (*models.ResponseRequestSemester, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Semester (title, year, organization_id) 
			  VALUES ($1, $2, $3) RETURNING id;`

	err := s.db.QueryRow(query, req.Title, req.Year, *req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to add semester: %w", err)
	}
	return &models.ResponseRequestSemester{
		Status: "OK",
		ID:     &newID,
	}, nil
}

func (s *AuthService) UpdateSemester(req models.RegisterRequestSemester) (*models.ResponseUpdate, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `UPDATE stu_tracker.Semester SET title = $1, year = $2 WHERE id = $3;`

	_, err := s.db.Exec(query, req.Title, req.Year, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update program: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteSemester(req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `DELETE FROM stu_tracker.Semester WHERE id = $1;`

	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update program: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
