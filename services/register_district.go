package services

import (
	"fmt"
	"tracker/app/models"
)

func (s *AuthService) AddDistrict(req models.RegisterRequestDistrict) (*models.ResponseRequestDistrict, error) {
	// Input validation
	if req.Name == "" || req.Region == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: Name, Adminid, Region")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.District (name, city, region, state, organization_id) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err := s.db.QueryRow(query, req.Name, req.City, req.Region, req.State, *req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.ResponseRequestDistrict{
		Status:     "OK",
		DistrictId: newID,
	}, nil
}

func (s *AuthService) UpdateDistrict(req models.RegisterRequestDistrict) (*models.ResponseUpdate, error) {
	// Input validation
	if req.Name == "" || req.Region == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: Name, Region, Orgid")
	}
	query := `UPDATE stu_tracker.District SET name = $1, city = $2, region = $3, state = $4, organization_id = $5 WHERE id = $6`
	_, err := s.db.Exec(query, req.Name, req.City, req.Region, req.State, *req.OrganizationId, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteDistrict(req models.RemoveRequest) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: ID")
	}
	query := `DELETE FROM stu_tracker.District WHERE id = $1`
	_, err := s.db.Exec(query, *req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
