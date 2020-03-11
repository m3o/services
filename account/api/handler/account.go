package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/account/api/proto/account"
	"github.com/micro/services/payments/provider"
	payment "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the account api proto interface
type Handler struct {
	name    string
	users   users.UsersService
	payment payment.ProviderService
}

// NewHandler returns an initialised handle
func NewHandler(srv micro.Service) *Handler {
	pay, err := provider.NewProvider("stripe", srv.Client())
	if err != nil {
		log.Fatalf("Error setting up payment provider: %v", err)
	}

	return &Handler{
		name:    srv.Name(),
		users:   users.NewUsersService("go.micro.srv.users", srv.Client()),
		payment: pay,
	}
}

// ReadUser retrieves a user from the users service
func (h *Handler) ReadUser(ctx context.Context, req *pb.ReadUserRequest, rsp *pb.ReadUserResponse) error {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.Id) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Lookup the user
	resp, err := h.users.Read(ctx, &users.ReadRequest{Id: acc.Id})
	if err != nil {
		return err
	}

	// Serialize the User
	rsp.User = serializeUser(resp.User)

	// Fetch the payment methods
	pRsp, err := h.payment.ListPaymentMethods(ctx, &payment.ListPaymentMethodsRequest{UserId: acc.Id})
	if err != nil {
		return errors.InternalServerError(h.name, "Error listing payment methods: %v", err)
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
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.Id) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Validate the Userequest
	if req.User == nil {
		return errors.BadRequest(h.name, "User is missing")
	}
	req.User.Id = acc.Id

	// Update the user
	resp, err := h.users.Update(ctx, &users.UpdateRequest{User: deserializeUser(req.User)})
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(resp.User)
	return nil
}

// DeleteUser the user service
func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, rsp *pb.DeleteUserResponse) error {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.Id) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Delete the user
	_, err = h.users.Delete(ctx, &users.DeleteRequest{Id: acc.Id})
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
	if len(acc.Id) == 0 {
		return errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Create a payment method
	pRsp, err := h.payment.CreatePaymentMethod(ctx, &payment.CreatePaymentMethodRequest{UserId: acc.Id, Id: req.Id})
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
