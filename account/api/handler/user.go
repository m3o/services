package handler

import (
	"context"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"

	pb "github.com/micro/services/account/api/proto/account"
	invite "github.com/micro/services/account/invite/proto"
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

func serializeUser(u *users.User) *pb.User {
	return &pb.User{
		Id:             u.Id,
		Created:        u.Created,
		Updated:        u.Updated,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		InviteVerified: u.InviteVerified,
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

func serializeToken(t *auth.Token) *pb.Token {
	return &pb.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Created:      t.Created.Unix(),
		Expiry:       t.Expiry.Unix(),
	}
}
