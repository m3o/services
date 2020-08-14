package handler

import (
	"context"
	"log"

	billing "github.com/m3o/services/billing/proto"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/micro/v3/service/config"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
)

type Billing struct {
	stripeClient *client.API // stripe api client
}

func NewBilling() *Billing {
	apiKey := config.Get("micro", "payments", "stripe", "api_key").String("")
	if len(apiKey) == 0 {
		log.Fatalf("Missing required config: micro.payments.stripe.api_key")
	}
	return &Billing{
		stripeClient: client.New(apiKey, nil),
	}
}

// Portal returns the billing portal address the customers can go to to manager their subscriptons
func (b *Billing) Portal(ctx context.Context, req *billing.PortalRequest, rsp *billing.PortalResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("billing.Portal", "Authentication failed")
	}
	params := &stripe.CustomerListParams{
		Email: &acc.ID,
	}
	params.Filters.AddFilter("limit", "", "3")
	customerIter := b.stripeClient.Customers.List(params)

	customerID := ""
	for customerIter.Next() {
		c := customerIter.Customer()
		customerID = c.ID
		break
	}
	if len(customerID) == 0 {
		return errors.BadRequest("billing.Portal", "No stripe customer found for account %v", acc.ID)
	}

	billParams := &stripe.BillingPortalSessionParams{
		Customer: stripe.String(customerID),
	}
	sess, err := b.stripeClient.BillingPortalSessions.New(billParams)
	if err != nil {
		return errors.InternalServerError("billing.Portal", "Could not create billing portal session: %v", err)
	}
	rsp.PortalUrl = sess.URL
	return nil
}
