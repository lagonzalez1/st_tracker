package models

import "time"

type PurchaseIntent struct {
	PriceID    *string `json:"price_id"`
	CustomerID *string `json:"stripe_customer_id"`
}

type StripeResponse struct {
	ID     *int64  `json:"id"`
	Status *string `json:"status"`
}

type StripeCreateCheckoutSessionResponse struct {
	URL    *string `json:"url"`
	Status string  `json:"status"`
}

type StripeCheckoutSession struct {
	URL    *string `json:"url"`
	Status string  `json:"status"`
}

type Admin struct {
	ID             *int64 `json:"id"`
	OrganizationID *int64 `json:"organization_id"`
}

type SubscriptionPlan struct {
	ID            *int64  `json:"id"`
	Code          *string `json:"code"`
	Name          *string `json:"name"`
	StripePriceID *string `json:"stripe_price_id"`
	IsActive      *string `json:"is_active"`
}

type AvilableSubscriptions struct {
	ID          *int64   `json:"id"`
	Code        *string  `json:"code"`
	Name        *string  `json:"name"`
	CostMonthly *float32 `json:"cost_monthly"`
	CostYearly  *float32 `json:"cost_yearly"`
	PriceID     *string  `json:"price_id"`
	IsActive    bool     `json:"is_active"`
}

type SubscriptionsEntitlements struct {
	ID         *int64  `json:"id"`
	PlanID     *string `json:"plan_id"`
	Key        *string `json:"key"`
	LimitValue *int64  `json:"limit_value"`
	Enabled    bool    `json:"enabled"`
	Enterprise bool    `json:"enterprise"`
	IsActive   bool    `json:"is_active"`
}

type OrganizationSubscription struct {
	ID                   *int64     `json:"id"`
	PlanID               *int64     `json:"plan_id"`
	Status               *string    `json:"status"`
	CurrentPeriodStart   *time.Time `json:"current_period_start"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end"`
	CanceledAt           *time.Time `json:"canceled_at"`
	StripeCustomerID     *string    `json:"stripe_customer_id"`
	StripeSubscriptionID *string    `json:"stripe_subscription_id"`
}

type OrganizationPlanEntitlement struct {
	PlanID     *int64  `json:"id"`
	ActionKey  *string `json:"key"`
	LimitValue *int64  `json:"limit_value"`
	Enabled    bool    `json:"enabled"`
	Enterprise bool    `json:"enterprise"`
}
