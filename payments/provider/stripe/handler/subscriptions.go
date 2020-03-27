package handler

import (
	"context"

	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/services/payments/provider/proto"
	stripe "github.com/stripe/stripe-go"
)

// CreateSubscription via the Stripe API, e.g. "Subscribe John Doe to Notes Gold"
func (h *Handler) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest, rsp *pb.CreateSubscriptionResponse) error {
	id, err := h.getStripeIDForUser(req.UserId)
	if err != nil {
		return err
	}

	_, err = h.client.Subscriptions.New(&stripe.SubscriptionParams{
		Customer: stripe.String(id),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(req.PlanId),
			},
		},
	})
	if err == nil {
		return nil
	}

	// Handle the error
	switch err.(*stripe.Error).Code {
	case stripe.ErrorCodeParameterInvalidEmpty:
		return errors.BadRequest("go.micro.service.payment.stripe", "missing arguments")
	default:
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", err)
	}
}

func (h *Handler) ListSubscriptions(ctx context.Context, req *pb.ListSubscriptionsRequest, rsp *pb.ListSubscriptionsResponse) error {
	return nil
}
