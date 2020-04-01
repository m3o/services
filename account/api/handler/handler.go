package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"

	payment "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the account api proto interface
type Handler struct {
	name    string
	auth    auth.Auth
	users   users.UsersService
	payment payment.ProviderService
}

// NewHandler returns an initialised handle
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:    srv.Name(),
		auth:    srv.Options().Auth,
		users:   users.NewUsersService("go.micro.service.users", srv.Client()),
		payment: payment.NewProviderService("go.micro.service.payment.stripe", srv.Client()),
	}
}

func (h *Handler) userFromContext(ctx context.Context) (*users.User, error) {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if len(acc.ID) == 0 {
		return nil, errors.Unauthorized(h.name, "A valid auth token is required")
	}

	// Lookup the user
	resp, err := h.users.Read(ctx, &users.ReadRequest{Email: acc.ID})
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}
