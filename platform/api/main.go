package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/platform/api/proto"
	platform "github.com/micro/services/platform/service/proto"
	users "github.com/micro/services/users/service/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.platform"),
	)
	service.Init()

	h := NewHandler(service)
	pb.RegisterPlatformHandler(service.Server(), h)

	if err := service.Run(); err != nil {
		log.Error(err)
	}
}

// Handler is an impementation of the platform api
type Handler struct {
	Platform platform.PlatformService
	Users    users.UsersService
}

// NewHandler returns an initialized Handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		Users:    users.NewUsersService("go.micro.service.users", srv.Client()),
		Platform: platform.NewPlatformService("go.micro.service.platform", srv.Client()),
	}
}

// CreateService deploys a service on the platform
func (h *Handler) CreateService(ctx context.Context, req *pb.CreateServiceRequest, rsp *pb.CreateServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	_, err := h.Platform.CreateService(ctx, &platform.CreateServiceRequest{
		Service: deserializeService(req.Service),
	})

	return err
}

// ReadService returns information about services matching the query
func (h *Handler) ReadService(ctx context.Context, req *pb.ReadServiceRequest, rsp *pb.ReadServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	resp, err := h.Platform.ReadService(ctx, &platform.ReadServiceRequest{
		Service: deserializeService(req.Service),
	})
	if err != nil {
		return err
	}

	rsp.Services = make([]*pb.Service, len(resp.Services))
	for i, s := range resp.Services {
		rsp.Services[i] = serializeService(s)
	}

	return nil
}

// UpdateService updates a service running on the platform
func (h *Handler) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest, rsp *pb.UpdateServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	_, err := h.Platform.UpdateService(ctx, &platform.UpdateServiceRequest{
		Service: deserializeService(req.Service),
	})

	return err
}

// DeleteService terminates a service running on the platform
func (h *Handler) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest, rsp *pb.DeleteServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	_, err := h.Platform.DeleteService(ctx, &platform.DeleteServiceRequest{
		Service: deserializeService(req.Service),
	})

	return err
}

// ListServices returns all the services running on the platform
func (h *Handler) ListServices(ctx context.Context, req *pb.ListServicesRequest, rsp *pb.ListServicesResponse) error {
	resp, err := h.Platform.ListServices(ctx, &platform.ListServicesRequest{})
	if err != nil {
		return err
	}

	rsp.Services = make([]*pb.Service, len(resp.Services))
	for i, s := range resp.Services {
		rsp.Services[i] = serializeService(s)
	}

	return nil
}

// ReadUser gets the current user
func (h *Handler) ReadUser(ctx context.Context, req *pb.ReadUserRequest, rsp *pb.ReadUserResponse) error {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	if acc.Metadata == nil {
		acc.Metadata = make(map[string]string)
	}

	uRsp, err := h.Users.Read(ctx, &users.ReadRequest{Id: acc.ID})
	if err != nil {
		return err
	}

	rsp.User = &pb.User{
		Email:                 uRsp.User.Email,
		Login:                 uRsp.User.FirstName,
		AvatarUrl:             uRsp.User.ProfilePictureUrl,
		Name:                  fmt.Sprintf("%v %v", uRsp.User.FirstName, uRsp.User.LastName),
		TeamName:              "Community",
		TeamUrl:               "https://github.com/orgs/micro/teams/community",
		OrganizationAvatarUrl: "https://avatars3.githubusercontent.com/u/5161210?v=4",
	}

	return nil
}

func serializeService(service *platform.Service) *pb.Service {
	return &pb.Service{
		Name:     service.Name,
		Version:  service.Version,
		Source:   service.Source,
		Metadata: service.Metadata,
		Type:     service.Type,
	}
}

func deserializeService(service *pb.Service) *platform.Service {
	return &platform.Service{
		Name:     service.Name,
		Version:  service.Version,
		Source:   service.Source,
		Metadata: service.Metadata,
		Type:     service.Type,
	}
}
