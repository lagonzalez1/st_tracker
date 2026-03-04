package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"tracker/app/models"

	"github.com/stripe/stripe-go/v83"
	portalsession "github.com/stripe/stripe-go/v83/billingportal/session"
	"github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/subscription"
)

// purchase_intent -> New user, create subscription and save the customer_id
// purchase_portal -> User exist, find user  by DB search.

// Checkout session returns a status str and URL with stripe session.
// Initially no customer_id has been attached.
func (s *AuthService) CreateCheckoutSession(ctx context.Context, model *models.PurchaseIntent, orgid *int64) (*models.StripeCreateCheckoutSessionResponse, error) {
	stripe.Key = os.Getenv("STRIPE_SECRET")
	var domain string
	if os.Getenv("APP_ENV") == "dev" {
		domain = fmt.Sprintf(`%s/dashboard`, os.Getenv("DEV_URL"))
	} else {
		domain = fmt.Sprintf(`%s/dashboard`, os.Getenv("PROD_URL"))
	}
	p := &stripe.CheckoutSessionLineItemParams{
		Price:    stripe.String(*model.PriceID),
		Quantity: stripe.Int64(1),
	}
	checkoutParams := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems:  []*stripe.CheckoutSessionLineItemParams{p},
		SuccessURL: stripe.String(domain + "?success=true&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(domain + "?canceled=true"),
		Customer:   model.CustomerID,
		Metadata: map[string]string{
			"orgid": fmt.Sprint(*orgid),
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"orgid": fmt.Sprint(*orgid),
			},
		},
	}
	stripeSession, err := session.New(checkoutParams)
	if err != nil {
		log.Printf("session.New: %v", err)
	}
	// Save the user session and customer id to reference in webhook
	query := `UPDATE stu_tracker.Admin_root SET stripe_session_id = $1
			  WHERE organization_id = $2;`

	_, err = s.db.ExecContext(ctx, query, &stripeSession.ID, orgid)
	if err != nil {
		return nil, err
	}
	return &models.StripeCreateCheckoutSessionResponse{
		Status: "OK",
		URL:    &stripeSession.URL,
	}, nil

}

// Return a url with a existing user.
func (s *AuthService) CreatePortalSession(ctx context.Context, orgid *int64) (*models.StripeCheckoutSession, error) {
	var stripe_customer_id *string
	query := `SELECT stripe_customer_id FROM stu_tracker.organization_subscription WHERE organization_id = $1`
	err := s.db.QueryRowContext(ctx, query, orgid).Scan(&stripe_customer_id)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch credentials: %v", err)
	}
	if stripe_customer_id == nil {
		return nil, fmt.Errorf("stripe customer not found")
	}
	var domain string
	if os.Getenv("APP_ENV") == "dev" {
		domain = fmt.Sprintf(`%s/dashboard`, os.Getenv("DEV_URL"))
	} else {
		domain = fmt.Sprintf(`%s/dashboard`, os.Getenv("PROD_URL"))
	}

	stripe.Key = os.Getenv("STRIPE_SECRET")
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(*stripe_customer_id),
		ReturnURL: stripe.String(domain),
	}
	ps, err := portalsession.New(params)
	if err != nil {
		return nil, err
	}
	return &models.StripeCheckoutSession{
		Status: "OK",
		URL:    &ps.URL,
	}, nil
}

func (s *AuthService) GetOrganizationIdByCustomerId(ctx context.Context, cid string) (*models.Admin, error) {
	q := `SELECT id, organization_id FROM stu_tracker.organization_subscription 
	WHERE stripe_customer_id = $1;`
	var admin models.Admin
	err := s.db.QueryRowContext(ctx, q, cid).Scan(&admin.ID, &admin.OrganizationID)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *AuthService) GetAdminByCustomerId(ctx context.Context, cid string) (*models.Admin, error) {
	q := `SELECT id, organization_id FROM stu_tracker.Admin_root 
	WHERE stripe_customer_id = $1;`
	var admin models.Admin
	err := s.db.QueryRowContext(ctx, q, cid).Scan(&admin.ID, &admin.OrganizationID)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *AuthService) GetSubscriptionById(spid string) (*models.SubscriptionPlan, error) {

	q := `SELECT id, code, name, is_active FROM stu_tracker.subscription_plan 
	WHERE stripe_price_id = $1;`
	var plan models.SubscriptionPlan
	err := s.db.QueryRow(q, spid).Scan(&plan.ID, &plan.Code, &plan.Name, &plan.IsActive)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *AuthService) GetStripePriceId(c context.Context, id int) (*string, error) {
	q := `SELECT stripe_price_id FROM stu_tracker.subscription_plan WHERE id = $1;`
	var stripeId *string
	err := s.db.QueryRowContext(c, q, id).Scan(&stripeId)
	if err != nil {
		return nil, err
	}
	return stripeId, nil
}

func (s *AuthService) GetSubscriptionPlanByPriceID(pid string) (*models.AvilableSubscriptions, error) {
	q := `SELECT id, code, name, stripe_price_id, is_active FROM stu_tracker.subscription_plan; 
	WHERE stripe_price_id = $1; `
	var d *models.AvilableSubscriptions
	err := s.db.QueryRow(q, pid).Scan(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func CreateStripeCustomer(email, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}
	return customer.New(params)
}

func CreateTrialSubscription(stripeCustomerID, priceID string, trialDays int64) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		TrialPeriodDays: stripe.Int64(trialDays),
		// Subscription won't charge until trial ends
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
	}
	return subscription.New(params)
}
