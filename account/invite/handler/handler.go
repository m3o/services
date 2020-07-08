package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/account/invite/proto"
)

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:  srv.Name(),
		store: srv.Options().Store,
	}
}

// Handler implements the invite service inteface
type Handler struct {
	name  string
	store store.Store
}

// Create an invite
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// TODO maybe send an email or something

	// write the email to the store
	return h.store.Write(&store.Record{
		Key: req.Email,
	})
}

// Validate an invite
func (h *Handler) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	// check if the email exists in the store
	_, err := h.store.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.BadRequest(h.name, "Invalid email")
	}
	return err
}
