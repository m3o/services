package handler

import (
	"context"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/account/api/proto/account"
	payment "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// ReadUser retrieves a user from the users service
func (h *Handler) ReadUser(ctx context.Context, req *pb.ReadUserRequest, rsp *pb.ReadUserResponse) error {
	// Get the user
	user, err := h.userFromContext(ctx)
	if err != nil {
		return err
	}

	// Get the account
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	// Serialize the User
	rsp.User = serializeUser(user)
	rsp.User.Roles = acc.Roles

	// Fetch the payment methods
	pRsp, err := h.payment.ListPaymentMethods(ctx, &payment.ListPaymentMethodsRequest{UserId: user.Id})
	if err != nil {
		log.Infof("Error listing payment methods: %v", err)
		return nil
	}

	// Serialize the payment methods
	rsp.User.PaymentMethods = make([]*pb.PaymentMethod, len(pRsp.PaymentMethods))
	for i, p := range pRsp.PaymentMethods {
		rsp.User.PaymentMethods[i] = serializePaymentMethod(p)
	}

	// Fetch the subscriptions
	sRsp, err := h.payment.ListSubscriptions(ctx, &payment.ListSubscriptionsRequest{UserId: user.Id})
	if err != nil {
		log.Infof("Error listing subscriptions: %v", err)
		return nil
	}

	// Serialize the subscriptions
	rsp.User.Subscriptions = make([]*pb.Subscription, len(sRsp.Subscriptions))
	for i, s := range sRsp.Subscriptions {
		rsp.User.Subscriptions[i] = serializeSubscription(s)
	}

	return nil
}

// UpdateUser modifies a user in the users service
func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, rsp *pb.UpdateUserResponse) error {
	// Validate the Userequest
	if req.User == nil {
		return errors.BadRequest(h.name, "User is missing")
	}

	// Get the user
	user, err := h.userFromContext(ctx)
	if err != nil {
		return err
	}
	req.User.Id = user.Id

	// Update the user
	uRsp, err := h.users.Update(ctx, &users.UpdateRequest{User: deserializeUser(req.User)})
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	return nil
}

// DeleteUser the user service
func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, rsp *pb.DeleteUserResponse) error {
	// Get the user
	user, err := h.userFromContext(ctx)
	if err != nil {
		return err
	}

	// Delete the user
	_, err = h.users.Delete(ctx, &users.DeleteRequest{Id: user.Id})
	return err
}

func serializeSubscription(s *payment.Subscription) *pb.Subscription {
	return &pb.Subscription{
		Id: s.Id,
		Plan: &pb.Plan{
			Id:       s.Plan.Id,
			Amount:   s.Plan.Amount,
			Interval: s.Plan.Interval.String(),
		},
	}
}
