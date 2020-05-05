package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/m3o/api/proto"
	payment "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// Account implments the M3O account service proto
type Account struct {
	name    string
	users   users.UsersService
	payment payment.ProviderService
}

// NewAccount returns an initialised account handler
func NewAccount(service micro.Service) *Account {
	return &Account{
		name:    service.Name(),
		users:   users.NewUsersService("go.micro.service.users", service.Client()),
		payment: payment.NewProviderService("go.micro.service.payment.stripe", service.Client()),
	}
}

// Read the current users info
func (a *Account) Read(ctx context.Context, req *pb.ReadAccountRequest, rsp *pb.ReadAccountResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(a.name, "Account Required")
	}

	uRsp, err := a.users.Read(ctx, &users.ReadRequest{Email: acc.ID})
	if err != nil {
		return err
	}

	rsp.User = serializeUser(uRsp.User)

	// Fetch the payment methods
	pRsp, err := a.payment.ListPaymentMethods(ctx, &payment.ListPaymentMethodsRequest{
		CustomerType: "user", CustomerId: acc.ID,
	})
	if err != nil {
		logger.Infof("Error listing payment methods: %v", err)
		return nil
	}

	// Serialize the payment methods
	rsp.User.PaymentMethods = make([]*pb.PaymentMethod, len(pRsp.PaymentMethods))
	for i, p := range pRsp.PaymentMethods {
		rsp.User.PaymentMethods[i] = serializePaymentMethod(p)
	}

	return nil
}

// CreatePaymentMethod via the provider
func (a *Account) CreatePaymentMethod(ctx context.Context, req *pb.CreatePaymentMethodRequest, rsp *pb.CreatePaymentMethodResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest(a.name, "Missing payment method ID")
	}

	// Identify the user
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(a.name, "A valid auth token is required")
	}

	// Create a payment method
	pRsp, err := a.payment.CreatePaymentMethod(ctx, &payment.CreatePaymentMethodRequest{
		Id: req.Id, CustomerType: "user", CustomerId: acc.ID,
	})
	if err != nil {
		return errors.InternalServerError(a.name, "Error creating payment method: %v", err)
	}

	// Serialize the payment method
	rsp.PaymentMethod = serializePaymentMethod(pRsp.PaymentMethod)

	// Check to see if this is the users only payment method
	lRsp, err := a.payment.ListPaymentMethods(ctx, &payment.ListPaymentMethodsRequest{
		CustomerType: "user", CustomerId: acc.ID,
	})
	if err != nil {
		logger.Infof("Error listing payment methods: %v", err)
		return nil
	}
	if len(lRsp.PaymentMethods) != 1 {
		return nil // no need to set the default
	}

	// Set the default
	_, err = a.payment.SetDefaultPaymentMethod(ctx, &payment.SetDefaultPaymentMethodRequest{
		PaymentMethodId: req.Id, CustomerType: "user", CustomerId: acc.ID,
	})
	if err != nil {
		logger.Infof("Error setting default payment method: %v", err)
		return nil
	}
	rsp.PaymentMethod.Default = true

	return nil
}

// DefaultPaymentMethod sets a users default payment method
func (a *Account) DefaultPaymentMethod(ctx context.Context, req *pb.DefaultPaymentMethodRequest, rsp *pb.DefaultPaymentMethodResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest(a.name, "Missing payment method ID")
	}

	// Identify the user
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(a.name, "A valid auth token is required")
	}

	// Set the default payment method
	_, err := a.payment.SetDefaultPaymentMethod(ctx, &payment.SetDefaultPaymentMethodRequest{
		PaymentMethodId: req.Id, CustomerType: "user", CustomerId: acc.ID,
	})
	if err != nil {
		return errors.InternalServerError(a.name, "Error setting default payment method: %v", err)
	}

	return nil
}

// DeletePaymentMethod via the provider
func (a *Account) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest, rsp *pb.DeletePaymentMethodResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest(a.name, "Missing payment method ID")
	}

	// Delete the payment method
	_, err := a.payment.DeletePaymentMethod(ctx, &payment.DeletePaymentMethodRequest{Id: req.Id})
	if err != nil {
		return errors.InternalServerError(a.name, "Error creating payment method: %v", err)
	}

	return nil
}

func serializeUser(u *users.User) *pb.User {
	return &pb.User{
		Id:                u.Id,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Email:             u.Email,
		ProfilePictureUrl: u.ProfilePictureUrl,
	}
}

func serializePaymentMethod(p *payment.PaymentMethod) *pb.PaymentMethod {
	return &pb.PaymentMethod{
		Id:           p.Id,
		Created:      p.Created,
		CardBrand:    p.CardBrand,
		CardExpMonth: p.CardExpMonth,
		CardExpYear:  p.CardExpYear,
		CardLast_4:   p.CardLast_4,
		Default:      p.Default,
	}
}
