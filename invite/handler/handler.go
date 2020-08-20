package handler

import (
	"context"
	"encoding/json"

	pb "github.com/m3o/services/invite/proto"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service"
	mstore "github.com/micro/micro/v3/service/store"
)

type invite struct {
	Email      string
	Deleted    bool
	Namespaces []string
}

// New returns an initialised handler
func New(srv *service.Service) *Invite {
	return &Invite{
		name: srv.Name(),
	}
}

// Invite implements the invite service inteface
type Invite struct {
	name string
}

// Create an invite
func (h *Invite) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest(h.name, "Unauthorized request")
	}

	// TODO maybe send an email or something
	b, _ := json.Marshal(invite{
		Email:      req.Email,
		Deleted:    false,
		Namespaces: []string{account.Issuer},
	})
	// write the email to the store
	return mstore.Write(&store.Record{
		Key:   req.Email,
		Value: b,
	})
}

// Delete an invite
func (h *Invite) Delete(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// soft delete by marking as deleted. Note, assumes email was present, doesn't error in case it was never created
	b, _ := json.Marshal(invite{Email: req.Email, Deleted: true})
	return mstore.Write(&store.Record{
		Key:   req.Email,
		Value: b,
	})
}

// Validate an invite
func (h *Invite) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	// check if the email exists in the store
	values, err := mstore.Read(req.Email)
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
	rsp.Namespaces = invite.Namespaces
	return nil
}
