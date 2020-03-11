package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro/v2"
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
	apiKey := os.Getenv("STRIPE_API_KEY")
	if len(apiKey) == 0 {
		log.Fatalf("Missing required env: STRIPE_API_KEY")
	}

	return &Handler{
		store:  store.DefaultStore,
		client: client.New(apiKey, nil),
		name:   srv.Name(),
	}
}

// User is the datatype stored in the store
type User struct {
	StripeID string `json:"stripe_id"`
}

// getStripeIDForUser returns the stripe ID from the store for the given user
func (h *Handler) getStripeIDForUser(userID string) (string, error) {
	recs, err := h.store.Read(userID)
	if err == store.ErrNotFound || len(recs) == 0 {
		return "", nil
	} else if err != nil {
		return "", errors.InternalServerError(h.name, "Could not read from store: %v", err)
	}

	var user *User
	if err := json.Unmarshal(recs[0].Value, &user); err != nil {
		return "", errors.InternalServerError(h.name, "Could not unmarshal json: %v", err)
	}

	fmt.Printf("User #%v has stripe ID: %v\n", userID, user.StripeID)
	return user.StripeID, nil
}

// setStripeIDForUser writes the stripe ID to the store for the given user
func (h *Handler) setStripeIDForUser(stripeID, userID string) error {
	bytes, err := json.Marshal(&User{StripeID: stripeID})
	if err != nil {
		return errors.InternalServerError(h.name, "Could not marshal json: %v", err)
	}

	if err := h.store.Write(&store.Record{Key: userID, Value: bytes}); err != nil {
		return errors.InternalServerError(h.name, "Could not write to store: %v", err)
	}

	return nil
}
