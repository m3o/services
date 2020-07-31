package handler

import (
	"context"
	"path"
	"strings"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service"
	mruntime "github.com/micro/micro/v3/service/runtime"

	pb "github.com/m3o/services/events/api/proto"
	event "github.com/m3o/services/events/service/proto"
	environments "github.com/m3o/services/projects/environments/proto"
	project "github.com/m3o/services/projects/service/proto"
)

const (
	githubBase    = "github.com"
	githubPkgBase = "docker.pkg.github.com"
)

// Handler implements the event api interface
type Handler struct {
	name         string
	event        event.EventsService
	project      project.ProjectsService
	environments environments.EnvironmentsService
}

// New returns an initialised handler
func New(service *service.Service) *Handler {
	return &Handler{
		name:         service.Name(),
		event:        event.NewEventsService("go.micro.service.events"),
		project:      project.NewProjectsService("go.micro.service.projects"),
		environments: environments.NewEnvironmentsService("go.micro.service.projects.environments"),
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
	if acc.Metadata == nil || len(acc.Metadata["project-id"]) == 0 {
		return errors.Unauthorized(h.name, "Invalid account used, missing project-id metadata")
	}

	// find the project the account belongs to
	pRsp, err := h.project.Read(ctx, &project.ReadRequest{Id: acc.Metadata["project-id"]})
	if err != nil {
		return err
	}

	// get the projects environments
	eRsp, err := h.environments.Read(ctx, &environments.ReadRequest{ProjectId: acc.Metadata["project-id"]})
	if err != nil {
		return err
	}

	// update the runtime
	for _, env := range eRsp.Environments {
		// create the event
		_, err = h.event.Create(ctx, &event.CreateRequest{
			EnvironmentId: env.Id,
			Metadata:      req.Metadata,
			Type:          evType,
		})
		if err != nil {
			return err
		}

		if err := h.updateRuntime(evType, req.Metadata, pRsp.Project, env); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) updateRuntime(evType event.EventType, md map[string]string, project *project.Project, env *environments.Environment) error {
	// we only care about these two events with regards to the runtime
	if evType != event.EventType_BuildFinished && evType != event.EventType_SourceDeleted {
		return nil
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
		opts := []runtime.DeleteOption{
			runtime.DeleteNamespace(env.Namespace),
		}

		if err := mruntime.Delete(service, opts...); err != nil {
			logger.Warnf("Failed to delete service %v/%v: %v", env.Namespace, srvName, err)
			return err
		}
		logger.Infof("Successfully deleted service %v/%v", env.Namespace, srvName)
	}

	// check if the service is already running, if it is we'll just update it
	srvs, err := mruntime.Read(
		runtime.ReadService(service.Name),
		runtime.ReadNamespace(env.Namespace),
	)
	if err != nil {
		logger.Warnf("Failed to read service %v/%v: %v", env.Namespace, srvName, err)
		return err
	}
	if len(srvs) > 0 {
		opts := []runtime.UpdateOption{
			runtime.UpdateNamespace(env.Namespace),
		}

		// the service already exists, we just need to update it
		if err := mruntime.Update(service, opts...); err != nil {
			logger.Warnf("Failed to update service %v/%v: %v", env.Namespace, srvName, err)
			return err
		}
		logger.Warnf("Successfully updated service %v/%v", env.Namespace, srvName)
	}

	image := path.Join(githubPkgBase, project.Repository, srvName)
	logger.Infof("Namespace: %v; Image: %v\n", env.Namespace, image)

	// the service doesn't exist, we must create it
	opts := []runtime.CreateOption{
		runtime.CreateImage(image),
		runtime.CreateType(typeFromServiceName(srvName)),
		runtime.CreateNamespace(env.Namespace),
	}
	if err := mruntime.Create(service, opts...); err != nil {
		logger.Warnf("Failed to create service %v/%v: %v", env.Namespace, srvName, err)
		return err
	}
	logger.Warnf("Successfully created service %v/%v", env.Namespace, srvName)

	return nil
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
