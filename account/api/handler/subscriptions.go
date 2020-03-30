package handler

import (
	"context"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/services/account/api/proto/account"
	payment "github.com/micro/services/payments/provider/proto"
)

var (
	// ProductID in stripe
	ProductID = "micro"
)

// ListPlans returns all the available plans
func (h *Handler) ListPlans(ctx context.Context, req *pb.ListPlansRequest, rsp *pb.ListPlansResponse) error {
	pRsp, err := h.payment.ListPlans(ctx, &payment.ListPlansRequest{ProductId: ProductID})
	if err != nil {
		return err
	}

	rsp.Plans = make([]*pb.Plan, 0, len(pRsp.Plans))
	for _, p := range pRsp.Plans {
		rsp.Plans = append(rsp.Plans, serializePlan(p))
	}

	return nil
}

// CreateSubscription for the user
func (h *Handler) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest, rsp *pb.CreateSubscriptionResponse) error {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Create the subscription
	_, err = h.payment.CreateSubscription(ctx, &payment.CreateSubscriptionRequest{UserId: acc.ID, PlanId: req.PlanId})
	return err
}

func serializePlan(p *payment.Plan) *pb.Plan {
	return &pb.Plan{
		Id:        p.Id,
		Name:      p.Name,
		Amount:    p.Amount,
		Interval:  p.Interval.String(),
		Available: p.Id != "team",
	}
}
