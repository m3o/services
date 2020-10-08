package handler

import (
	"encoding/json"
	"log"

	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	mstore "github.com/micro/micro/v3/service/store"

	"github.com/stripe/stripe-go/v71/client"
)

// Provider implements the payments provider interface for stripe
type Provider struct {
	name   string      // name of the service
	store  store.Store // go-micro store (key/value)
	client *client.API // stripe api client
}

// NewProvider returns an initialised Provider, it will error if any of
// the required enviroment variables are not set
func New(srv *service.Service) *Provider {
	val, err := config.Get("micro.payments.stripe.api_key")
	if err != nil {
		logger.Warnf("Error getting config: %v", err)
	}
	apiKey := val.String("")

	if len(apiKey) == 0 {
		log.Fatalf("Missing required config: micro.payments.stripe.api_key")
	}

	return &Provider{
		name:   srv.Name(),
		client: client.New(apiKey, nil),
	}
}

// Customer is the datatype stored in the store
type Customer struct {
	StripeID string `json:"stripe_id"`
}

// getStripeIDForCustomer returns the stripe ID from the store for the given customer
func (h *Provider) getStripeIDForCustomer(customerType, customerID string) (string, error) {
	recs, err := mstore.Read(customerType + "/" + customerID)
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
func (h *Provider) setStripeIDForCustomer(stripeID, customerType, customerID string) error {
	bytes, err := json.Marshal(&Customer{StripeID: stripeID})
	if err != nil {
		return errors.InternalServerError(h.name, "Could not marshal json: %v", err)
	}

	if err := mstore.Write(&store.Record{Key: customerType + "/" + customerID, Value: bytes}); err != nil {
		return errors.InternalServerError(h.name, "Could not write to store: %v", err)
	}

	return nil
}
