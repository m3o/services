package handler

import (
	"context"
	"encoding/json"
	"fmt"

	custpb "github.com/m3o/services/customers/proto"
	stripepb "github.com/m3o/services/stripe/proto"
	api "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
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
	var err error
	response.NewBalance, err = s.affectBalance(request.CustomerId, -request.Delta)
	return err
}

func (s *Stripe) affectBalance(custID string, delta int64) (int64, error) {
	recs, err := store.Read(fmt.Sprintf(prefixM3OID, custID))
	if err != nil {
		return 0, err
	}
	var cm CustomerMapping
	if err := json.Unmarshal(recs[0].Value, &cm); err != nil {
		return 0, err
	}
	tx, err := bt.New(&stripe.CustomerBalanceTransactionParams{
		Params: stripe.Params{
			IdempotencyKey: nil,
		},
		Amount:   stripe.Int64(delta), //negative is credit
		Customer: stripe.String(cm.StripeID),
	})
	if err != nil {
		return 0, err
	}
	return tx.EndingBalance, nil
}

func (s *Stripe) DecrementCustomerBalance(ctx context.Context, request *stripepb.DecrementRequest, response *stripepb.DecrementResponse) error {
	var err error
	response.NewBalance, err = s.affectBalance(request.CustomerId, request.Delta)
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
	// TODO update their balance
	var cm CustomerMapping
	if err := json.Unmarshal(recs[0].Value, &cm); err != nil {
		return err
	}
	log.Infof("TODO update customer balance for %s", cm.ID)

	return nil
}
