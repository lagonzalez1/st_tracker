package services

import (
	"fmt"
	"tracker/app/config"
	"tracker/app/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) AddOrganization(req models.RegisterOrganization) (*models.RegisterOrganizationResponse, error) {
	// Input validation
	env_config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load env file")
	}
	if req.Email == nil || req.Password == nil || req.Code == nil {
		return nil, fmt.Errorf("missing required fields: email, password")
	}
	fmt.Print(env_config.SignUp.ORG_ADD_KEY)
	if *req.Code == env_config.SignUp.ORG_ADD_KEY || *req.Code == env_config.SignUp.ORG_ADD_KEY_1 || *req.Code == env_config.SignUp.ORG_ADD_KEY_2 || *req.Code == env_config.SignUp.ORG_ADD_KEY_3 {
		var orgID *int64
		query := `INSERT INTO stu_tracker.Organization (title, address, zip_code, state, city) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id;`

		err = s.db.QueryRow(query, req.OrganizationName, req.Address, req.ZipCode, req.State, req.City).Scan(&orgID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
		}

		unhashed_password := []byte(*req.Password)
		hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("unable to generate password %w", err)
		}
		var AdminID int64
		query2 := `INSERT INTO stu_tracker.Admin_root(fullname, email, password_hash, organization_id, organization_name)
				VALUES ($1,$2,$3,$4,$5) RETURNING id;`
		err = s.db.QueryRow(query2, req.Fullname, req.Email, hash_password, orgID, req.OrganizationName).Scan(&AdminID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
		}
		return &models.RegisterOrganizationResponse{
			Status: "OK",
		}, nil
	}
	return &models.RegisterOrganizationResponse{Status: "Error"}, nil
}
