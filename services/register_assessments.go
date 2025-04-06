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
	query := `INSERT INTO stu_tracker.Assessments
			(title, description, letter, cycle, alpha_identifier, external_link, max_score, subject_id, material_id, 
			organization_id, visible, program_id, version, pre, post, mid) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id;`

	err := s.db.QueryRow(query, req.Title, req.Description, req.Letter,
		req.Cycle, req.AlphaIdentifier, req.ExternalLink, req.MaxScore,
		req.SubjectId, req.MaterialID, req.OrganizationID, req.Visible,
		req.ProgramId, req.Version, req.Pre, req.Post, req.Mid).Scan(&newID)
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
	query := `UPDATE stu_tracker.Assessments SET title = $1, description = $2, letter = $3, 
	cycle = $4, alpha_identifier = $5, external_link = $6, max_score = $7, subject_id = $8, 
	material_id = $9, organization_id = $10, visible = $11, program_id = $12, version = $13, pre = $14, post = $15, mid = $16 
	WHERE id = $17;`

	_, err := s.db.Exec(query, req.Title, req.Description, req.Letter, req.Cycle,
		req.AlphaIdentifier, req.ExternalLink, req.MaxScore, req.SubjectId, req.MaterialID,
		req.OrganizationID, req.Visible, req.ProgramId, req.Version, req.Pre, req.Post, req.Mid, req.ID)
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
