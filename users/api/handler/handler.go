package handler

import (
	"context"
	"strings"

	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/services/users/api/proto"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the users api interface
type Handler struct {
	auth  auth.Auth
	users users.UsersService
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		auth:  srv.Options().Auth,
		users: users.NewUsersService("go.micro.srv.users", srv.Client()),
	}
}

const (
	// BearerScheme is the prefix used in the Authorization header
	BearerScheme = "Bearer "
)

// userFromContext retrives the auth account from the context (req headers).
// TOOD: Refactor this to be part of go-micro/auth
func (h *Handler) userFromContext(ctx context.Context) (*auth.Account, error) {
	// Extract the token if present. Note: if noop is being used
	// then the token can be blank without erroring
	var token string
	if header, ok := metadata.Get(ctx, "Authorization"); ok {
		// Ensure the correct scheme is being used
		if !strings.HasPrefix(header, BearerScheme) {
			return nil, errors.Unauthorized("go.micro.api.users", "invalid authorization header. expected Bearer schema")
		}

		token = header[len(BearerScheme):]
	}

	// Verify the token
	return h.auth.Verify(token)
}

// Read retrieves a user from the users service
func (h *Handler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// Identify the user
	a, err := h.userFromContext(ctx)
	if err != nil {
		return err
	}

	// Lookup the user
	resp, err := h.users.Read(ctx, &users.ReadRequest{Id: a.Id})
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
	a, err := h.userFromContext(ctx)
	if err != nil {
		return err
	}

	// Validate the request
	if req.User == nil {
		return errors.BadRequest("go.micro.api.users", "User is missing")
	}
	req.User.Id = a.Id

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
	a, err := h.userFromContext(ctx)
	if err != nil {
		return err
	}

	// Delete the user
	_, err = h.users.Delete(ctx, &users.DeleteRequest{Id: a.Id})
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
