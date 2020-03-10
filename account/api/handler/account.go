package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"

	pb "github.com/micro/services/account/api/proto/account"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the account api proto interface
type Handler struct {
	users users.UsersService
}

// NewHandler returns an initialised handle
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		users: users.NewUsersService("go.micro.srv.users", srv.Client()),
	}
}

// ReadUser retrieves a user from the users service
func (h *Handler) ReadUser(ctx context.Context, req *pb.ReadUserRequest, rsp *pb.ReadUserResponse) error {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if acc == nil {
		return errors.Unauthorized("go.micro.api.users", "A valid auth token is required")
	}

	// Lookup the user
	resp, err := h.users.Read(ctx, &users.ReadRequest{Id: acc.Id})
	if err != nil {
		return err
	}

	// Serialize the Userresponse
	rsp.User = h.serializeUser(resp.User)
	return nil
}

// UpdateUser modifies a user in the users service
func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, rsp *pb.UpdateUserResponse) error {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if acc == nil {
		return errors.Unauthorized("go.micro.api.users", "A valid auth token is required")
	}

	// Validate the Userequest
	if req.User == nil {
		return errors.BadRequest("go.micro.api.users", "User is missing")
	}
	req.User.Id = acc.Id

	// Update the user
	resp, err := h.users.Update(ctx, &users.UpdateRequest{User: h.deserializeUser(req.User)})
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = h.serializeUser(resp.User)
	return nil
}

// DeleteUser the user service
func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, rsp *pb.DeleteUserResponse) error {
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if acc == nil {
		return errors.Unauthorized("go.micro.api.users", "A valid auth token is required")
	}

	// Delete the user
	_, err = h.users.Delete(ctx, &users.DeleteRequest{Id: acc.Id})
	return err
}

func (h *Handler) serializeUser(u *users.User) *pb.User {
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

func (h *Handler) deserializeUser(u *pb.User) *users.User {
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
