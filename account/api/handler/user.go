package handler

import (
	"context"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"

	pb "github.com/micro/services/account/api/proto/account"
	invite "github.com/micro/services/account/invite/proto"
	payment "github.com/micro/services/payments/provider/proto"
	teams "github.com/micro/services/teams/service/proto/teams"
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

	// Get the users teams
	tRsp, err := h.teams.ListMemberships(ctx, &teams.ListMembershipsRequest{MemberId: user.Id})
	if err != nil {
		return err
	}
	rsp.User.Teams = make([]*pb.Team, 0, len(tRsp.Teams))
	for _, t := range tRsp.Teams {
		rsp.User.Teams = append(rsp.User.Teams, h.serializeTeam(ctx, t))
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

	// Construct the update params
	updateParams := deserializeUser(req.User)
	updateParams.Id = user.Id

	// Verify the users invite token
	if err := h.verifyInviteToken(ctx, user, req.User.InviteCode); err != nil {
		return err
	}
	updateParams.InviteVerified = true

	// Update the user
	uRsp, err := h.users.Update(ctx, &users.UpdateRequest{User: updateParams})
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

func (h *Handler) verifyInviteToken(ctx context.Context, user *users.User, token string) error {
	if user.InviteVerified {
		return nil
	}
	_, err := h.invite.Validate(ctx, &invite.ValidateRequest{Code: token})
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
