package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddProgram(c context.Context, req models.RegisterRequestProgram) (*models.ResponseRequestProgram, error) {
	// Input validation
	if req.ProgramName == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Programs (program_name, organization_id, timeframe_required, survey_required) 
			  VALUES ($1, $2, $3, $4) RETURNING id;`

	err := s.db.QueryRowContext(c, query, req.ProgramName, *req.OrganizationId, req.TimeFrameRequired, req.SurveyRequired).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.ResponseRequestProgram{
		Status:    "OK",
		ProgramId: newID,
	}, nil
}

func (s *AuthService) UpdateProgram(c context.Context, req models.RegisterRequestProgram) (*models.ResponseRequestProgram, error) {
	// Input validation
	if req.ProgramName == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `UPDATE stu_tracker.Programs SET program_name = $1, timeframe_required = $2, survey_required = $3 WHERE id = $4`
	_, err := s.db.ExecContext(c, query, req.ProgramName, req.TimeFrameRequired, req.SurveyRequired, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to update program: %w", err)
	}
	return &models.ResponseRequestProgram{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteProgram(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
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

func (s *AuthService) CreateProgramLocation(c context.Context, req models.RegisterLocationProgram) (*models.RegisterResponseLocationProgram, error) {
	// Input validation
	if req.LocationId == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: ID, Program name")
	}
	query := `INSERT INTO stu_tracker.Location_programs (location_id, program_id, organization_id) 
			  VALUES ($1, $2, $3);`

	_, err := s.db.ExecContext(c, query, req.LocationId, req.ProgramId, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RegisterResponseLocationProgram{
		Status: "OK",
	}, nil
}

func (s *AuthService) DeleteProgramLocation(c context.Context, req models.RemoveLocationProgram) (*models.RemoveResponse, error) {
	// Input validation
	if req.LocationId == nil || req.ProgramId == nil {
		return nil, fmt.Errorf("missing required fields: ID")
	}
	query := `DELETE FROM stu_tracker.Location_programs WHERE location_id = $1 AND program_id = $2`

	_, err := s.db.ExecContext(c, query, *req.LocationId, *req.ProgramId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RemoveResponse{
		Status: "OK",
	}, nil
}

func (s *AuthService) CreateProgramSurvey(c context.Context, req models.RegisterProgramSurvey) (*models.RegisterResponseProgramSurvey, error) {
	// Input validation
	if req.SurveyId == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: surveyid programid and organization")
	}
	query := `INSERT INTO stu_tracker.Program_survey (survey_id, program_id, organization_id) 
			  VALUES ($1, $2, $3);`

	_, err := s.db.ExecContext(c, query, req.SurveyId, req.ProgramId, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RegisterResponseProgramSurvey{
		Status: "Created",
	}, nil
}

func (s *AuthService) DeleteProgramSurvey(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {

	query := `DELETE FROM stu_tracker.Program_survey WHERE organization_id = $1 AND id = $2`

	_, err := s.db.ExecContext(c, query, req.OrganizationId, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
