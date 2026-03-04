package services

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tracker/app/config"
	"tracker/app/models"

	"github.com/stripe/stripe-go/v83"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) AddOrganization(c context.Context, req models.RegisterOrganization) (*models.RegisterOrganizationResponse, error) {
	// Input validation
	env_config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load env file")
	}
	if req.Email == nil || req.Password == nil || req.Code == nil {
		return nil, fmt.Errorf("missing required fields: email, password")
	}
	if *req.Code == env_config.SignUp.ORG_ADD_KEY || *req.Code == env_config.SignUp.ORG_ADD_KEY_1 || *req.Code == env_config.SignUp.ORG_ADD_KEY_2 || *req.Code == env_config.SignUp.ORG_ADD_KEY_3 {
		var orgID *int64
		query := `INSERT INTO stu_tracker.Organization(title, address, zip_code, state, city) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id;`

		err = s.db.QueryRowContext(c, query, req.OrganizationName, req.Address, req.ZipCode, req.State, req.City).Scan(&orgID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
		}

		unhashed_password := []byte(*req.Password)
		hash_password, err := bcrypt.GenerateFromPassword(unhashed_password, bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("unable to generate password %w", err)
		}
		var AdminID int64
		email := strings.ToLower(*req.Email)
		query2 := `INSERT INTO stu_tracker.Admin_root(fullname, email, password_hash, organization_id, organization_name)
				VALUES ($1,$2,$3,$4,$5) RETURNING id;`
		err = s.db.QueryRow(query2, req.Fullname, email, hash_password, orgID, req.OrganizationName).Scan(&AdminID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert student: %w", err)
		}
		return &models.RegisterOrganizationResponse{Status: "Ok", OrganizationId: orgID}, nil
	}
	return nil, fmt.Errorf("unable to create organization")
}

func (s *AuthService) AddOrganizationInitSubscription(c context.Context, orgID *int64, customer *stripe.Customer, subscription *stripe.Subscription) error {
	now := time.Now().UTC()
	const (
		Plan_id            = 1
		Status             = "Trial active"
		SubscriptionStatus = "Status ok"
	)
	CurrentPeriodStart := now.Truncate(time.Second)
	CurrentPeriodEnd := now.AddDate(0, 0, 15).Truncate(time.Second)
	subscribeBasicPlan := `INSERT INTO stu_tracker.organization_subscription
		(organization_id, plan_id, status, current_period_start, current_period_end, stripe_customer_id, stripe_subscription_id, subscription_status)
		VALUES ($1,$2,$3,$4, $5, $6,$7,$8);`
	_, err := s.db.ExecContext(c, subscribeBasicPlan, orgID, Plan_id, Status, CurrentPeriodStart, CurrentPeriodEnd, customer.ID, subscription.ID, SubscriptionStatus)
	if err != nil {
		return err
	}
	return nil
}
