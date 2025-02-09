package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"tracker/app/config"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) LoginAction(req models.LoginRequest) (*models.LoginResponse, error) {
	var user *models.User
	if req.Type == "ROOT" {
		// 1. Find the user by email
		root_user, err := s.findUserByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		user = root_user
		user.Type = "ROOT"
	} else if req.Type == "ADMIN" {
		// 1. Find the user by email
		admin_user, err := s.findAdminByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		user = admin_user
		user.Type = "ADMIN"
	} else if req.Type == "TUTOR" {
		// 1. Find the user by email
		tutor_user, err := s.findATutorByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		user = tutor_user
		user.Type = "TUTOR"
	} else {
		return nil, errors.New("login type was not specified")
	}

	// 2. Compare passwords
	if err := comparePasswords(user.Password, req.Password); err != nil {
		return nil, errors.New("wrong password")
	}

	// 3. Generate JWT token
	token, err := generateJWTToken(user)
	if err != nil {
		return nil, errors.New("unable to create JWT token")
	}
	refreshToken, err := generateJWTTokenRefresh(user)
	if err != nil {
		return nil, errors.New("unable to create refresh JWT token")
	}

	// 4. Return by reference the models.LoginResonse { Token: User: {ID, Username, Email}, nil}
	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: models.User{
			ID:             user.ID,
			Email:          user.Email,
			OrganizationId: user.OrganizationId,
			Permissions:    user.Permissions,
		},
		Permissions: models.LoginResponsePermissions{
			DisableUpdate: false,
			DisableCreate: false,
			DisableDelete: false,
		},
	}, nil
}

func comparePasswords(hashedPassword, inputString string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(inputString),
	)
}

func generateJWTTokenRefresh(user *models.User) (string, error) {
	env_config, err := config.LoadConfig()
	if err != nil {
		return fmt.Sprintf("unable to load config env"), err
	}
	jwt_token := env_config.JWT
	secret_key := []byte(jwt_token)
	// Need to create a role to return here ?
	jwt_object := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         user.Email,
		"exp":         time.Now().Add(5 * time.Hour).Unix(),
		"iat":         time.Now().Unix(),
		"type":        user.Type,
		"permissions": user.Permissions,
	})
	token_string, err := jwt_object.SignedString(secret_key)
	if err != nil {
		return fmt.Sprintf("Unable to create JWT token"), err
	}
	return token_string, nil
}

func generateJWTToken(user *models.User) (string, error) {
	env_config, err := config.LoadConfig()
	if err != nil {
		return fmt.Sprintf("unable to load config env"), err
	}
	jwt_token := env_config.JWT
	secret_key := []byte(jwt_token)
	// Need to create a role to return here ?
	jwt_object := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         user.Email,
		"exp":         time.Now().Add(10 * time.Minute).Unix(),
		"iat":         time.Now().Unix(),
		"type":        user.Type,
		"permissions": user.Permissions,
	})
	token_string, err := jwt_object.SignedString(secret_key)
	if err != nil {
		return fmt.Sprintf("Unable to create JWT token"), err
	}
	return token_string, nil
}

func (s *AuthService) findUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, organization_id from stu_tracker.Admin_root WHERE email = $1`
	var user models.User
	// Create copies of the id, username, email, password
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.OrganizationId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email: %s", email)
		}
		return nil, err
	}
	user.Permissions = []string{"all"}
	return &user, nil
}

func (s *AuthService) findAdminByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, organization_id from stu_tracker.Admin_staff WHERE email = $1`
	var user models.User
	// Create copies of the id, username, email, password
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.OrganizationId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find admin: %s", email)
		}
		return nil, err
	}

	permissionQuery := `SELECT name from stu_tracker.permissions p LEFT JOIN
						stu_tracker.admin_role_permissions ad ON p.id = ad.permission_id WHERE
						ad.user_id = $1 AND p.organization_id = $2;`
	rows, err := s.db.Query(permissionQuery, user.ID, user.OrganizationId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to get permissions: %s", email)
		}
		return nil, err
	}
	defer rows.Close()
	var permissions []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		permissions = append(permissions, name)
	}
	user.Permissions = permissions
	return &user, nil
}

func (s *AuthService) findATutorByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, organization_id from stu_tracker.Tutor WHERE email = $1`
	var user models.User
	// Create copies of the id, username, email, password
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.OrganizationId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email: %s", email)
		}
		return nil, err
	}

	permissionQuery := `SELECT name from stu_tracker.permissions p LEFT JOIN
						stu_tracker.admin_role_permissions ad ON p.id = ad.permission_id WHERE
						ad.user_id = $1 AND p.organization_id = $2;`
	rows, err := s.db.Query(permissionQuery, user.ID, user.OrganizationId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to get permissions: %s", email)
		}
		return nil, err
	}
	defer rows.Close()
	var permissions []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		permissions = append(permissions, name)
	}
	user.Permissions = permissions
	return &user, nil
}
