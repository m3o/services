package handler

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/account/api/proto/account"
	login "github.com/micro/services/login/service/proto/login"
	payment "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// ReadUser retrieves a user from the users service
func (h *Handler) ReadUser(ctx context.Context, req *pb.ReadUserRequest, rsp *pb.ReadUserResponse) error {
	// Generate a context with elevated privelages
	privCtx, err := auth.ContextWithToken(ctx, h.authToken)
	if err != nil {
		return err
	}

	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Lookup the user
	resp, err := h.users.Read(privCtx, &users.ReadRequest{Id: acc.ID})
	if err != nil {
		return err
	}

	// Serialize the User
	rsp.User = serializeUser(resp.User)
	rsp.User.Roles = make([]string, 0, len(acc.Roles))
	for _, r := range acc.Roles {
		rsp.User.Roles = append(rsp.User.Roles, strings.Title(r))
	}

	// Fetch the payment methods
	pRsp, err := h.payment.ListPaymentMethods(privCtx, &payment.ListPaymentMethodsRequest{UserId: acc.ID})
	if err != nil {
		log.Infof("Error listing payment methods: %v", err)
		return nil
	}

	// Serialize the payment methods
	rsp.User.PaymentMethods = make([]*pb.PaymentMethod, len(pRsp.PaymentMethods))
	for i, p := range pRsp.PaymentMethods {
		rsp.User.PaymentMethods[i] = serializePaymentMethod(p)
	}

	return nil
}

// UpdateUser modifies a user in the users service
func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, rsp *pb.UpdateUserResponse) error {
	// Generate a context with elevated privelages
	privCtx, err := auth.ContextWithToken(ctx, h.authToken)
	if err != nil {
		return err
	}

	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Validate the Userequest
	if req.User == nil {
		return errors.BadRequest(h.name, "User is missing")
	}
	req.User.Id = acc.ID

	// Get the user
	rRsp, err := h.users.Read(privCtx, &users.ReadRequest{Id: acc.ID})
	if err != nil {
		return err
	}

	// Update the user
	uRsp, err := h.users.Update(privCtx, &users.UpdateRequest{User: deserializeUser(req.User)})
	if err != nil {
		return err
	}

	// If the users email changed, notify the login service
	// TODO: Remove this once it's handled by event consumption
	if rRsp.User.Email != uRsp.User.Email {
		h.login.UpdateEmail(ctx, &login.UpdateEmailRequest{
			OldEmail: rRsp.User.Email,
			NewEmail: uRsp.User.Email,
		})
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	return nil
}

// DeleteUser the user service
func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, rsp *pb.DeleteUserResponse) error {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Delete the user
	_, err = h.users.Delete(ctx, &users.DeleteRequest{Id: acc.ID})
	return err
}

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
		Username:  u.Username,
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
		Username:  u.Username,
	}
}
