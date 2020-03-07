package handler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/errors"

	pb "github.com/micro/services/distributed/sprints/proto/sprints"
)

// CreateObjective creates an objective for a given sprint
func (h *Handler) CreateObjective(ctx context.Context, req *pb.CreateObjectiveRequest, rsp *pb.CreateObjectiveResponse) error {
	// validate an objective was passed to the request
	obj := req.Objective
	if obj == nil {
		return errors.BadRequest(h.name, "Missing objective")
	}

	// validate the sprint exists
	if _, err := h.store.ReadSprint(req.SprintId); err != nil {
		return errors.BadRequest(h.name, "Sprint ID invalid, sprint not found")
	}

	// add the auto-generated attributes
	obj.Id = uuid.New().String()
	obj.Created = time.Now().Unix()

	// write the objective to the store
	if err := h.store.CreateObjective(req.SprintId, obj); err != nil {
		return errors.InternalServerError(h.name, "Unable to create objectiv: %v", err)
	}

	// encode the response
	rsp.Objective = obj
	return nil
}

// UpdateObjective amends an objective in the store
func (h *Handler) UpdateObjective(ctx context.Context, req *pb.UpdateObjectiveRequest, rsp *pb.UpdateObjectiveResponse) error {
	// validate an objective was passed to the request
	if req.Objective == nil {
		return errors.BadRequest(h.name, "Missing objective")
	}

	// verify the objective exists
	obj, err := h.store.ReadObjective(req.SprintId, req.Objective.Id)
	if err != nil {
		return errors.BadRequest(h.name, "Objective ID invalid, objective not found: %v", err)
	}

	// amend the updatable attributes
	obj.Name = req.Objective.Name
	obj.Status = req.Objective.Status

	// update the objective in the store
	if err := h.store.UpdateObjective(req.SprintId, obj); err != nil {
		return errors.InternalServerError(h.name, "Unable to update objective: %v", err)
	}

	return nil
}

// DeleteObjective removes an objective
func (h *Handler) DeleteObjective(ctx context.Context, req *pb.DeleteObjectiveRequest, rsp *pb.DeleteObjectiveResponse) error {
	// validate an objective id was passed to the request
	if req.SprintId == "" {
		return errors.BadRequest(h.name, "Missing sprint id")
	}
	if req.ObjectiveId == "" {
		return errors.BadRequest(h.name, "Missing objective id")
	}

	// delete the objective
	if err := h.store.DeleteObjective(req.SprintId, req.ObjectiveId); err != nil {
		return errors.InternalServerError(h.name, "Unable to delete objective: %v", err)
	}

	return nil
}
