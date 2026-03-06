package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"tracker/app/helpers"
	"tracker/app/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/webhook"
)

func (h *AuthHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error found reading body")
		return
	}
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:*")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}
	var model *models.PurchaseIntent
	if err := json.Unmarshal([]byte(body), &model); err != nil {
		http.Error(w, "Unable to unmarshal body", http.StatusInternalServerError)
		return
	}
	// Dev mode, change to prod for domain.
	res, err := h.authService.CreateCheckoutSession(ctx, model, &orgid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": res}
	json.NewEncoder(w).Encode(response)
	return
}

func (h *AuthHandler) CreatePortalSession(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	claims, ok := r.Context().Value("props").(jwt.MapClaims)
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	valid, err := validateRequest(claims, "write:*")
	if err != nil || !valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	orgid, err := helpers.ExtractInt64Claim(claims, "orgid")
	if err != nil {
		http.Error(w, "Unable to parse id", http.StatusInternalServerError)
		return
	}

	// Dev mode, change to prod for domain.
	res, err := h.authService.CreatePortalSession(ctx, &orgid)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
			return
		}
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Printf("service error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": res}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {

	fmt.Println("=== Webhook Debug Info ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("Content-Type: %s\n", r.Header.Get("Content-Type"))
	fmt.Printf("Stripe-Signature: %s\n", r.Header.Get("Stripe-Signature"))

	const tolerance = 300 // seconds
	endpointSecret := os.Getenv("STRIPE_WEBHOOK")
	ctx := r.Context()

	const MaxBodyBytes = int64(65536)
	bodyReader := http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(bodyReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer r.Body.Close()
	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEventWithOptions(payload, signatureHeader, endpointSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		http.Error(w, "Stripe signature error found", http.StatusInternalServerError)
		return
	}
	switch event.Type {
	case "checkout.session.completed":
		// Triggers when a checkout session occurs.
		var sess *stripe.CheckoutSession
		event_id := event.ID
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			http.Error(w, "bad session payload", http.StatusBadRequest)
			return
		}
		if idop, err := h.authService.HasProcessedEvent(ctx, &sess.Customer.ID, &event_id); err != nil && idop {
			fmt.Printf("Event %s already processed. Skipping.\n", event.ID)
			w.WriteHeader(http.StatusOK)
			return
		}
		o := sess.Metadata["orgid"]
		orgid, err := strconv.ParseInt(o, 10, 64)
		if err != nil {
			fmt.Print("Unable to parse")
			http.Error(w, "Unable to link checkout ", http.StatusInternalServerError)
			return
		}
		// Link account to user, no access until payable is complete.
		if err := h.authService.LinkCheckout(sess, &orgid, event_id); err != nil {
			http.Error(w, "Unable to link checkout ", http.StatusInternalServerError)
			return
		}

	case "customer.subscription.created":
		// Triggers when subscription is created.
		var sub stripe.Subscription
		event_id := event.ID
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			http.Error(w, "bad subscription payload", http.StatusBadRequest)
			return
		}
		if idop, err := h.authService.HasProcessedEvent(ctx, &sub.Customer.ID, &event_id); err != nil && idop {
			fmt.Printf("Event %s already processed. Skipping.\n", event.ID)
			w.WriteHeader(http.StatusOK)
			return
		}

		o := sub.Metadata["orgid"]
		orgid, err := strconv.ParseInt(o, 10, 64)
		if err != nil {
			fmt.Print("Unable to parse")
			return
		}
		plan, err := h.authService.GetSubscriptionById(sub.Items.Data[0].Price.ID)
		if err != nil {
			fmt.Print("Error on GetSubscriptionById", err)
			http.Error(w, "Unable to link checkout ", http.StatusInternalServerError)
			return
		}
		fmt.Println("GetSubscriptionById", plan.ID, plan.Code)

		// Grant access to user
		if err := h.authService.GrantAccess(&sub, &orgid, plan); err != nil {
			http.Error(w, "grant access failed", http.StatusInternalServerError)
			return
		}

	case "customer.subscription.updated":
		var sub stripe.Subscription
		event_id := event.ID
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			http.Error(w, "bad subscription payload", http.StatusBadRequest)
			return
		}
		if idop, err := h.authService.HasProcessedEvent(ctx, &sub.Customer.ID, &event_id); err != nil && idop {
			fmt.Printf("Event %s already processed. Skipping.\n", event.ID)
			w.WriteHeader(http.StatusOK)
			return
		}

		admin, err := h.authService.GetOrganizationIdByCustomerId(ctx, sub.Customer.ID)
		if err != nil {
			fmt.Println("Error GetAdminByCustomerID", err)
			http.Error(w, "GetOrganizationIdByCustomerId error ", http.StatusInternalServerError)
			return
		}
		plan, err := h.authService.GetSubscriptionById(sub.Items.Data[0].Price.ID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "GetSubscriptionById error ", http.StatusInternalServerError)
			return
		}
		fmt.Println("GetSubscriptionById", plan.ID, plan.Code)
		// Handle upgrades/downgrades & status transitions here.
		if err := h.authService.UpdateAccess(&sub, admin, plan); err != nil {
			http.Error(w, "grant access failed", http.StatusInternalServerError)
			return
		}

	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			http.Error(w, "bad subscription payload", http.StatusBadRequest)
			return
		}
		admin, err := h.authService.GetOrganizationIdByCustomerId(ctx, sub.Customer.ID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to link checkout ", http.StatusInternalServerError)
			return
		}
		plan, err := h.authService.GetSubscriptionById(sub.Items.Data[0].Price.ID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to link checkout ", http.StatusInternalServerError)
			return
		}
		// Revoke access to user
		if err := h.authService.RevokeAccess(&sub, admin, plan); err != nil {
			http.Error(w, "grant access failed", http.StatusInternalServerError)
			return
		}

	case "invoice.paid":
		var inv stripe.Invoice
		event_id := event.ID
		if err := json.Unmarshal(event.Data.Raw, &inv); err != nil {
			http.Error(w, "bad invoice payload", http.StatusBadRequest)
			return
		}
		if idop, err := h.authService.HasProcessedEvent(ctx, &inv.Customer.ID, &event_id); err != nil && idop {
			fmt.Printf("Event %s already processed. Skipping.\n", event.ID)
			w.WriteHeader(http.StatusOK)
			return
		}
		admin, err := h.authService.GetOrganizationIdByCustomerId(ctx, inv.Customer.ID)
		if err != nil {
			http.Error(w, "Unable to link checkout ", http.StatusInternalServerError)
			return
		}
		//  Verify access to user
		if err := h.authService.VerifyAccess(ctx, &inv, admin); err != nil {
			http.Error(w, "grant access failed", http.StatusInternalServerError)
			return
		}

	case "invoice.payment_failed":
		var inv stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &inv); err != nil {
			http.Error(w, "bad invoice payload", http.StatusBadRequest)
			return
		}
		admin, err := h.authService.GetOrganizationIdByCustomerId(ctx, inv.Customer.ID)
		if err != nil {
			http.Error(w, "Unable to link checkout ", http.StatusInternalServerError)
			return
		}
		if err := h.authService.FailedPayment(&inv, admin); err != nil {
			http.Error(w, "grant access failed", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Errorf("unhandled event type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
