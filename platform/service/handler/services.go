package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/runtime"
	pb "github.com/micro/services/platform/service/proto"
)

// CreateService deploys a service on the platform
func (h *Handler) CreateService(ctx context.Context, req *pb.CreateServiceRequest, rsp *pb.CreateServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	service := deserializeService(req.Service)

	// specify the image to use
	formattedName := strings.ReplaceAll(service.Name, "/", "-")
	image := fmt.Sprintf("%v/%v", Image, formattedName)

	opts := []runtime.CreateOption{
		// create a specific service type
		runtime.CreateType(req.Service.Type),
		runtime.CreateImage(image),
	}

	if err := h.Runtime.Create(service, opts...); err != nil {
		return err
	}

	go h.Event.Publish(ctx, &pb.Event{
		Type:    pb.EventType_ServiceCreated,
		Service: req.Service,
	})

	return nil
}

// ReadService returns information about services matching the query
func (h *Handler) ReadService(ctx context.Context, req *pb.ReadServiceRequest, rsp *pb.ReadServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	resp, err := h.Runtime.Read(
		runtime.ReadType(req.Service.Type),
		runtime.ReadService(req.Service.Name),
		runtime.ReadVersion(req.Service.Version),
	)
	if err != nil {
		return err
	}

	rsp.Services = make([]*pb.Service, len(resp))
	for i, s := range resp {
		rsp.Services[i] = serializeService(s)
	}

	return nil
}

// UpdateService updates a service running on the platform
func (h *Handler) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest, rsp *pb.UpdateServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	go h.Event.Publish(ctx, &pb.Event{
		Type:    pb.EventType_ServiceUpdated,
		Service: req.Service,
	})

	return h.Runtime.Update(deserializeService(req.Service))
}

// DeleteService terminates a service running on the platform
func (h *Handler) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest, rsp *pb.DeleteServiceResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.api.platform", "service required")
	}

	go h.Event.Publish(ctx, &pb.Event{
		Type:    pb.EventType_ServiceDeleted,
		Service: req.Service,
	})

	return h.Runtime.Delete(deserializeService(req.Service))
}

// ListServices returns all the services running on the platform
func (h *Handler) ListServices(ctx context.Context, req *pb.ListServicesRequest, rsp *pb.ListServicesResponse) error {
	resp, err := h.Runtime.List()
	fmt.Println("whaatt", resp, err)
	if err != nil {
		return err
	}

	rsp.Services = make([]*pb.Service, len(resp))
	for i, s := range resp {
		rsp.Services[i] = serializeService(s)
	}

	return nil
}

func serializeService(srv *runtime.Service) *pb.Service {
	return &pb.Service{
		Name:     srv.Name,
		Version:  srv.Version,
		Source:   srv.Source,
		Metadata: srv.Metadata,
	}
}

func deserializeService(srv *pb.Service) *runtime.Service {
	return &runtime.Service{
		Name:     srv.Name,
		Version:  srv.Version,
		Source:   srv.Source,
		Metadata: srv.Metadata,
	}
}
