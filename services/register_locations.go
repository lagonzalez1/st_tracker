package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddLocation(c context.Context, req models.RegisterRequestLocation) (*models.ResponseRequestLocation, error) {
	// Input validation
	if req.Name == "" || req.Address == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Locations(name, address, city, state, zip_code, district_id, organization_id)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	err := s.db.QueryRowContext(c, query, req.Name, req.Address, req.City, req.State, req.ZipCode, req.DistrictId, *req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseRequestLocation{
		Status:     "OK",
		LocationId: newID,
	}, nil
}

func (s *AuthService) UpdateLocation(c context.Context, req models.RegisterRequestLocation) (*models.ResponseUpdate, error) {
	// Input validation
	if req.Name == "" || req.Address == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `UPDATE stu_tracker.Locations SET name = $1, address = $2, city = $3, state = $4, zip_code = $5, district_id = $6, organization_id = $7
              WHERE id = $8`

	_, err := s.db.ExecContext(c, query, req.Name, req.Address, req.City, req.State, req.ZipCode, req.DistrictId, *req.OrganizationId, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteLocation(c context.Context, req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `DELETE FROM stu_tracker.Locations WHERE id = $1`

	_, err := s.db.ExecContext(c, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
