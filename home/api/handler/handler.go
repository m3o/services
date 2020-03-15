package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"

	pb "github.com/micro/services/home/api/proto/home"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the home api interface
type Handler struct {
	name  string
	users users.UsersService
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:  srv.Name(),
		users: users.NewUsersService("go.micro.srv.users", srv.Client()),
	}
}

// ReadUser returns information about the user currently logged in
func (h *Handler) ReadUser(ctx context.Context, req *pb.ReadUserRequest, rsp *pb.ReadUserResponse) error {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	uRsp, err := h.users.Read(ctx, &users.ReadRequest{Id: acc.Id})
	if err != nil {
		return err
	}

	rsp.User = &pb.User{
		FirstName: uRsp.User.FirstName,
		LastName:  uRsp.User.LastName,
	}

	return nil
}
