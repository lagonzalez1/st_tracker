package services

import (
	"database/sql"
	"fmt"
	"net/mail"
	"tracker/app/database"
	"tracker/app/models"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rabbitmq/amqp091-go"
	"github.com/valkey-io/valkey-go"
	"golang.org/x/crypto/bcrypt"
)

/*
	//Business logic for user creation
	//Password hashing
	//Database insertion

*/

type MQChannels struct {
	Connection *amqp091.Connection
	Channels   map[string]*amqp091.Channel // Keyed by task type
}

type AuthService struct {
	db *sql.DB
	s3 *s3.Client
	mq *database.MQChannels
	vk valkey.Client
}

// Return the struct object by reference.
func NewAuthService(db *sql.DB, client *s3.Client, mq *database.MQChannels, vk valkey.Client) *AuthService {
	return &AuthService{
		db: db,
		s3: client,
		mq: mq,
		vk: vk,
	}
}

// This is a hacky way of getting the Valkey reference
// Can this be improved ??
func (s *AuthService) Valkey() valkey.Client {
	return s.vk
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

	// Query to insert into Admin_root
	query := `INSERT INTO 
			  stu_tracker.Admin_root(email, password_hash, organization_name, organization_id) 
			  VALUES ($1, $2, $3, $4) RETURNING id, email, organization_id`

	// Query to insert into Organization
	query2 := `INSERT INTO
				stu_tracker.Organization(title, address, city, zip_code, state)
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	// Hash the password
	unhashedPassword := []byte(req.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(unhashedPassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %v", err)
	}

	// Step 2: Use the ID from the first query to insert into Organization
	var organizationID int
	err = s.db.QueryRow(query2,
		req.OrganizationName,
		req.Address,
		req.City,
		req.ZipCode,
		req.State,
	).Scan(&organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %v", err)
	}

	var user models.RegisterResponseAdminRoot
	err = s.db.QueryRow(query, req.Email, string(hashedPassword), req.OrganizationName, organizationID).Scan(
		&user.ID,
		&user.Email,
		&user.OrganizationId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin_root: %v", err)
	}

	// Optionally include the organization ID in the response or log it
	fmt.Printf("Created Organization with ID: %d\n", organizationID)

	return &user, nil
}
