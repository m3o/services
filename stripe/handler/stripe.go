package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	custpb "github.com/m3o/services/customers/proto"
	stripepb "github.com/m3o/services/stripe/proto"
	api "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/stripe/stripe-go/v71"
	bt "github.com/stripe/stripe-go/v71/customerbalancetransaction"
)

const (
	prefixStripeID = "mappingByStripeID:%s"
	prefixM3OID    = "mappingByID:%s"
)

type WebhookResponse struct{}

type CustomerMapping struct {
	ID       string
	StripeID string
}

type Stripe struct {
	custSvc custpb.CustomersService
}

func NewHandler(serv *service.Service) stripepb.StripeHandler {
	return &Stripe{
		custSvc: custpb.NewCustomersService("customers", serv.Client()),
	}
}

func (s *Stripe) IncrementCustomerBalance(ctx context.Context, request *stripepb.IncrementRequest, response *stripepb.IncrementResponse) error {
	if err := verifyAdmin(ctx, "stripe.IncrementCustomerBalance"); err != nil {
		return err
	}
	var err error
	response.NewBalance, err = s.affectBalance(request.CustomerId, -request.Delta, request.IdempotencyKey)
	return err
}

func (s *Stripe) affectBalance(custID string, delta int64, idempotencyKey string) (int64, error) {
	recs, err := store.Read(fmt.Sprintf(prefixM3OID, custID))
	if err != nil {
		return 0, err
	}
	var cm CustomerMapping
	if err := json.Unmarshal(recs[0].Value, &cm); err != nil {
		return 0, err
	}
	now := time.Now()
	tx, err := bt.New(&stripe.CustomerBalanceTransactionParams{
		Params: stripe.Params{
			IdempotencyKey: stripe.String(idempotencyKey),
		},
		Amount:   stripe.Int64(delta), //negative is credit
		Customer: stripe.String(cm.StripeID),
		Currency: stripe.String(string(stripe.CurrencyUSD)), // TODO is everything USD?
	})
	log.Infof("Updating balance took %s", time.Since(now))
	if err != nil {
		return 0, err
	}
	return tx.EndingBalance, nil
}

func (s *Stripe) DecrementCustomerBalance(ctx context.Context, request *stripepb.DecrementRequest, response *stripepb.DecrementResponse) error {
	if err := verifyAdmin(ctx, "stripe.DecrementCustomerBalance"); err != nil {
		return err
	}
	var err error
	response.NewBalance, err = s.affectBalance(request.CustomerId, request.Delta, request.IdempotencyKey)
	return err
}

func (s *Stripe) Webhook(ctx context.Context, ev *stripe.Event, rsp *api.Response) error {
	log.Infof("Received event %s:%s", ev.ID, ev.Type)
	switch ev.Type {
	case "customer.created":
		return s.customerCreated(ctx, ev)
	case "charge.succeeded":
		return s.chargeSucceeded(ctx, ev)
	default:
		log.Infof("Discarding event %s:%s", ev.ID, ev.Type)
	}
	return nil
}

func (s *Stripe) customerCreated(ctx context.Context, event *stripe.Event) error {
	// correlate customer based on email
	// store mapping stripe id to our id
	var cust stripe.Customer
	if err := json.Unmarshal(event.Data.Raw, &cust); err != nil {
		return err
	}
	// lookup customer on email

	rsp, err := s.custSvc.Read(ctx, &custpb.ReadRequest{Email: cust.Email}, client.WithAuthToken())
	if err != nil {
		// TODO check if not found error
		log.Errorf("Error looking up customer %s", cust.Email)
		return err
	}
	cm := CustomerMapping{
		ID:       rsp.Customer.Id,
		StripeID: cust.ID,
	}

	// persist it
	return s.storeMapping(&cm)
}

func (s *Stripe) storeMapping(cm *CustomerMapping) error {
	b, _ := json.Marshal(cm)
	// index on both stripe id and our id
	if err := store.Write(
		&store.Record{
			Key:   fmt.Sprintf(prefixM3OID, cm.ID),
			Value: b,
		},
	); err != nil {
		return err
	}
	return store.Write(
		&store.Record{
			Key:   fmt.Sprintf(prefixStripeID, cm.StripeID),
			Value: b,
		},
	)
}

func (s *Stripe) chargeSucceeded(ctx context.Context, event *stripe.Event) error {
	var ch stripe.Charge
	if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
		return err
	}
	// lookup the customer
	recs, err := store.Read(fmt.Sprintf(prefixStripeID, ch.Customer.ID))
	if err != nil {
		if err == store.ErrNotFound {
			log.Errorf("Unrecognised customer for charge %s", ch.ID)
			return nil
		} else {
			log.Errorf("Error looking up customer for charge %s", ch.ID)
		}
		return err
	}
	var cm CustomerMapping
	if err := json.Unmarshal(recs[0].Value, &cm); err != nil {
		return err
	}

	events.Publish("stripe", &stripepb.Event{
		Type: "ChargeSucceeded",
		ChargeSucceeded: &stripepb.ChargeSuceededEvent{
			CustomerId: cm.ID,
			Currency:   string(ch.Currency), // TOOD
			Ammount:    ch.Amount,
		},
	})
	return nil
}

func verifyAdmin(ctx context.Context, method string) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != "micro" {
		return errors.Forbidden(method, "Forbidden")
	}
	for _, s := range acc.Scopes {
		if s == "admin" || s == "service" {
			return nil
		}
	}
	return errors.Forbidden(method, "Forbidden")
}
