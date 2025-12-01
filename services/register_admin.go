package services

import (
	"context"
	"fmt"
	"strings"
	"tracker/app/models"

	"golang.org/x/crypto/bcrypt"
)

/**
 Admin types:
 For Placment info regarding substitutes and logistics
 View: locations, materail, district, programs, students, subjects, tutors
 Delete: Tutor-locations

**/

const (
	Root  = "root"
	Admin = "admin"
	Tutor = "tutor"
)

func (s *AuthService) AddAdminStaff(c context.Context, req models.RegisterRequestAdmin) (*models.ResponseRequestAdmin, error) {
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
	query := `INSERT INTO stu_tracker.Admin_staff(fullname, email, password_hash, region, state, organization_id, district_id, active) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`
	email := strings.ToLower(req.Email)
	err = s.db.QueryRowContext(c, query, req.Fullname, email, string(hash_password), req.Region, req.State, *req.OrganizationId, req.DistrictId, req.Active).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}
	return &models.ResponseRequestAdmin{
		Status:   "OK",
		Admin_id: newID,
	}, nil
}

func (s *AuthService) AddAdminLocation(c context.Context, req models.RegisterAdminLocation, orgid *int64) (*models.ResponseRequestAdmin, error) {
	// Input validation
	if req.LocationID == nil {
		return nil, fmt.Errorf("missing required fields: location_id")
	}
	query := `INSERT INTO stu_tracker.admin_location_access(admin_id, location_id, organization_id) 
			  VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(c, query, req.AdminID, req.LocationID, orgid)
	if err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}

	return &models.ResponseRequestAdmin{
		Status:   "OK",
		Admin_id: *req.AdminID,
	}, nil
}

func (s *AuthService) UpdateAdminStaff(c context.Context, req models.RegisterRequestAdmin) (*models.ResponseUpdateAdmin, error) {
	if req.ID == nil || req.Email == "" || req.OrganizationId == nil {
		return nil, fmt.Errorf("missing required fields: id, email, org_id")
	}
	if req.Password != "" {
		unhashed_password := []byte(req.Password)
		hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("unable to hash password: %v", err)
		}
		query := `UPDATE stu_tracker.Admin_staff SET fullname = $1, email = $2, password_hash = $3, region = $4, state = $5, organization_id = $6, district_id = $7, active = $8
			  WHERE id = $9;`
		email := strings.ToLower(req.Email)
		_, err = s.db.Exec(query, req.Fullname, email, string(hash_password), req.Region, req.State, *req.OrganizationId, req.DistrictId, req.Active, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert admin: %w", err)
		}
	} else {
		query := `UPDATE stu_tracker.Admin_staff SET fullname = $1, email = $2, region = $3, state = $4, organization_id = $5, district_id = $6, active = $7
		WHERE id = $8;`
		email := strings.ToLower(req.Email)
		_, err := s.db.ExecContext(c, query, req.Fullname, email, req.Region, req.State, *req.OrganizationId, req.DistrictId, req.Active, req.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert admin: %w", err)
		}
	}
	return &models.ResponseUpdateAdmin{
		ID:     req.ID,
		Status: "Updated",
	}, nil
}

func (s *AuthService) DeleteAdminStaff(c context.Context, req models.RemoveAdmin) (*models.RemoveResponse, error) {
	if req.ID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `DELETE FROM stu_tracker.Admin_staff WHERE id = $1;`
	_, err := s.db.ExecContext(c, query, req.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete staff: %w", err)

	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}

func (s *AuthService) DeleteAdminLocation(c context.Context, req models.RemoveAdminLocation) (*models.RemoveResponse, error) {
	if req.AdminID == nil {
		return nil, fmt.Errorf("missing required fields: first_name, last_name, or email")
	}
	query := `DELETE FROM stu_tracker.admin_location_access WHERE admin_id = $1 AND location_id = $2;`
	_, err := s.db.ExecContext(c, query, req.AdminID, req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("unable to delete staff: %w", err)

	}
	return &models.RemoveResponse{
		Status: "Removed",
	}, nil
}
