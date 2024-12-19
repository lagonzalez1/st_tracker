package services

import (
	"database/sql"
	"fmt"
	"net/mail"
	"tracker/app/models"

	"golang.org/x/crypto/bcrypt"
)

/*
	//Business logic for user creation
	//Password hashing
	//Database insertion

*/

type AuthService struct {
	db *sql.DB
}

// Return the struct object by reference.
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *AuthService) RegisterRootUser(req models.RegisterRequestAdminRoot) (*models.RegisterResponseAdminRoot, error) {
	if req.Password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}
	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}
	query := `INSERT INTO 
			  stu_tracker.Admin_root(email, password_hash, organization_name) 
			  VALUES ($1, $2, $3) RETURNING id, email`
	// Convert password to byte. Do i need to check for empty password.
	unhashed_password := []byte(req.Password)
	hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %v", err)
	}
	var user models.RegisterResponseAdminRoot
	err = s.db.QueryRow(query, req.Email, string(hash_password), req.Organization).Scan(
		&user.ID,
		&user.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %v", err)
	}
	return &user, nil
}
