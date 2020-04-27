package handler

import (
	"context"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"

	pb "github.com/micro/services/event/api/proto"
	event "github.com/micro/services/event/service/proto"
	project "github.com/micro/services/project/service/proto"
)

// Handler implements the event api interface
type Handler struct {
	name    string
	auth    auth.Auth
	event   event.EventService
	project project.ProjectService
}

// New returns an initialised handler
func New(service micro.Service) *Handler {
	return &Handler{
		name:    service.Name(),
		auth:    service.Options().Auth,
		event:   event.NewEventService("go.micro.service.event", service.Client()),
		project: project.NewProjectService("go.micro.service.project", service.Client()),
	}
}

// Create a new event
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// determine the event type
	var eType event.EventType
	switch req.Type {
	case "build_started":
		eType = event.EventType_BuildStarted
	case "build_finished":
		eType = event.EventType_BuildFinished
	case "build_failed":
		eType = event.EventType_BuildFailed
	default:
		return errors.BadRequest(h.name, "Invalid type")
	}

	// lookup the account
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	// find the namespace the account belongs to
	pRsp, err := h.project.Read(ctx, &project.ReadRequest{Namespace: acc.Namespace})
	if err != nil {
		return err
	}

	// create the event
	_, err = h.event.Create(ctx, &event.CreateRequest{
		ProjectId: pRsp.Project.Id,
		Metadata:  req.Metadata,
		Type:      eType,
	})
	return err
}
