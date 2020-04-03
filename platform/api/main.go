package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"

	logproto "github.com/micro/micro/v2/debug/log/proto"
	statsproto "github.com/micro/micro/v2/debug/stats/proto"
	traceproto "github.com/micro/micro/v2/debug/trace/proto"
	pb "github.com/micro/services/platform/api/proto"
	platform "github.com/micro/services/platform/service/proto"
	users "github.com/micro/services/users/service/proto"
)

const (
	name = "go.micro.api.platform"
)

func main() {
	service := micro.NewService(
		micro.Name(name),
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
	service  micro.Service
}

// NewHandler returns an initialized Handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		service:  srv,
		Users:    users.NewUsersService("go.micro.service.users", srv.Client()),
		Platform: platform.NewPlatformService("go.micro.service.platform", srv.Client()),
	}
}

// CreateService deploys a service on the platform
func (h *Handler) CreateService(ctx context.Context, req *pb.CreateServiceRequest, rsp *pb.CreateServiceResponse) error {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(name, "A valid auth token is required")
	}
	if req.Service == nil {
		return errors.BadRequest(name, "service required")
	}

	_, err = h.Platform.CreateService(ctx, &platform.CreateServiceRequest{
		Service: deserializeService(req.Service),
	})

	return err
}

// ReadService returns information about services matching the query
func (h *Handler) ReadService(ctx context.Context, req *pb.ReadServiceRequest, rsp *pb.ReadServiceResponse) error {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(name, "A valid auth token is required")
	}
	if req.Service == nil {
		return errors.BadRequest(name, "service required")
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
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(name, "A valid auth token is required")
	}
	if req.Service == nil {
		return errors.BadRequest(name, "service required")
	}

	_, err = h.Platform.UpdateService(ctx, &platform.UpdateServiceRequest{
		Service: deserializeService(req.Service),
	})

	return err
}

// DeleteService terminates a service running on the platform
func (h *Handler) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest, rsp *pb.DeleteServiceResponse) error {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(name, "A valid auth token is required")
	}
	if req.Service == nil {
		return errors.BadRequest(name, "service required")
	}

	_, err = h.Platform.DeleteService(ctx, &platform.DeleteServiceRequest{
		Service: deserializeService(req.Service),
	})

	return err
}

// ListServices returns all the services running on the platform
func (h *Handler) ListServices(ctx context.Context, req *pb.ListServicesRequest, rsp *pb.ListServicesResponse) error {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(name, "A valid auth token is required")
	}
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
	// Identify the user
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	if len(acc.ID) == 0 {
		return errors.Unauthorized(name, "A valid auth token is required")
	}

	// Lookup the user
	uRsp, err := h.Users.Read(ctx, &users.ReadRequest{Email: acc.ID})
	if err != nil {
		return err
	}

	if acc.Metadata == nil {
		acc.Metadata = make(map[string]string)
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

func (h *Handler) Logs(ctx context.Context, req *pb.LogsRequest, rsp *pb.LogsResponse) error {
	if len(req.GetService()) == 0 {
		return errors.BadRequest(name, "Service missing")
	}
	client := h.service.Client()
	request := client.NewRequest("go.micro.debug", "Log.Read", &logproto.ReadRequest{
		Service: req.GetService(),
	})
	resp := &logproto.ReadResponse{}
	if err := client.Call(ctx, request, resp); err != nil {
		return err
	}
	rsp.Records = make([]*pb.Record, 0, len(resp.GetRecords()))
	for _, v := range resp.GetRecords() {
		rsp.Records = append(rsp.Records, &pb.Record{
			Timestamp: v.GetTimestamp(),
			Metadata:  v.GetMetadata(),
		})
	}
	return nil
}

func (h *Handler) Stats(ctx context.Context, req *pb.StatsRequest, rsp *pb.StatsResponse) error {
	if len(req.GetService().GetName()) == 0 {
		return errors.BadRequest(name, "Service missing")
	}
	client := h.service.Client()
	preq := &statsproto.ReadRequest{
		Service: &statsproto.Service{
			Name: req.GetService().GetName(),
		},
		Past: true,
	}
	version := req.GetService().GetVersion()
	if len(version) > 0 {
		preq.Service.Version = version
	}
	request := client.NewRequest("go.micro.debug", "Stats.Read", preq)
	resp := &statsproto.ReadResponse{}
	if err := client.Call(ctx, request, resp); err != nil {
		return err
	}
	rsp.Stats = make([]*pb.StatSnapshot, 0, len(resp.Stats))
	for _, v := range resp.GetStats() {
		rsp.Stats = append(rsp.Stats, &pb.StatSnapshot{
			Service: &pb.Service{
				Name:    v.GetService().GetName(),
				Version: v.GetService().GetVersion(),
			},
			Started:   v.GetStarted(),
			Uptime:    v.GetUptime(),
			Memory:    v.GetMemory(),
			Threads:   v.GetThreads(),
			Gc:        v.GetGc(),
			Requests:  v.GetRequests(),
			Errors:    v.GetErrors(),
			Timestamp: v.GetTimestamp(),
		})
	}
	return nil
}

func (h *Handler) Traces(ctx context.Context, req *pb.TracesRequest, rsp *pb.TracesResponse) error {
	reqProto := &traceproto.ReadRequest{
		Past: true,
	}
	var limit int64 = 1000
	if req.GetLimit() > 0 {
		limit = req.GetLimit()
	}
	if len(req.GetService().GetName()) > 0 {
		reqProto.Service = &traceproto.Service{
			Name: req.GetService().GetName(),
		}
		reqProto.Limit = limit
	}
	client := h.service.Client()
	request := client.NewRequest("go.micro.debug", "Trace.Read", reqProto)
	resp := &traceproto.ReadResponse{}
	if err := client.Call(ctx, request, resp); err != nil {
		return err
	}
	rsp.Spans = make([]*pb.Span, 0, len(resp.GetSpans()))
	for _, v := range resp.GetSpans() {
		rsp.Spans = append(rsp.Spans, &pb.Span{
			Trace:    v.GetTrace(),
			Id:       v.GetId(),
			Parent:   v.GetParent(),
			Name:     v.GetName(),
			Started:  v.GetStarted(),
			Duration: v.GetDuration(),
			Metadata: v.GetMetadata(),
			Type:     pb.SpanType(v.GetType()),
		})
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
