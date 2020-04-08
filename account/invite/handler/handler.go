package handler

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/account/invite/proto"
)

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:  srv.Name(),
		store: store.DefaultStore,
	}
}

// Handler implements the invite service inteface
type Handler struct {
	name  string
	store store.Store
}

// Create an invite
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// generate a six digit code
	for i := 0; i < 6; i++ {
		rsp.Code = rsp.Code + strconv.Itoa(rand.Intn(8)+1)
	}

	// write the code to the store
	return h.store.Write(&store.Record{Key: rsp.Code})
}

// Validate an invite
func (h *Handler) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	// check if the code exists in the store
	_, err := h.store.Read(req.Code)
	if err == store.ErrNotFound {
		return errors.BadRequest(h.name, "Invalid invite code")
	}
	return err
}
