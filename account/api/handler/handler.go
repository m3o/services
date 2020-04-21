package handler

import (
	"context"

	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"

	invite "github.com/micro/services/account/invite/proto"
	payment "github.com/micro/services/payments/provider/proto"
	teamInvite "github.com/micro/services/teams/invites/proto/invites"
	teams "github.com/micro/services/teams/service/proto/teams"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the account api proto interface
type Handler struct {
	name       string
	auth       auth.Auth
	users      users.UsersService
	teams      teams.TeamsService
	invite     invite.InviteService
	payment    payment.ProviderService
	teamInvite teamInvite.InvitesService
}

// NewHandler returns an initialised handle
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:       srv.Name(),
		auth:       srv.Options().Auth,
		users:      users.NewUsersService("go.micro.service.users", srv.Client()),
		teams:      teams.NewTeamsService("go.micro.service.teams", srv.Client()),
		invite:     invite.NewInviteService("go.micro.service.account.invite", srv.Client()),
		payment:    payment.NewProviderService("go.micro.service.payment.stripe", srv.Client()),
		teamInvite: teamInvite.NewInvitesService("go.micro.service.teams.invites", srv.Client()),
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

// verifyTeamMembership returns a boolean indicating if a user belongs to a team
func (h *Handler) verifyTeamMembership(ctx context.Context, userID, teamID string) bool {
	rsp, err := h.teams.ListMemberships(ctx, &teams.ListMembershipsRequest{MemberId: userID})
	if err != nil {
		logger.Error(err)
		return false
	}

	for _, t := range rsp.Teams {
		if t.Id == teamID {
			return true
		}
	}

	return false
}
