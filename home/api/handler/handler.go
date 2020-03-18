package handler

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"

	apps "github.com/micro/services/apps/service/proto/apps"
	pb "github.com/micro/services/home/api/proto/home"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the home api interface
type Handler struct {
	name  string
	apps  apps.AppsService
	users users.UsersService
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:  srv.Name(),
		apps:  apps.NewAppsService("go.micro.service.apps", srv.Client()),
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

// ListApps returns all the apps a user has access to
func (h *Handler) ListApps(ctx context.Context, req *pb.ListAppsRequest, rsp *pb.ListAppsResponse) error {
	_, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	aRsp, err := h.apps.List(ctx, &apps.ListRequest{OnlyActive: true})
	if err != nil {
		return err
	}

	rsp.Apps = make([]*pb.App, len(aRsp.Apps))
	for i, a := range aRsp.Apps {
		// Asset are served from root, e.g.icon.png
		// would become /distributed/icon.png
		var icon string
		if len(a.Icon) > 0 {
			icon = fmt.Sprintf("/%v/%v", a.Id, a.Icon)
		}

		rsp.Apps[i] = &pb.App{
			Id:       a.Id,
			Name:     a.Name,
			Category: a.Category,
			Icon:     icon,
		}
	}

	return nil
}
