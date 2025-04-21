package services

import (
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddProgram(req models.RegisterRequestProgram) (*models.ResponseRequestProgram, error) {
	// Input validation
	if req.ProgramName == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Programs (program_name, organization_id) 
			  VALUES ($1, $2) RETURNING id;`

	err := s.db.QueryRow(query, req.ProgramName, *req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.ResponseRequestProgram{
		Status:    "OK",
		ProgramId: newID,
	}, nil
}

func (s *AuthService) UpdateProgram(req models.RegisterRequestProgram) (*models.ResponseRequestProgram, error) {
	// Input validation
	if req.ProgramName == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `UPDATE stu_tracker.Programs SET program_name = $1 WHERE id = $2`
	_, err := s.db.Exec(query, req.ProgramName, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to update program: %w", err)
	}
	return &models.ResponseRequestProgram{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteProgram(req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `DELETE FROM stu_tracker.Programs WHERE id = $1;`
	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete program: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}

func (s *AuthService) CreateProgramLocation(req models.RegisterLocationProgram) (*models.RegisterResponseLocationProgram, error) {
	// Input validation
	if req.LocationId == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `INSERT INTO stu_tracker.Location_programs (location_id, program_id, organization_id) 
			  VALUES ($1, $2, $3);`

	_, err := s.db.Exec(query, req.LocationId, req.ProgramId, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RegisterResponseLocationProgram{
		Status: "OK",
	}, nil
}

func (s *AuthService) DeleteProgramLocation(req models.RemoveLocationProgram) (*models.RemoveResponse, error) {
	// Input validation
	if req.LocationId == nil || req.ProgramId == nil {
		return nil, fmt.Errorf("missing required fields: ID")
	}
	query := `DELETE FROM stu_tracker.Location_programs WHERE location_id = $1 AND program_id = $2`

	_, err := s.db.Exec(query, req.LocationId, req.ProgramId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RemoveResponse{
		Status: "OK",
	}, nil
}
