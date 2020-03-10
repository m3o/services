package handler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/errors"

	pb "github.com/micro/services/distributed/sprints/proto/sprints"
)

// CreateTask creates a task for a given sprint
func (h *Handler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest, rsp *pb.CreateTaskResponse) error {
	// validate a task was passed to the request
	task := req.Task
	if task == nil {
		return errors.BadRequest(h.name, "Missing task")
	}

	// validate the sprint exists
	if _, err := h.store.ReadSprint(req.SprintId); err != nil {
		return errors.BadRequest(h.name, "Sprint ID invalid, sprint not found")
	}

	// add the auto-generated attributes
	task.Id = uuid.New().String()
	task.Created = time.Now().Unix()

	// write the task to the store
	if err := h.store.CreateTask(req.SprintId, task); err != nil {
		return errors.InternalServerError(h.name, "Unable to create task: %v", err)
	}

	// encode the response
	rsp.Task = task
	return nil
}

// UpdateTask amends a task in the store
func (h *Handler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest, rsp *pb.UpdateTaskResponse) error {
	// validate a task was passed to the request
	if req.Task == nil {
		return errors.BadRequest(h.name, "Missing task")
	}

	// verify the task exists
	task, err := h.store.ReadTask(req.SprintId, req.Task.Id)
	if err != nil {
		return errors.BadRequest(h.name, "Task ID invalid, task not found: %v", err)
	}

	// amend the updatable attributes
	task.Name = req.Task.Name
	task.Status = req.Task.Status

	// update the task in the store
	if err := h.store.UpdateTask(req.SprintId, task); err != nil {
		return errors.InternalServerError(h.name, "Unable to update task: %v", err)
	}

	return nil
}

// DeleteTask removes a task
func (h *Handler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest, rsp *pb.DeleteTaskResponse) error {
	// validate a task id was passed to the request
	if req.SprintId == "" {
		return errors.BadRequest(h.name, "Missing sprint id")
	}
	if req.TaskId == "" {
		return errors.BadRequest(h.name, "Missing task id")
	}

	// delete the task
	if err := h.store.DeleteTask(req.SprintId, req.TaskId); err != nil {
		return errors.InternalServerError(h.name, "Unable to delete task: %v", err)
	}

	return nil
}
