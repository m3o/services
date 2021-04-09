package handler

import (
	"context"
	"encoding/json"

	stripepb "github.com/m3o/services/stripe/proto"
	"github.com/micro/go-micro/v3/logger"
	api "github.com/micro/micro/v3/proto/api"
	log "github.com/micro/micro/v3/service/logger"

	"github.com/stripe/stripe-go/v71"
)

type WebhookResponse struct{}

type Stripe struct{}

func (s *Stripe) IncrementCustomerBalance(ctx context.Context, request *stripepb.IncrementRequest, response *stripepb.IncrementResponse) error {
	panic("implement me")
}

func (s *Stripe) DecrementCustomerBalance(ctx context.Context, request *stripepb.DecrementRequest, response *stripepb.DecrementResponse) error {
	panic("implement me")
}

func (s *Stripe) Webhook(ctx context.Context, req *api.Request, rsp *api.Response) error {

	var ev stripe.Event
	if err := json.Unmarshal([]byte(req.Body), &ev); err != nil {
		log.Errorf("Error unmarshalling Stripe webhook event %s", err)
	}
	log.Infof("Received event %s:%s", ev.ID, ev.Type)
	switch ev.Type {
	case "customer.created":
		return s.customerCreated(ev)
	case "charge.succeeded":
	default:
		log.Infof("Discarding event %s:%s", ev.ID, ev.Type)
	}
	return nil
}

func (s *Stripe) customerCreated(event stripe.Event) error {
	// correlate customer based on email
	// store mapping stripe id to our id
	var cust stripe.Customer
	if err := json.Unmarshal(event.Data.Raw, &cust); err != nil {
		return err
	}
	logger.Infof("Customer created %s", cust.Email)
	return nil
}
