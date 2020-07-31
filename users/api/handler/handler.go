package handler

import (
	"context"

	pb "github.com/m3o/services/users/api/proto"
	users "github.com/m3o/services/users/service/proto"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/micro/v3/service"
)

// Handler implements the users api interface
type Handler struct {
	users users.UsersService
}

// NewHandler returns an initialised handler
func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		users: users.NewUsersService("go.micro.service.users"),
	}
}

// Read retrieves a user from the users service
func (h *Handler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// Identify the user
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("go.micro.api.users", "account not found")
	}
	if acc == nil {
		return errors.Unauthorized("go.micro.api.users", "A valid auth token is required")
	}

	// Lookup the user
	resp, err := h.users.Read(ctx, &users.ReadRequest{Id: acc.ID})
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = h.serializeUser(resp.User)
	return nil
}

// Update modifies a user in the users service
func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// Identify the user
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("go.micro.api.users", "account not found")
	}
	if acc == nil {
		return errors.Unauthorized("go.micro.api.users", "A valid auth token is required")
	}

	// Validate the request
	if req.User == nil {
		return errors.BadRequest("go.micro.api.users", "User is missing")
	}
	req.User.Id = acc.ID

	// Update the user
	resp, err := h.users.Update(ctx, &users.UpdateRequest{User: h.deserializeUser(req.User)})
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = h.serializeUser(resp.User)
	return nil
}

// Delete a user in the store
func (h *Handler) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Identify the user
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("go.micro.api.users", "account not found")
	}
	if acc == nil {
		return errors.Unauthorized("go.micro.api.users", "A valid auth token is required")
	}

	// Delete the user
	_, err := h.users.Delete(ctx, &users.DeleteRequest{Id: acc.ID})
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
	}
}
