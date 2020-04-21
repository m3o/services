package handler

import (
	"context"

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
	// Get the user
	user, err := h.userFromContext(ctx)
	if err != nil {
		return err
	}

	// Validate the user has access to the team
	if !h.verifyTeamMembership(ctx, user.Id, req.TeamId) {
		return errors.Forbidden(h.name, "Forbidden team")
	}

	// Create the subscription
	sReq := &payment.CreateSubscriptionRequest{PlanId: req.PlanId, CustomerType: "team", CustomerId: req.TeamId}
	_, err = h.payment.CreateSubscription(ctx, sReq)
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
