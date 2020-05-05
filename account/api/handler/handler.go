package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"

	invite "github.com/micro/services/account/invite/proto"
	payment "github.com/micro/services/payments/provider/proto"
	projectInvite "github.com/micro/services/projects/invite/proto"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the account api proto interface
type Handler struct {
	name          string
	auth          auth.Auth
	users         users.UsersService
	invite        invite.InviteService
	payment       payment.ProviderService
	projectInvite projectInvite.InviteService
}

// NewHandler returns an initialised handle
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:          srv.Name(),
		auth:          srv.Options().Auth,
		users:         users.NewUsersService("go.micro.service.users", srv.Client()),
		invite:        invite.NewInviteService("go.micro.service.account.invite", srv.Client()),
		payment:       payment.NewProviderService("go.micro.service.payment.stripe", srv.Client()),
		projectInvite: projectInvite.NewInviteService("go.micro.service.projects.invite", srv.Client()),
	}
}

func (h *Handler) userFromContext(ctx context.Context) (*users.User, error) {
	// Identify the user
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized(h.name, "account not found")
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
