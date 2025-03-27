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
	query := `INSERT INTO stu_tracker.Tutors(first_name, last_name, email, password_hash, organization_id, location_id) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	err = s.db.QueryRow(query, req.FirstName, req.LastName, req.Email, hash_password, *req.OrganizationId, *req.LocationId).Scan(&newID)
	if err != nil {
		return nil, err
	}
	if req.LocationId != nil {
		queryLocationLink := `INSERT INTO stu_tracker.Tutor_locations(tutor_id, location_id, organization_id) 
							  VALUES ($1, $2, $3); `
		_, err = s.db.Exec(queryLocationLink, newID, *req.LocationId, *req.OrganizationId)
		if err != nil {
			fmt.Printf("Database query failed at AddTutor: %v", err)
			return nil, err
		}
	}

	permissionQuery := `INSERT INTO stu_tracker.Tutor_Permissions (tutor_id, permission_id)
						VALUES ($1, 25), ($1, 24), ($1, 6), ($1, 39);`
	_, err = s.db.Exec(permissionQuery, newID)
	if err != nil {
		fmt.Printf("Unable to add permissions to tutor: %v", err)
		return nil, err
	}

	return &models.ResponseRequestTutor{
		Status:  "OK",
		TutorId: newID,
	}, nil
}

func (s *AuthService) AddTutorLocation(req models.RegisterTutorLocation) (*models.RegisterTutorResponse, error) {
	// Input validation
	if req.LocationId == nil || req.TutorId == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	query := `INSERT INTO stu_tracker.Tutor_locations (tutor_id, location_id, organization_id) 
			  VALUES ($1, $2, $3) ON CONFLICT (tutor_id, location_id, organization_id) DO NOTHING;`
	_, err := s.db.Exec(query, *req.TutorId, *req.LocationId, *req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RegisterTutorResponse{
		Status: "OK",
	}, nil
}

func (s *AuthService) UpdateTutor(req models.RegisterRequestTutor) (*models.ResponseUpdate, error) {
	// Input validation
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}

	// This can be optimzed....
	if req.Password != "" {
		unhashed_password := []byte(req.Password)
		hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		values := []interface{}{}
		query := `UPDATE stu_tracker.Tutors 
		SET first_name = $1, last_name = $2, organization_id = $3, location_id = $4, password_hash = $5`
		values = append(values, req.FirstName, req.LastName, *req.OrganizationId, *req.LocationId, string(hash_password))
		if req.EmailChange != "" {
			query += `, email = $6 WHERE id = $7`
			values = append(values, req.EmailChange, req.ID)
		}
		query += " WHERE id = $6"
		values = append(values, req.ID)
		_, err = s.db.Exec(query, values...)
		if err != nil {
			return nil, err
		}
	} else {
		values := []interface{}{}
		query := `UPDATE stu_tracker.Tutors 
		SET first_name = $1, last_name = $2, organization_id = $3, location_id = $4`
		values = append(values, req.FirstName, req.LastName, *req.OrganizationId, *req.LocationId)
		if req.EmailChange != "" {
			query += `, email = $5 WHERE id = $6`
			values = append(values, req.EmailChange, req.ID)
		}
		query += " WHERE id = $5"
		values = append(values, req.ID)
		_, err := s.db.Exec(query, values...)
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
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

func (s *AuthService) DeleteTutorLocation(req models.RemoveTutorLocation) (*models.RemoveResponse, error) {
	// Input validation
	if req.LocationId == nil || req.TutorId == nil || req.OrganizationID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, email, password")
	}
	query := `DELETE FROM stu_tracker.Tutor_locations 
			  WHERE tutor_id = $1 
			  AND location_id = $2 
			  AND organization_id = $3;`
	_, err := s.db.Exec(query, *req.TutorId, *req.LocationId, *req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
