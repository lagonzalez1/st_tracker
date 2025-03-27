package services

import (
	"fmt"
	"tracker/app/models"

	"golang.org/x/crypto/bcrypt"
)

/**
 Admin types:

 For Placment info regarding substitutes and logistics
 View: locations, materail, district, programs, students, subjects, tutors
 Delete: Tutor-locations

**/

func (s *AuthService) AddAdminStaff(req models.RegisterRequestAdmin) (*models.ResponseRequestAdmin, error) {
	// Input validation
	if req.Fullname == "" || req.Password == "" || req.Email == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	unhashed_password := []byte(req.Password)
	hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %v", err)
	}
	var newID int64
	query := `INSERT INTO stu_tracker.Admin_staff (fullname, email, password_hash, region, state, organization_id) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	err = s.db.QueryRow(query, req.Fullname, req.Email, string(hash_password), req.Region, req.State, *req.OrganizationId).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}

	permissions := `INSERT INTO stu_tracker.Admin_Permissions (admin_id, permission_id)
					VALUES ($1, 1), ($1, 12), ($1, 39), ($1, 9), ($1, 15), 
					($1, 30), ($1, 32), ($1, 33), ($1, 34), ($1, 31);`

	_, err = s.db.Exec(permissions, newID)
	if err != nil {
		fmt.Printf("Unable to add permissions in admin staff: %v", err)
		return nil, err
	}

	return &models.ResponseRequestAdmin{
		Status:   "OK",
		Admin_id: newID,
	}, nil
}

func (s *AuthService) UpdateAdminStaff(req models.RegisterRequestAdmin) (*models.ResponseUpdateAdmin, error) {
	if req.ID == nil || req.Email == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: id, email, org_id")
	}
	if req.Password != "" {
		unhashed_password := []byte(req.Password)
		hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("unable to hash password: %v", err)
		}
		query := `UPDATE stu_tracker.Admin_staff SET fullname = $1, email = $2, password_hash = $3, region = $4, state = $5, organization_id = $6 
			  WHERE id = $7;`

		_, err = s.db.Exec(query, req.Fullname, req.Email, string(hash_password), req.Region, req.State, *req.OrganizationId, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert admin: %w", err)
		}
	} else {
		query := `UPDATE stu_tracker.Admin_staff SET fullname = $1, admin_id = $2, email = $3, region = $4, state = $5, organization_id = $6 
		WHERE id = $7;`
		_, err := s.db.Exec(query, req.Fullname, req.RootId, req.Email, req.Region, req.State, *req.OrganizationId, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert admin: %w", err)
		}
	}
	return &models.ResponseUpdateAdmin{
		ID:     req.ID,
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteAdminStaff(req models.RemoveAdmin) (*models.RemoveResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `DELETE FROM stu_tracker.Admin_staff WHERE id = $1;`
	_, err := s.db.Exec(query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete staff: %w", err)

	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
