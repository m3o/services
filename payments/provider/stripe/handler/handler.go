package handler

import (
	"encoding/json"
	"log"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	"github.com/stripe/stripe-go/client"
)

// Handler implements the payments provider interface for stripe
type Handler struct {
	name   string      // name of the service
	store  store.Store // go-micro store (key/value)
	client *client.API // stripe api client
}

// NewHandler returns an initialised Handler, it will error if any of
// the required enviroment variables are not set
func NewHandler(srv micro.Service) *Handler {
	apiKey := config.Get("micro", "payments", "stripe", "api_key").String("")
	if len(apiKey) == 0 {
		log.Fatalf("Missing required config: micro.payments.stripe.api_key")
	}

	return &Handler{
		name:   srv.Name(),
		store:  srv.Options().Store,
		client: client.New(apiKey, nil),
	}
}

// Customer is the datatype stored in the store
type Customer struct {
	StripeID string `json:"stripe_id"`
}

// getStripeIDForCustomer returns the stripe ID from the store for the given customer
func (h *Handler) getStripeIDForCustomer(customerType, customerID string) (string, error) {
	recs, err := h.store.Read(customerType + "/" + customerID)
	if err == store.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", errors.InternalServerError(h.name, "Could not read from store: %v", err)
	}

	var c *Customer
	if err := json.Unmarshal(recs[0].Value, &c); err != nil {
		return "", errors.InternalServerError(h.name, "Could not unmarshal json: %v", err)
	}

	return c.StripeID, nil
}

// setStripeIDForCustomer writes the stripe ID to the store for the given customer
func (h *Handler) setStripeIDForCustomer(stripeID, customerType, customerID string) error {
	bytes, err := json.Marshal(&Customer{StripeID: stripeID})
	if err != nil {
		return errors.InternalServerError(h.name, "Could not marshal json: %v", err)
	}

	if err := h.store.Write(&store.Record{Key: customerType + "/" + customerID, Value: bytes}); err != nil {
		return errors.InternalServerError(h.name, "Could not write to store: %v", err)
	}

	return nil
}
