package services

import (
	"context"
	"fmt"
	"tracker/app/models"
)

/**

location_id INT REFERENCES stu_tracker.Locations(id) ON DELETE CASCADE,
    program_id INT REFERENCES stu_tracker.Programs(id) ON DELETE CASCADE,
    description VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone TEXT CHECK(phone ~ '^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$')

**/

func (s *AuthService) CreateLocationContact(ctx context.Context, req models.RegisterLocationContact) (*models.SimpleReponse, error) {
	if req.OrganizationId == nil {
		return nil, fmt.Errorf("unable to find orgid")
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Location_contacts (program_id, first_name, last_name, email, phone, notes, organization_id, location_id, description) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	err := s.db.QueryRow(query, req.ProgramID, req.FirstName, req.LastName, req.Email, req.Phone, req.Notes, req.OrganizationId, req.LocationID, req.Description).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.SimpleReponse{
		Status: "OK",
		ID:     &newID,
	}, nil

}

func (s *AuthService) UpdateLocationContact(c context.Context, req models.RegisterLocationContact) (*models.SimpleReponse, error) {
	if req.ID == nil || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `UPDATE stu_tracker.Location_contacts SET program_id = $1, first_name = $2, last_name = $3, email = $4, phone = $5, notes = $6, description = $7, program_id = $8
              WHERE id = $9;`

	_, err := s.db.ExecContext(c, query, req.ProgramID, req.FirstName, req.LastName, req.Phone, req.Email, req.Phone, req.Notes, req.Description, req.ProgramID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.SimpleReponse{
		Status: "Updated",
		ID:     req.ID,
	}, nil
}

func (s *AuthService) DeleteLocationContact(c context.Context, req models.RemoveRequestLocation) (*models.RemoveResponse, error) {
	// Input validation
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `DELETE FROM stu_tracker.Location_contacts WHERE id = $1 AND organization_id = $2`

	_, err := s.db.ExecContext(c, query, req.ID, req.OrganizationId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert Locations: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
