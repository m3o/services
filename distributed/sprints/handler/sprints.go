package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/services/distributed/sprints/proto/sprints"
)

// CreateSprint generates a new sprint
func (h *Handler) CreateSprint(ctx context.Context, req *pb.CreateSprintRequest, rsp *pb.CreateSprintResponse) error {
	// Valiate the request
	var sprint *pb.Sprint
	if req.Sprint != nil {
		sprint = req.Sprint
	}

	// Generate a default name
	if sprint.Name == "" {
		c, err := h.store.CountSprints()
		if err != nil {
			return errors.InternalServerError(h.name, "Unable to count sprints: %v", err)
		}
		sprint.Name = fmt.Sprintf("Sprint #%v", c+1)
	}

	// Generate a default start and end time
	if sprint.StartTime == 0 {
		sprint.StartTime = time.Now().Unix()
	}
	if sprint.EndTime == 0 {
		sprint.EndTime = time.Unix(sprint.StartTime, 0).Add(time.Hour * 24 * 7).Unix()
	}

	// Add the auto-generated attributes
	sprint.Id = uuid.New().String()
	sprint.Created = time.Now().Unix()

	// Write to the store
	if err := h.store.CreateSprint(sprint); err != nil {
		return errors.InternalServerError(h.name, "Unable to create sprint: %v", err)
	}

	// Write the response
	rsp.Sprint = sprint
	rsp.Sprint.Tasks = make([]*pb.Task, 0)
	rsp.Sprint.Objectives = make([]*pb.Objective, 0)
	return nil
}

// ListSprints returns all the sprints recorded in the service, without the nested attributes
func (h *Handler) ListSprints(ctx context.Context, req *pb.ListSprintsRequest, rsp *pb.ListSprintsResponse) error {
	sprints, err := h.store.ListSprints()
	if err != nil {
		return errors.InternalServerError(h.name, "Unable to list sprints: %v", err)
	}
	rsp.Sprints = sprints
	return nil
}

// ReadSprint looks up a sprint using a given ID
func (h *Handler) ReadSprint(ctx context.Context, req *pb.ReadSprintRequest, rsp *pb.ReadSprintResponse) error {
	// Lookup the sprint
	sprint, err := h.store.ReadSprint(req.Id)
	if err != nil {
		return errors.NotFound(h.name, "Sprint not found: %v", err)
	}

	// Lookup the tasks
	tasks, err := h.store.ListTasks(req.Id)
	if err != nil {
		return errors.InternalServerError(h.name, "Could not list tasks: %v", err)
	}

	// Lookup the objectives
	objs, err := h.store.ListObjectives(req.Id)
	if err != nil {
		return errors.InternalServerError(h.name, "Could not list objectives: %v", err)
	}

	// Encode the response
	rsp.Sprint = sprint
	rsp.Sprint.Tasks = tasks
	rsp.Sprint.Objectives = objs
	return nil
}
