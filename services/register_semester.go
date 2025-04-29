package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddSemester(c context.Context, req models.RegisterRequestSemester) (*models.ResponseRequestSemester, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Semester (title, year, organization_id, date_start, date_end, active) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	err := s.db.QueryRowContext(c, query, req.Title, req.Year, *req.OrganizationId, req.DateStart, req.DateEnd, req.Active).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to add semester: %w", err)
	}
	return &models.ResponseRequestSemester{
		Status: "OK",
		ID:     &newID,
	}, nil
}

func (s *AuthService) AddSemesterLocation(c context.Context, req models.RegisterRequestSemesterLocation) (*models.ResponseRequestSemester, error) {
	// Input validation
	if req.LocationID == nil || req.OrganizationId == nil || req.SemesterID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Semester_Location (semester_id, location_id, organization_id) 
			  VALUES ($1, $2, $3) ON CONFLICT(semester_id, location_id, organization_id) DO NOTHING RETURNING id;`

	err := s.db.QueryRowContext(c, query, req.SemesterID, req.LocationID, req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to add semester: %w", err)
	}
	return &models.ResponseRequestSemester{
		Status: "OK",
		ID:     &newID,
	}, nil
}

func (s *AuthService) UpdateSemesterLocation(c context.Context, req models.RegisterRequestSemesterLocation) (*models.ResponseRequestSemester, error) {
	// Input validation
	if req.LocationID == nil || req.OrganizationId == nil || req.SemesterID == nil || req.ID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	var newID int64
	query := `UPDATE stu_tracker.Semester_Location SET semester_id = $1, location_iD = $2 WHERE id = $3`

	err := s.db.QueryRowContext(c, query, req.SemesterID, req.LocationID, req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to add semester: %w", err)
	}
	return &models.ResponseRequestSemester{
		Status: "OK",
		ID:     &newID,
	}, nil
}

func (s *AuthService) DeleteSemesterLocation(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	var newID int64
	query := `DELETE FROM stu_tracker.Semester_Location WHERE id = $1`

	err := s.db.QueryRowContext(c, query, req.ID).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to add semester: %w", err)
	}
	return &models.RemoveResponse{
		Status: "OK",
	}, nil
}

func (s *AuthService) UpdateSemester(c context.Context, req models.RegisterRequestSemester) (*models.ResponseUpdate, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `UPDATE stu_tracker.Semester SET title = $1, year = $2, date_start = $3, date_end = $4, active = $5 WHERE id = $6;`
	_, err := s.db.ExecContext(c, query, req.Title, req.Year, req.DateStart, req.DateEnd, req.Active, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update program: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteSemester(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `DELETE FROM stu_tracker.Semester WHERE id = $1;`

	_, err := s.db.ExecContext(c, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update program: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
