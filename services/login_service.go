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

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// 1. Find the user by email
	user, err := s.findUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("unable to find with email")
	}
	// 2. Compare passwords
	if err := comparePasswords(user.Password, req.Password); err != nil {
		return nil, errors.New("password does not match")
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
			ID:    user.ID,
			Email: user.Email,
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
		"sub": user.Email,
		"exp": time.Now().Add(time.Minute).Unix(),
		"iat": time.Now().Unix(),
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
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	token_string, err := jwt_object.SignedString(secret_key)
	if err != nil {
		return fmt.Sprintf("Unable to create JWT token"), err
	}
	return token_string, nil
}

func (s *AuthService) findUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash from stu_tracker.Admin_root WHERE email = $1`
	var user models.User
	// Create copies of the id, username, email, password
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with email: %s", email)
		}
		return nil, err
	}
	return &user, nil
}
