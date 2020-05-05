package handler

import (
	"context"
	"path"
	"strings"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"

	pb "github.com/micro/services/events/api/proto"
	event "github.com/micro/services/events/service/proto"
	project "github.com/micro/services/projects/service/proto"
)

const (
	githubBase    = "github.com"
	githubPkgBase = "docker.pkg.github.com"
)

// Handler implements the event api interface
type Handler struct {
	name    string
	auth    auth.Auth
	runtime runtime.Runtime
	event   event.EventsService
	project project.ProjectsService
}

// New returns an initialised handler
func New(service micro.Service) *Handler {
	return &Handler{
		name:    service.Name(),
		auth:    service.Options().Auth,
		runtime: runtime.DefaultRuntime,
		event:   event.NewEventsService("go.micro.service.events", service.Client()),
		project: project.NewProjectsService("go.micro.service.projects", service.Client()),
	}
}

// Create a new event
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Metadata == nil {
		return errors.BadRequest(h.name, "Missing metadata")
	}

	// determine the event type
	var evType event.EventType
	switch req.Type {
	case "build_started":
		evType = event.EventType_BuildStarted
	case "build_finished":
		evType = event.EventType_BuildFinished
	case "build_failed":
		evType = event.EventType_BuildFailed
	case "source_created":
		evType = event.EventType_SourceCreated
	case "source_updated":
		evType = event.EventType_SourceUpdated
	case "source_deleted":
		evType = event.EventType_SourceDeleted
	default:
		return errors.BadRequest(h.name, "Invalid type")
	}

	// lookup the account
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(h.name, "account not found")
	}

	// find the namespace the account belongs to
	pRsp, err := h.project.Read(ctx, &project.ReadRequest{Namespace: acc.Namespace})
	if err != nil {
		return err
	}

	// update the runtime
	go h.updateRuntime(acc, evType, req.Metadata, pRsp.Project)

	// create the event
	_, err = h.event.Create(ctx, &event.CreateRequest{
		ProjectId: pRsp.Project.Id,
		Metadata:  req.Metadata,
		Type:      evType,
	})
	return err
}

func (h *Handler) updateRuntime(acc *auth.Account, evType event.EventType, md map[string]string, project *project.Project) {
	// update the runtime. We create a blank context
	// with the account so that the downstream services
	// (e.g. the runtime) will use the namespace only
	// from the account
	ctx := auth.ContextWithAccount(context.Background(), acc)

	// we only care about these two events with regards to the runtime
	if evType != event.EventType_BuildFinished && evType != event.EventType_SourceDeleted {
		return
	}

	// construct the service object
	srvName := md["service"]
	service := &runtime.Service{
		Name:    srvName,
		Source:  path.Join(githubBase, project.Repository, srvName),
		Version: "latest",
		Metadata: map[string]string{
			"commit":      md["commit"],
			"build":       md["build"],
			"repo":        project.Repository,
			"deployed_by": "go.micro.api.events",
		},
	}

	// if the service was deleted, remove it from the runtime
	if evType == event.EventType_SourceDeleted {
		if err := h.runtime.Delete(service, runtime.DeleteContext(ctx)); err != nil {
			logger.Warnf("Failed to delete service %v/%v: %v", acc.Namespace, srvName, err)
		} else {
			logger.Infof("Successfully deleted service %v/%v: %v", acc.Namespace, srvName, err)
		}
		return
	}

	// check if the service is already running, if it is we'll just update it
	srvs, err := h.runtime.Read(
		runtime.ReadService(service.Name),
		runtime.ReadContext(ctx),
	)
	if err != nil {
		logger.Warnf("Failed to read service %v/%v: %v", acc.Namespace, srvName, err)
		return
	}
	if len(srvs) > 0 {
		// the service already exists, we just need to update it
		if err := h.runtime.Update(service, runtime.UpdateContext(ctx)); err != nil {
			logger.Warnf("Failed to update service %v/%v: %v", acc.Namespace, srvName, err)
		} else {
			logger.Warnf("Successfully updated service %v/%v", acc.Namespace, srvName)
		}
		return
	}

	// the service doesn't exist, we must create it
	opts := []runtime.CreateOption{
		runtime.CreateType(typeFromServiceName(srvName)),
		runtime.CreateImage(path.Join(githubPkgBase, project.Repository, srvName)),
		runtime.CreateContext(ctx),
	}
	if err := h.runtime.Create(service, opts...); err != nil {
		logger.Warnf("Failed to create service %v/%v: %v", acc.Namespace, srvName, err)
	} else {
		logger.Warnf("Successfully created service %v/%v", acc.Namespace, srvName)
	}
}

func typeFromServiceName(name string) string {
	if strings.Contains(name, "api") {
		return "api"
	}
	if strings.Contains(name, "web") {
		return "web"
	}
	return "service"
}
