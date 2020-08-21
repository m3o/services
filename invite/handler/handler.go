package handler

import (
	"context"
	"encoding/json"
	"path"

	pb "github.com/m3o/services/invite/proto"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service"
	mstore "github.com/micro/micro/v3/service/store"
)

const (
	// This is defined in internal/namespace/namespace.go so we can't import that
	defaultNamespace = "micro"
	// namespace invite count
	namespaceCountPrefix = "namespace-count"
	// user invite count
	userCountPrefix     = "user-count"
	maxUserInvites      = 5
	maxNamespaceInvites = 5
)

type invite struct {
	Email      string
	Deleted    bool
	Namespaces []string
}

type inviteCount struct {
	Count int
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
// Some cases to think about with this function:
// - a micro admin invites someone to enable signup
// - a user invites a user without sharing namespace ie "hey join micro"
// - a user invites a user to share a namespace ie "hey join my namespace on micro"
func (h *Invite) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(h.name, "Unauthorized request")
	}

	namespaces := []string{}
	// When admins invite from "micro", we don't save
	// the namespace because that would enable users to join the
	// micro (admin) namespace which  we do not want.
	if account.Issuer != defaultNamespace && len(req.Namespace) > 0 {
		namespaces = append(namespaces, account.Issuer)
	}
	if account.Issuer != defaultNamespace {
		err := h.canInvite(account.ID, namespaces)
		if err != nil {
			return err
		}
	}
	// TODO maybe send an email or something
	b, _ := json.Marshal(invite{
		Email:      req.Email,
		Deleted:    false,
		Namespaces: namespaces,
	})
	// write the email to the store
	err := mstore.Write(&store.Record{
		Key:   req.Email,
		Value: b,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "Failed to save invite %v", err)
	}

	if account.Issuer != defaultNamespace {
		return h.increaseInviteCount(account.ID, namespaces)
	}
	return nil
}

// has user invited more than 5 invites sent out already
// || does namespace have more than 5 invite
// -> { forbidden }
func (h *Invite) canInvite(userID string, namespaces []string) error {
	userCounts, err := mstore.Read(path.Join(userCountPrefix, userID))
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError(h.name, "can't read user invite count")
	}
	userCount := &inviteCount{}
	if err == nil {
		if err := json.Unmarshal(userCounts[0].Value, userCount); err != nil {
			return err
		}
	}
	if userCount.Count >= maxUserInvites {
		return errors.BadRequest(h.name, "user invite limit reached")
	}

	if len(namespaces) == 0 {
		return nil
	}

	namespaceCounts, err := mstore.Read(path.Join(namespaceCountPrefix, userID))
	if err != nil && err != store.ErrNotFound {
		return errors.BadRequest(h.name, "can''t read namespace invite count")
	}
	namespaceCount := &inviteCount{}
	if err == nil {
		if err := json.Unmarshal(namespaceCounts[0].Value, userCount); err != nil {
			return err
		}
	}
	if namespaceCount.Count > maxNamespaceInvites {
		return errors.BadRequest(h.name, "user invite limit reached")
	}

	return nil
}

func (h *Invite) increaseInviteCount(userID string, namespaces []string) error {
	userCounts, err := mstore.Read(path.Join(userCountPrefix, userID))
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError(h.name, "can't read user invite count")
	}
	userCount := &inviteCount{}
	if err == nil {
		if err := json.Unmarshal(userCounts[0].Value, userCount); err != nil {
			return err
		}
	}
	userCount.Count++
	b, _ := json.Marshal(userCount)
	err = mstore.Write(&store.Record{
		Key:   path.Join(userCountPrefix, userID),
		Value: b,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "can't increase user invite count: %v", err)
	}

	if len(namespaces) == 0 {
		return nil
	}

	namespaceCounts, err := mstore.Read(path.Join(namespaceCountPrefix, userID))
	if err != nil && err != store.ErrNotFound {
		return errors.BadRequest(h.name, "can''t read namespace invite count")
	}
	namespaceCount := &inviteCount{}
	if err == nil {
		if err := json.Unmarshal(namespaceCounts[0].Value, userCount); err != nil {
			return err
		}
	}
	namespaceCount.Count++
	b, _ = json.Marshal(namespaceCount)
	err = mstore.Write(&store.Record{
		Key:   path.Join(namespaceCountPrefix, namespaces[0]),
		Value: b,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "can't increase namespace invite count: %v", err)
	}
	return nil
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
