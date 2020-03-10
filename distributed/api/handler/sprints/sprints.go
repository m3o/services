package sprints

import (
	"context"

	"github.com/micro/go-micro/v2"
	pb "github.com/micro/services/distributed/api/proto/sprints"
	sprints "github.com/micro/services/distributed/sprints/proto/sprints"
)

// Handler implements the distributed sprints proto
type Handler struct {
	sprints sprints.SprintsService
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		sprints: sprints.NewSprintsService("go.micro.srv.distributed.sprints", srv.Client()),
	}
}

// CreateSprint creates a new sprint in the sprints service
func (h *Handler) CreateSprint(ctx context.Context, req *pb.CreateSprintRequest, rsp *pb.CreateSprintResponse) error {
	r := &sprints.CreateSprintRequest{
		Sprint: deserializeSprint(req.Sprint),
	}

	resp, err := h.sprints.CreateSprint(ctx, r)
	if err != nil {
		return err
	}

	rsp.Sprint = serializeSprint(resp.Sprint)
	return nil
}

// ListSprints returns all the sprints in the sprints service
func (h *Handler) ListSprints(ctx context.Context, req *pb.ListSprintsRequest, rsp *pb.ListSprintsResponse) error {
	resp, err := h.sprints.ListSprints(ctx, &sprints.ListSprintsRequest{})
	if err != nil {
		return err
	}

	rsp.Sprints = make([]*pb.Sprint, len(resp.Sprints))
	for i, s := range resp.Sprints {
		rsp.Sprints[i] = serializeSprint(s)
	}

	return nil
}

// ReadSprint looks up a sprint given an ID
func (h *Handler) ReadSprint(ctx context.Context, req *pb.ReadSprintRequest, rsp *pb.ReadSprintResponse) error {
	r := &sprints.ReadSprintRequest{Id: req.Id}
	resp, err := h.sprints.ReadSprint(ctx, r)
	if err != nil {
		return err
	}

	rsp.Sprint = serializeSprint(resp.Sprint)
	return nil
}

// CreateTask creates a new task in the sprints service for the given sprint
func (h *Handler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest, rsp *pb.CreateTaskResponse) error {
	r := &sprints.CreateTaskRequest{
		Task:     deserializeTask(req.Task),
		SprintId: req.SprintId,
	}

	resp, err := h.sprints.CreateTask(ctx, r)
	if err != nil {
		return err
	}

	rsp.Task = serializeTask(resp.Task)
	return nil
}

// UpdateTask updates a task in the sprints service for the given sprint
func (h *Handler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest, rsp *pb.UpdateTaskResponse) error {
	r := &sprints.UpdateTaskRequest{
		Task:     deserializeTask(req.Task),
		SprintId: req.SprintId,
	}

	_, err := h.sprints.UpdateTask(ctx, r)
	return err
}

// DeleteTask removes a task in the sprints service for the given sprint
func (h *Handler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest, rsp *pb.DeleteTaskResponse) error {
	r := &sprints.DeleteTaskRequest{
		TaskId:   req.TaskId,
		SprintId: req.SprintId,
	}

	_, err := h.sprints.DeleteTask(ctx, r)
	return err
}

// CreateObjective creates a new objective in the sprints service for the given sprint
func (h *Handler) CreateObjective(ctx context.Context, req *pb.CreateObjectiveRequest, rsp *pb.CreateObjectiveResponse) error {
	r := &sprints.CreateObjectiveRequest{
		SprintId:  req.SprintId,
		Objective: deserializeObjective(req.Objective),
	}

	resp, err := h.sprints.CreateObjective(ctx, r)
	if err != nil {
		return err
	}

	rsp.Objective = serializeObjective(resp.Objective)
	return nil
}

// UpdateObjective updates an objective in the sprints service for the given sprint
func (h *Handler) UpdateObjective(ctx context.Context, req *pb.UpdateObjectiveRequest, rsp *pb.UpdateObjectiveResponse) error {
	r := &sprints.UpdateObjectiveRequest{
		SprintId:  req.SprintId,
		Objective: deserializeObjective(req.Objective),
	}

	_, err := h.sprints.UpdateObjective(ctx, r)
	return err
}

// DeleteObjective removes an objective in the sprints service for the given sprint
func (h *Handler) DeleteObjective(ctx context.Context, req *pb.DeleteObjectiveRequest, rsp *pb.DeleteObjectiveResponse) error {
	r := &sprints.DeleteObjectiveRequest{
		SprintId:    req.SprintId,
		ObjectiveId: req.ObjectiveId,
	}

	_, err := h.sprints.DeleteObjective(ctx, r)
	return err
}

func deserializeSprint(s *pb.Sprint) *sprints.Sprint {
	if s == nil {
		return nil
	}

	return &sprints.Sprint{
		Name:      s.Name,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
	}
}

func deserializeTask(s *pb.Task) *sprints.Task {
	if s == nil {
		return nil
	}

	return &sprints.Task{
		Name:   s.Name,
		Status: deserializeStatus(s.Status),
	}
}

func deserializeObjective(o *pb.Objective) *sprints.Objective {
	if o == nil {
		return nil
	}

	return &sprints.Objective{
		Name:   o.Name,
		Status: deserializeStatus(o.Status),
	}
}

func serializeSprint(s *sprints.Sprint) *pb.Sprint {
	if s == nil {
		return nil
	}

	var tasks []*pb.Task
	if s.Tasks != nil {
		tasks = make([]*pb.Task, len(s.Tasks))
		for i, t := range s.Tasks {
			tasks[i] = serializeTask(t)
		}
	}

	var objs []*pb.Objective
	if s.Objectives != nil {
		objs = make([]*pb.Objective, len(s.Objectives))
		for i, j := range s.Objectives {
			objs[i] = serializeObjective(j)
		}
	}

	return &pb.Sprint{
		Id:         s.Id,
		Name:       s.Name,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
		Created:    s.Created,
		Objectives: objs,
		Tasks:      tasks,
	}
}

func serializeObjective(o *sprints.Objective) *pb.Objective {
	if o == nil {
		return nil
	}

	return &pb.Objective{
		Id:      o.Id,
		Name:    o.Name,
		Created: o.Created,
		Status:  serializeStatus(o.Status),
	}
}

func serializeTask(t *sprints.Task) *pb.Task {
	if t == nil {
		return nil
	}

	return &pb.Task{
		Id:      t.Id,
		Name:    t.Name,
		Created: t.Created,
		Status:  serializeStatus(t.Status),
	}
}

func deserializeStatus(s pb.Status) sprints.Status {
	switch s {
	case pb.Status_IN_PROGRESS:
		return sprints.Status_IN_PROGRESS
	case pb.Status_BLOCKED:
		return sprints.Status_BLOCKED
	case pb.Status_COMPLETED:
		return sprints.Status_COMPLETED
	case pb.Status_CANCELLED:
		return sprints.Status_CANCELLED
	default:
		return sprints.Status_PENDING
	}
}

func serializeStatus(s sprints.Status) pb.Status {
	switch s {
	case sprints.Status_IN_PROGRESS:
		return pb.Status_IN_PROGRESS
	case sprints.Status_BLOCKED:
		return pb.Status_BLOCKED
	case sprints.Status_COMPLETED:
		return pb.Status_COMPLETED
	case sprints.Status_CANCELLED:
		return pb.Status_CANCELLED
	default:
		return pb.Status_PENDING
	}
}
