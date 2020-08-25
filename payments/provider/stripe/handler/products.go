package handler

import (
	"context"
	"strings"

	pb "github.com/m3o/services/payments/provider/proto"
	"github.com/micro/go-micro/v3/errors"
	stripe "github.com/stripe/stripe-go/v71"
)

// CreateProduct via the Stripe API, e.g. "Notes"
func (h *Provider) CreateProduct(ctx context.Context, req *pb.CreateProductRequest, rsp *pb.CreateProductResponse) error {
	if req.Product == nil {
		return errors.BadRequest(h.name, "Product required")
	}

	// Construct the stripe product params
	params := &stripe.ProductParams{
		ID:          stripe.String(req.Product.Id),
		Name:        stripe.String(req.Product.Name),
		Description: stripe.String(req.Product.Description),
		Active:      stripe.Bool(req.Product.Active),
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
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}
}

// CreatePlan via the Stripe API, e.g. "Gold"
func (h *Provider) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest, rsp *pb.CreatePlanResponse) error {
	if req.Plan == nil {
		return errors.BadRequest(h.name, "Plan required")
	}

	// Format the interval
	interval := strings.ToLower(req.Plan.Interval.String())

	// Construct the stripe product plan params
	params := &stripe.PlanParams{
		ID:        stripe.String(req.Plan.Id),
		Nickname:  stripe.String(req.Plan.Name),
		Currency:  stripe.String(req.Plan.Currency),
		ProductID: stripe.String(req.Plan.ProductId),
		Interval:  stripe.String(interval),
		Amount:    stripe.Int64(req.Plan.Amount),
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
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}
}

// ListPlans returns all the plans for a product
func (h *Provider) ListPlans(ctx context.Context, req *pb.ListPlansRequest, rsp *pb.ListPlansResponse) error {
	iter := h.client.Plans.List(&stripe.PlanListParams{Product: stripe.String(req.ProductId)})
	if iter.Err() != nil {
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", iter.Err())
	}

	// Loop through and serialize
	rsp.Plans = make([]*pb.Plan, 0)
	for {
		if !iter.Next() {
			break
		}

		pm := serializePlan(iter.Plan())
		rsp.Plans = append(rsp.Plans, pm)
	}

	return nil
}
