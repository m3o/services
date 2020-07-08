package handler

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/account/invite/proto"
)

type invite struct {
	Email   string
	Deleted bool
}

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
	b, _ := json.Marshal(invite{Email: req.Email, Deleted: false})
	// write the email to the store
	return h.store.Write(&store.Record{
		Key:   req.Email,
		Value: b,
	})
}

// Delete an invite
func (h *Handler) Delete(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// soft delete by marking as deleted. Note, assumes email was present, doesn't error in case it was never created
	b, _ := json.Marshal(invite{Email: req.Email, Deleted: true})
	return h.store.Write(&store.Record{
		Key:   req.Email,
		Value: b,
	})
}

// Validate an invite
func (h *Handler) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	// check if the email exists in the store
	values, err := h.store.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.BadRequest(h.name, "invalid email")
	} else if err != nil {
		return err
	}
	invite := &invite{}
	if err := json.Unmarshal(values[0].Value, invite); err != nil {
		return err
	}
	if invite.Deleted {
		return errors.BadRequest(h.name, "invalid email")
	}
	return nil
}
