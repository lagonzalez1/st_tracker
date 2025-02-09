package services

import (
	"fmt"
	"tracker/app/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) AddTutor(req models.RegisterRequestTutor) (*models.ResponseRequestTutor, error) {
	// Input validation
	if req.FirstName == "" || req.LastName == "" || req.Password == "" || req.Email == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	unhashed_password := []byte(req.Password)
	hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %v", err)
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Tutors (first_name, last_name, email, password_hash, organization_id, location_id) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	err = s.db.QueryRow(query, req.FirstName, req.LastName, req.Email, hash_password, *req.OrganizationId, *req.LocationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	if req.LocationId != nil {
		queryLocationLink := `INSERT INTO stu_tracker.Tutor_locations(tutor_id, location_id, organization_id) 
							  VALUES ($1, $2, $3)`
		err := s.db.QueryRow(queryLocationLink, newID, *req.LocationId, *req.OrganizationId)
		if err != nil {
			fmt.Printf("Database query failed at AddTutor: %v", err)
			return nil, fmt.Errorf("failed to create tutor_location")
		}
	}
	return &models.ResponseRequestTutor{
		Status:  "OK",
		TutorId: newID,
	}, nil
}

func (s *AuthService) UpdateTutor(req models.RegisterRequestTutor) (*models.ResponseUpdate, error) {
	// Input validation
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}

	if req.Password != "" {
		unhashed_password := []byte(req.Password)
		hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("unable to hash password: %v", err)
		}
		query := `UPDATE TABLE stu_tracker.Tutors SET first_name = $1, last_name = $2, email = $3, organization_id = $4, location_id = $5, password_hash = $6`

		_, err = s.db.Exec(query, req.FirstName, req.LastName, req.Email, *req.OrganizationId, *req.LocationId, string(hash_password))
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
		}
	} else {
		query := `UPDATE TABLE stu_tracker.Tutors SET first_name = $1, last_name = $2, email = $3, organization_id = $4, location_id = $5`

		_, err := s.db.Exec(query, req.FirstName, req.LastName, req.Email, *req.OrganizationId, *req.LocationId)
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
		}
	}
	if req.LocationId != nil {
		queryLocationLink := `INSERT INTO stu_tracker.Tutor_locations (tutor_id, location_id, organization_id)
							  VALUES ($1, $2, $3)
							  ON CONFLICT (tutor_id, location_id) DO NOTHING;`
		_, err := s.db.Exec(queryLocationLink, req.ID, *req.LocationId, *req.OrganizationId)
		if err != nil {
			fmt.Printf("Database query failed at AddTutor: %v", err)
			return nil, fmt.Errorf("failed to create tutor_location")
		}
	}
	return &models.ResponseUpdate{
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteTutor(req models.RemoveRequest) (*models.RemoveResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	query := `DELETE FROM stu_tracker.Tutors WHERE id = $1`

	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove tutor: %w", err)
	}
	queryLocation := `DELETE FROM stu_tracker.Tutor_locations WHERE tutor_id = $1`
	_, err = s.db.Exec(queryLocation, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove tutor locations: %w", err)
	}

	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
