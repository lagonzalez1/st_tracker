package services

import (
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddAssessment(req models.RegisterAssessment) (*models.RegisterResponseAssessment, error) {
	// Input validation
	if req.Title == "" || req.MaxScore == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Assessments (title, description, letter, cycle, alpha_identifier, external_link, max_score, subject, material_id, organization_id, visable, program_id) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;`

	err := s.db.QueryRow(query, req.Title, req.Description, req.Letter, req.Cycle, req.AlphaIdentifier, req.ExternalLink, req.MaxScore, req.Subject, req.MaterialID, req.OrganizationID, req.Visable, req.ProgramId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RegisterResponseAssessment{
		Status:       "OK",
		AssessmentId: newID,
	}, nil
}

func (s *AuthService) UpdateAssessment(req models.RegisterAssessment) (*models.RegisterResponseAssessment, error) {
	// Input validation
	if req.Title == "" || req.MaxScore == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	var newID int64
	query := `UPDATE stu_tracker.Assessments SET title = $2, description = $3, letter = $4, 
	cycle = $5, alpha_identifier = $6, external_link = $7, max_score = $8, subject = $8, material_id = $9, organization_id = $10, visable = $11, program_id = $12) 
	WHERE id = $1;`

	_, err := s.db.Exec(query, req.ID, req.Title, req.Description, req.Letter, req.Cycle, req.AlphaIdentifier, req.ExternalLink, req.MaxScore, req.Subject, req.MaterialID, req.OrganizationID, req.Visable)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RegisterResponseAssessment{
		Status:       "OK",
		AssessmentId: newID,
	}, nil
}

func (s *AuthService) DeleteAssessment(req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	query := `DELETE FROM stu_tracker.Assessments WHERE id = $1`
	_, err := s.db.Exec(query, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
