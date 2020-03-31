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

// CreatePaymentMethod via the provider
func (h *Handler) CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest, rsp *pb.CreatePaymentMethodResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest(h.name, "Missing payment method ID")
	}

	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Create a payment method
	pRsp, err := h.payment.CreatePaymentMethod(ctx, &payment.CreatePaymentMethodRequest{UserId: acc.ID, Id: req.Id})
	if err != nil {
		return errors.InternalServerError(h.name, "Error creating payment method: %v", err)
	}

	// Serialize the payment method
	rsp.PaymentMethod = serializePaymentMethod(pRsp.PaymentMethod)

	// Check to see if this is the users only payment method
	lRsp, err := h.payment.ListPaymentMethods(ctx, &payment.ListPaymentMethodsRequest{UserId: acc.ID})
	if err != nil {
		log.Infof("Error listing payment methods: %v", err)
		return nil
	}
	if len(lRsp.PaymentMethods) != 1 {
		return nil // no need to set the default
	}

	// Set the default
	_, err = h.payment.SetDefaultPaymentMethod(ctx, &payment.SetDefaultPaymentMethodRequest{UserId: acc.ID, PaymentMethodId: pRsp.PaymentMethod.Id})
	if err != nil {
		log.Infof("Error setting default payment method: %v", err)
		return nil
	}
	rsp.PaymentMethod.Default = true

	return nil
}

// DefaultPaymentMethod sets a users default payment method
func (h *Handler) DefaultPaymentMethod(ctx context.Context, req *pb.DefaultPaymentMethodRequest, rsp *pb.DefaultPaymentMethodResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest(h.name, "Missing payment method ID")
	}

	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Set the default payment method
	_, err = h.payment.SetDefaultPaymentMethod(ctx, &payment.SetDefaultPaymentMethodRequest{UserId: acc.ID, PaymentMethodId: req.Id})
	if err != nil {
		return errors.InternalServerError(h.name, "Error setting default payment method: %v", err)
	}

	return nil
}

// DeletePaymentMethod via the provider
func (h *Handler) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest, rsp *pb.DeletePaymentMethodResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest(h.name, "Missing payment method ID")
	}

	// Delete the payment method
	_, err := h.payment.DeletePaymentMethod(ctx, &payment.DeletePaymentMethodRequest{Id: req.Id})
	if err != nil {
		return errors.InternalServerError(h.name, "Error creating payment method: %v", err)
	}

	return nil
}

func serializeToken(t *auth.Token) *pb.Token {
	return &pb.Token{
		Token:   t.Token,
		Expires: t.Expiry.Unix(),
	}
}

func serializePaymentMethod(p *payment.PaymentMethod) *pb.PaymentMethod {
	return &pb.PaymentMethod{
		Id:           p.Id,
		Created:      p.Created,
		UserId:       p.UserId,
		Type:         p.Type,
		CardBrand:    p.CardBrand,
		CardExpMonth: p.CardExpMonth,
		CardExpYear:  p.CardExpYear,
		CardLast_4:   p.CardLast_4,
		Default:      p.Default,
	}
}

func serializeUser(u *users.User) *pb.User {
	return &pb.User{
		Id:        u.Id,
		Created:   u.Created,
		Updated:   u.Updated,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

func deserializeUser(u *pb.User) *users.User {
	return &users.User{
		Id:        u.Id,
		Created:   u.Created,
		Updated:   u.Updated,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
