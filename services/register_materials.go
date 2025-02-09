package services

import (
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddMaterial(req models.RegisterRequestMaterials) (*models.ResponseRequestMaterials, error) {
	// Input validation
	if req.Title == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Materials(title, external_link, description, pre, mid, post, visable, version, organization_id, location_id)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;`

	err := s.db.QueryRow(query, req.Title, req.ExternalLink, req.Description, req.Pre, req.Mid, req.Post, req.Visable, req.Version, *req.OrganizationId, *req.LocationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseRequestMaterials{
		Status:     "OK",
		MaterialId: newID,
	}, nil
}

func (s *AuthService) UpdateMaterial(req models.RegisterRequestMaterials) (*models.ResponseUpdate, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	query := `UPDATE stu_tracker.Materials SET
			  title = $1, external_link = $2, description = $3, pre = $4, mid = $5, post = $6, visable = $7, version = $8, organization_id = $9, location_id = $10
              WHERE id = $11`

	_, err := s.db.Exec(query, req.Title, req.ExternalLink, req.Description, req.Pre, req.Mid, req.Post, req.Visable, req.Version, *req.OrganizationId, req.LocationId, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteMaterial(req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: id")
	}
	query := `DELETE FROM stu_tracker.Materials WHERE id = $1`
	_, err := s.db.Exec(query, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
