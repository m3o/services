package handler

import (
	"context"
	"encoding/json"
	"fmt"

	custpb "github.com/m3o/services/customers/proto"
	stripepb "github.com/m3o/services/stripe/proto"
	api "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/stripe/stripe-go/v71/setupintent"

	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	stripeclient "github.com/stripe/stripe-go/v71/client"
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
	custSvc    custpb.CustomersService
	client     *stripeclient.API // stripe api client
	successURL string
	cancelURL  string
}

func NewHandler(serv *service.Service) stripepb.StripeHandler {
	configObj := struct {
		ApiKey     string `json:"api_key"`
		SuccessURL string `json:"success_url"`
		CancelURL  string `json:"cancel_url"`
	}{}
	val, err := config.Get("micro.stripe")
	if err != nil {
		log.Warnf("Error getting config: %v", err)
	}
	if err := val.Scan(&configObj); err != nil {
		log.Fatalf("Error retrieving config %s", err)
	}

	if len(configObj.ApiKey) == 0 || len(configObj.CancelURL) == 0 || len(configObj.SuccessURL) == 0 {
		log.Fatalf("Missing required config: micro.stripe")
	}

	return &Stripe{
		custSvc:    custpb.NewCustomersService("customers", serv.Client()),
		client:     stripeclient.New(configObj.ApiKey, nil),
		successURL: configObj.SuccessURL,
		cancelURL:  configObj.CancelURL,
	}
}

func (s *Stripe) Webhook(ctx context.Context, ev *stripe.Event, rsp *api.Response) error {
	log.Infof("Received event %s:%s", ev.ID, ev.Type)
	switch ev.Type {
	case "customer.created":
		return s.customerCreated(ctx, ev)
	case "charge.succeeded":
		return s.chargeSucceeded(ctx, ev)
	case "checkout.session.completed":
		return s.checkoutSessionCompleted(ctx, ev)
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

func (s *Stripe) checkoutSessionCompleted(ctx context.Context, event *stripe.Event) error {
	//
	var ch stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &ch); err != nil {
		log.Errorf("Error unmarshalling event %s", err)
		return err
	}
	intent, err := setupintent.Get(ch.SetupIntent.ID, nil)
	if err != nil {
		log.Errorf("Error looking up setup intent %s %s", ch.SetupIntent.ID, err)
		return err
	}
	// if no existing customer create customer and attach
	log.Infof("Intent %+v", intent)
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

func (s *Stripe) CreateCheckoutSession(ctx context.Context, request *stripepb.CreateCheckoutSessionRequest, response *stripepb.CreateCheckoutSessionResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("stripe.CreateCheckoutSession", "Unauthorized")
	}
	if !request.SaveCard && request.Amount < 500 { // min spend
		return errors.BadRequest("stripe.CreateCheckoutSession", "Amount must be at least 500")
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		SuccessURL: stripe.String(s.successURL),
		CancelURL:  stripe.String(s.cancelURL),
	}

	if request.SaveCard {
		params.Mode = stripe.String(string(stripe.CheckoutSessionModeSetup))

	} else {
		params.Mode = stripe.String(string(stripe.CheckoutSessionModeSetup))
		params.LineItems = []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("M3O credit"),
					},
					UnitAmount: stripe.Int64(request.Amount),
				},
				Quantity: stripe.Int64(1),
			}}
	}

	// lookup customer
	recs, err := store.Read(fmt.Sprintf(prefixM3OID, acc.ID))
	if err != nil && err != store.ErrNotFound {
		log.Errorf("Error looking up stripe customer %s", err)
		return errors.InternalServerError("stripe.CreateCheckoutSession", "Error creating checkout session")

	}
	if len(recs) == 0 {
		// use email from account
		params.CustomerEmail = stripe.String(acc.Name)
	} else {
		// use existing customer obj
		var cm CustomerMapping
		json.Unmarshal(recs[0].Value, &cm)
		params.Customer = stripe.String(cm.StripeID)
	}

	session, err := session.New(params)
	if err != nil {
		return err
	}

	response.Id = session.ID
	return nil
}
