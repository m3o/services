package handler

import (
	"context"

	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/services/payments/provider/proto"
	stripe "github.com/stripe/stripe-go"
)

// CreateUser via the Stripe API, e.g. "John Doe"
func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest, rsp *pb.CreateUserResponse) error {
	if req.User == nil {
		return errors.BadRequest(h.name, "User required")
	}
	if req.User.Metadata == nil {
		req.User.Metadata = make(map[string]string, 0)
	}

	// Check to see if the user has already been created
	stripeID, err := h.getStripeIDForUser(req.User.Id)
	if err != nil {
		return err
	}

	// Construct the params
	var params stripe.CustomerParams
	if email := req.User.Metadata["email"]; len(email) > 0 {
		params.Email = stripe.String(email)
	}
	if name := req.User.Metadata["name"]; len(name) > 0 {
		params.Name = stripe.String(name)
	}
	if phone := req.User.Metadata["phone"]; len(phone) > 0 {
		params.Phone = stripe.String(phone)
	}

	// If the user already exists, update using the existing attrbutes
	if len(stripeID) > 0 {
		if _, err := h.client.Customers.Update(stripeID, &params); err != nil {
			return errors.InternalServerError(h.name, "Unexepcted stripe update error: %v", err)
		}
		return nil
	}

	// Create the user in stripe
	c, err := h.client.Customers.New(&params)
	if err != nil {
		return errors.InternalServerError(h.name, "Unexepcted stripe create error: %v", err)
	}

	// Write the ID to the database
	return h.setStripeIDForUser(c.ID, req.User.Id)
}
