package services

import (
	"context"
	"database/sql"
	"time"
	"tracker/app/models"

	"github.com/stripe/stripe-go/v83"
)

func (s *AuthService) HasProcessedEvent(ctx context.Context, stripeCustomerId *string, eventID *string) (bool, error) {
	var exist bool
	query := `SELECT EXISTS (SELECT 1 FROM stu_tracker.organization_subscription WHERE stripe_customer_id = $1 AND event_id = $2);`
	err := s.db.QueryRowContext(ctx, query, *stripeCustomerId, eventID).Scan(exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (s *AuthService) LinkCheckout(sess *stripe.CheckoutSession, OrganizationID *int64, eventID string) error {
	// Search for the user with session_id and customer_id ?
	query := `UPDATE stu_tracker.organization_subscription SET
		stripe_customer_id = $1, stripe_subscription_id = $2, stripe_session_id = $3, event_id = $4
		WHERE organization_id = $5`
	_, err := s.db.Exec(query, sess.Customer.ID, sess.Subscription.ID, sess.ID, eventID, OrganizationID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GrantAccess(sess *stripe.Subscription, OrganizationID *int64, plan *models.SubscriptionPlan) error {
	query := `UPDATE stu_tracker.organization_subscription 
	SET plan_id = $1, status = $2, current_period_start = $3, current_period_end = $4, cancel_at = $5, stripe_session_id = $6
	WHERE organization_id = $7 AND stripe_customer_id = $8`
	startDate := sess.Items.Data[0].CurrentPeriodStart
	endDate := sess.Items.Data[0].CurrentPeriodEnd
	tStartDate := time.Unix(startDate, 0)
	tEndDate := time.Unix(endDate, 0)
	_, err := s.db.Exec(query, plan.ID, "ACTIVE", tStartDate, tEndDate, nil, nil, OrganizationID, sess.Customer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) VerifyAccess(ctx context.Context, inv *stripe.Invoice, admin *models.Admin) error {
	var verify bool
	query := `UPDATE stu_tracker.organization_subscription SET latest_invoice_url = $1 WHERE organization_id = $2;`
	// Reset limits
	err := s.db.QueryRowContext(ctx, query, inv.HostedInvoiceURL, admin.OrganizationID).Scan(&verify)
	if err != nil {
		return err
	}
	// Reset limits
	mLimits := `UPDATE stu_tracker.Generate_materials_task SET input_tokens = 0, output_tokens = 0 WHERE organization_id = $1;`
	qLimits := `UPDATE stu_tracker.Generate_questions_task SET input_tokens = 0, output_tokens = 0 WHERE organization_id = $1;`
	_, err = s.db.ExecContext(ctx, mLimits, admin.OrganizationID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, qLimits, admin.OrganizationID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) UpdateAccess(sess *stripe.Subscription, admin *models.Admin, plan *models.SubscriptionPlan) error {
	query := `UPDATE stu_tracker.organization_subscription 
	SET plan_id = $1, status = $2, current_period_start = $3, current_period_end = $4, cancel_at = $5
	WHERE organization_id = $6 AND stripe_customer_id = $7`
	startDate := sess.Items.Data[0].CurrentPeriodStart
	endDate := sess.Items.Data[0].CurrentPeriodEnd
	tStartDate := time.Unix(startDate, 0)
	tEndDate := time.Unix(endDate, 0)
	status := sess.Status
	var cancelAt sql.NullTime
	if status == stripe.SubscriptionStatusCanceled || status == stripe.SubscriptionStatusUnpaid || status == stripe.SubscriptionStatusPastDue {
		cancelAt = sql.NullTime{Time: time.Unix(sess.CanceledAt, 0).UTC(), Valid: true}
	} else {
		cancelAt = sql.NullTime{Valid: false}
	}
	_, err := s.db.Exec(query, plan.ID, status, tStartDate, tEndDate, cancelAt, admin.OrganizationID, sess.Customer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) RevokeAccess(sess *stripe.Subscription, admin *models.Admin, plan *models.SubscriptionPlan) error {
	query := `UPDATE stu_tracker.organization_subscription 
	SET plan_id = $1, status = $2, current_period_start = $3, current_period_end = $4, cancel_at = $5, stripe_session_id = $6
	WHERE organization_id = $7 AND stripe_customer_id = $8`
	startDate := time.Now()
	endDate := time.Now()
	_, err := s.db.Exec(query, 0, "CANCELED", startDate, endDate, time.Now(), nil, admin.OrganizationID, sess.Customer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) FailedPayment(sess *stripe.Invoice, admin *models.Admin) error {
	query := `UPDATE stu_tracker.organization_subscription 
	SET status = $1, current_period_start = $2, current_period_end = $3, cancel_at = $4, stripe_session_id = $5
	WHERE organization_id = $6 AND stripe_customer_id = $7`
	startDate := time.Now()
	endDate := time.Now()
	_, err := s.db.Exec(query, "FAILED", startDate, endDate, time.Now().Add(time.Hour*48), nil, admin.OrganizationID, sess.Customer.ID)
	if err != nil {
		return err
	}
	return nil
}
