package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	"tracker/app/config"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CookieSettings struct {
}

type UserFinder interface {
	FindByEmail(email string) (*models.User, error)
	GetPermissions(userID, orgID int64) ([]string, error)
}

func (s *AuthService) LoginAction(req models.LoginRequest) (*models.LoginResponse, error) {
	var user *models.User
	var err error
	var locations_list []models.TutorLocationList
	var program_list []models.ResponseRequestProgramList

	switch req.Type {
	case "ROOT":
		// PERMISSIONS MIGHT BE WONKY
		user, err = s.findRootUser(req.Email)
	case "ADMIN":
		user, err = s.findAdminUser(req.Email)
	case "TUTOR":
		user, locations_list, program_list, err = s.findTutorUser(req.Email)
	default:
		return nil, errors.New("login type was not specified")
	}
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("wrong password")
	}
	// 20 min
	token, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("unable to create access token: %w", err)
	}
	// 5 + hours
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("unable to create refresh token: %w", err)
	}

	return &models.LoginResponse{
		Token:        &token,
		RefreshToken: &refreshToken,
		User: models.User{
			ID:             user.ID,
			Email:          user.Email,
			OrganizationId: user.OrganizationId,
			Permissions:    user.Permissions,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Type:           user.Type,
		},
		Permissions: models.LoginResponsePermissions{
			DisableUpdate: false,
			DisableCreate: false,
			DisableDelete: false,
		},
		TutorLocations: locations_list,
		TutorPrograms:  program_list,
	}, nil
}

// Helper functions HARD CODEDE
func (s *AuthService) findRootUser(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, organization_id, fullname FROM stu_tracker.Admin_root WHERE email = $1`
	user := &models.User{Type: "ROOT"}
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.OrganizationId,
		&user.FirstName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no root user found with email: %s", email)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	user.Permissions = []string{"write:*", "delete:*", "view:*"}
	return user, nil
}

func (s *AuthService) findAdminUser(email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, organization_id, fullname FROM stu_tracker.Admin_staff WHERE email = $1`
	user := &models.User{Type: "ADMIN"}
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.OrganizationId,
		&user.FirstName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no admin found with email: %s", email)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	user.Permissions, err = s.getAdminPermissions(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return user, nil
}

func (s *AuthService) findTutorUser(email string) (*models.User, []models.TutorLocationList, []models.ResponseRequestProgramList, error) {
	query := `SELECT id, email, password_hash, organization_id, first_name, last_name FROM stu_tracker.Tutors WHERE email = $1`
	user := &models.User{Type: "TUTOR"}
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.OrganizationId,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil, fmt.Errorf("no tutor found with email: %s", email)
		}
		return nil, nil, nil, fmt.Errorf("database error: %w", err)
	}

	// Get permissions
	user.Permissions, err = s.getTutorPermissions(user.ID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	// Get locations and programs
	locations, err := s.GetTutorLocations(user.ID, user.OrganizationId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get tutor locations: %w", err)
	}

	var programs []models.ResponseRequestProgramList
	if len(locations) > 0 {
		locationIDs := make([]int64, len(locations))
		for i, loc := range locations {
			locationIDs[i] = loc.ID
		}

		programs, err = s.GetProgramsByIds(locationIDs, user.OrganizationId)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to get programs: %w", err)
		}
	}

	return user, locations, programs, nil
}

func (s *AuthService) getAdminPermissions(userID int64) ([]string, error) {
	query := `SELECT name FROM stu_tracker.Permissions p 
              INNER JOIN stu_tracker.Admin_permissions ad 
              ON p.id = ad.permission_id 
              WHERE ad.admin_id = $1;`

	return s.queryPermissions(query, userID)
}

func (s *AuthService) getTutorPermissions(id int64) ([]string, error) {
	query := `SELECT name FROM stu_tracker.Permissions p 
              LEFT JOIN stu_tracker.Tutor_permissions tp 
              ON p.id = tp.permission_id
		  	  WHERE tp.tutor_id = $1;`

	return s.queryPermissions(query, id)
}

func (s *AuthService) queryPermissions(query string, args ...interface{}) ([]string, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		permissions = append(permissions, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return permissions, nil
}

func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	envConfig, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	claims := jwt.MapClaims{
		"sub":   user.Email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"type":  user.Type,
		"id":    user.ID,
		"orgid": user.OrganizationId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(envConfig.JWT))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	log.Printf("Generated access token for %s (type: %s)", user.Email, user.Type)
	log.Printf("Token length for %s: %d", user.Type, len(tokenString))

	return tokenString, nil
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	envConfig, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	// About 6 months of access until
	claims := jwt.MapClaims{
		"sub":   user.Email,
		"exp":   time.Now().Add(4380 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"type":  user.Type,
		"id":    user.ID,
		"orgid": user.OrganizationId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(envConfig.JWT))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}
