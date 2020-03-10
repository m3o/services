package handler

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"

	pb "github.com/micro/services/payments/provider/proto"
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

// CreateProduct via the Stripe API, e.g. "Notes"
func (h *Handler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest, rsp *pb.CreateProductResponse) error {
	if req.Product == nil {
		return errors.BadRequest(h.name, "Product required")
	}

	// Construct the stripe product params
	params := &stripe.ProductParams{
		ID:          &req.Product.Id,
		Name:        &req.Product.Name,
		Description: &req.Product.Description,
		Active:      &req.Product.Active,
	}

	// Create the product
	_, err := h.client.Products.New(params)
	if err == nil {
		return nil
	}

	// Handle the error
	switch err.(*stripe.Error).Code {
	case stripe.ErrorCodeResourceAlreadyExists:
		// the product already exists, update it
		params.ID = nil // don't pass ID again in req body
		_, updateErr := h.client.Products.Update(req.Product.Id, params)
		return updateErr
	default:
		// the error was not expected
		return err
	}
}

// CreatePlan via the Stripe API, e.g. "Gold"
func (h *Handler) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest, rsp *pb.CreatePlanResponse) error {
	if req.Plan == nil {
		return errors.BadRequest(h.name, "Plan required")
	}

	// Format the interval
	interval := strings.ToLower(req.Plan.Interval.String())

	// Construct the stripe product plan params
	params := &stripe.PlanParams{
		ID:        &req.Plan.Id,
		Nickname:  &req.Plan.Name,
		Amount:    &req.Plan.Amount,
		Currency:  &req.Plan.Currency,
		ProductID: &req.Plan.ProductId,
		Interval:  &interval,
	}

	// Create the product plan
	_, err := h.client.Plans.New(params)
	if err == nil {
		return nil
	}

	// Handle the error
	switch err.(*stripe.Error).Code {
	case stripe.ErrorCodeResourceAlreadyExists:
		// the product plan already exists and it cannot be updated
		return nil
	default:
		// the error was not expected
		return err
	}
}

// CreateCustomer via the Stripe API, e.g. "John Doe"
func (h *Handler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest, rsp *pb.CreateCustomerResponse) error {
	return nil
}

// CreateSubscription via the Stripe API, e.g. "Subscribe John Doe to Notes Gold"
func (h *Handler) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest, rsp *pb.CreateSubscriptionResponse) error {
	return nil
}
