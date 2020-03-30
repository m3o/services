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
	id, err := h.getStripeIDForUser(req.UserId)
	if err != nil {
		return err
	}

	iter := h.client.Subscriptions.List(&stripe.SubscriptionListParams{Customer: id, Plan: req.PlanId})
	if iter.Err() != nil {
		return errors.InternalServerError(h.name, "Unexpected stripe error: %v", iter.Err())
	}

	// Loop through and serialize
	rsp.Subscriptions = make([]*pb.Subscription, 0)
	for {
		if !iter.Next() {
			break
		}

		pm := serializeSubscription(iter.Subscription())
		rsp.Subscriptions = append(rsp.Subscriptions, pm)
	}

	return nil
}

func serializeSubscription(pm *stripe.Subscription) *pb.Subscription {
	rsp := &pb.Subscription{
		Id: pm.ID,
	}

	if pm.Items == nil || len(pm.Items.Data) == 0 {
		return rsp
	}

	plan := pm.Items.Data[0].Plan
	if plan != nil {
		rsp.Plan = serializePlan(plan)
	}
	if plan != nil && plan.Product != nil {
		rsp.Product = serializeProduct(plan.Product)
	}

	return rsp
}

func serializePlan(pm *stripe.Plan) *pb.Plan {
	var interval pb.PlanInterval
	switch pm.Interval {
	case stripe.PlanIntervalDay:
		interval = pb.PlanInterval_DAY
	case stripe.PlanIntervalWeek:
		interval = pb.PlanInterval_WEEK
	case stripe.PlanIntervalMonth:
		interval = pb.PlanInterval_MONTH
	case stripe.PlanIntervalYear:
		interval = pb.PlanInterval_YEAR
	}

	return &pb.Plan{
		Id:       pm.ID,
		Name:     pm.Nickname,
		Amount:   pm.Amount,
		Currency: string(pm.Currency),
		Interval: interval,
	}
}

func serializeProduct(prod *stripe.Product) *pb.Product {
	return &pb.Product{
		Id:   prod.ID,
		Name: prod.Name,
	}
}
