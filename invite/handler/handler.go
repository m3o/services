package handler

import (
	"context"
	"encoding/json"
	"path"
	"strings"
	"time"

	eproto "github.com/m3o/services/emails/proto"
	pb "github.com/m3o/services/invite/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	mconfig "github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
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
	Created    int64
}

type sendgridConfig struct {
	InviteTemplateID string `json:"invite_template_id"`
}

type inviteConfig struct {
	EmailFrom string         `json:"email_from"`
	Sendgrid  sendgridConfig `json:"sendgrid"`
}

// New returns an initialised handler
func New(srv *service.Service) *Invite {
	conf := inviteConfig{}
	values, err := mconfig.Get("micro.invite")
	if err != nil {
		logger.Warn(err)
	}
	err = values.Scan(&conf)
	if err != nil {
		logger.Warn(err)
	}

	eSvc := eproto.NewEmailsService("emails", srv.Client())
	return &Invite{
		config:   conf,
		name:     srv.Name(),
		emailSvc: eSvc,
	}
}

// Invite implements the invite service inteface
type Invite struct {
	config   inviteConfig
	name     string
	emailSvc eproto.EmailsService
}

// Invite a user
// Some cases to think about with this function:
// - a micro admin invites someone to enable signup
// - a user invites a user without sharing namespace ie "hey join micro"
// - a user invites a user to share a namespace ie "hey join my namespace on micro"
func (h *Invite) User(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return errors.BadRequest("invite.user.validation", "Valid email is required")
	}

	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(h.name, "Unauthorized request")
	}

	namespaces := []string{}

	// When admins invite from "micro", we don't save
	// the namespace because that would enable users to join the
	// micro (admin) namespace which  we do not want.
	if ns := req.Namespace; len(ns) > 0 {
		// if its an admin making the request or the namespace matches then append
		if account.Issuer == defaultNamespace || account.Issuer == req.Namespace {
			// append the requested namespace
			namespaces = append(namespaces, ns)
		} else {
			return errors.Unauthorized(h.name, "Unauthorized request")
		}
	}

	// check for the invite limit
	if account.Issuer != defaultNamespace {
		err := h.canInvite(account.ID, namespaces)
		if err != nil {
			return err
		}
	}

	if !req.Resend {
		// in normal circumstances we want to check in case we've sent an invite already
		recs, err := mstore.Read(req.Email)
		if err != nil && err != mstore.ErrNotFound {
			return err
		}
		if len(recs) > 0 {
			inv := &invite{}
			if err := json.Unmarshal(recs[0].Value, inv); err != nil {
				return err
			}
			if !inv.Deleted {
				// we've already sent one and it's not been deleted, return success
				logger.Infof("Invite already sent for user %s. Skipping ", req.Email)
				return nil
			}
		}
	}
	b, _ := json.Marshal(invite{
		Email:      req.Email,
		Deleted:    false,
		Namespaces: namespaces,
		Created:    time.Now().Unix(),
	})
	// write the email to the store
	err := mstore.Write(&mstore.Record{
		Key:   req.Email,
		Value: b,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "Failed to save invite %v", err)
	}

	err = h.sendEmail(ctx, req.Email, h.config.Sendgrid.InviteTemplateID)
	if err != nil {
		return errors.InternalServerError(h.name, "Failed to send email: %v", err)
	}

	if account.Issuer != defaultNamespace {
		return h.increaseInviteCount(account.ID, namespaces, req.Email)
	}
	return nil
}

func (e *Invite) sendEmail(ctx context.Context, email, token string) error {
	templateData := map[string]string{
		"token": token,
	}
	b, _ := json.Marshal(templateData)

	_, err := e.emailSvc.Send(ctx, &eproto.SendRequest{From: e.config.EmailFrom, To: email, TemplateId: e.config.Sendgrid.InviteTemplateID, TemplateData: b}, client.WithAuthToken())
	return err
}

// has user invited more than 5 invites sent out already
// || does namespace have more than 5 invite
// -> { forbidden }
func (h *Invite) canInvite(userID string, namespaces []string) error {
	userCounts, err := mstore.Read("", mstore.Prefix(path.Join(userCountPrefix, userID)))
	if err != nil && err != mstore.ErrNotFound {
		return errors.InternalServerError(h.name, "can't read user invite count")
	}
	if len(userCounts) >= maxUserInvites {
		return errors.BadRequest(h.name, "user invite limit reached")
	}

	if len(namespaces) == 0 {
		return nil
	}

	namespaceCounts, err := mstore.Read("", mstore.Prefix(path.Join(namespaceCountPrefix, userID)))
	if err != nil && err != mstore.ErrNotFound {
		return errors.BadRequest(h.name, "can''t read namespace invite count")
	}
	if len(namespaceCounts) >= maxNamespaceInvites {
		return errors.BadRequest(h.name, "user invite limit reached")
	}

	return nil
}

func (h *Invite) increaseInviteCount(userID string, namespaces []string, emailToBeInvited string) error {
	err := mstore.Write(&mstore.Record{
		Key:   path.Join(userCountPrefix, userID, emailToBeInvited),
		Value: nil,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "can't increase user invite count: %v", err)
	}

	if len(namespaces) == 0 {
		return nil
	}

	err = mstore.Write(&mstore.Record{
		Key:   path.Join(namespaceCountPrefix, namespaces[0], emailToBeInvited),
		Value: nil,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "can't increase namespace invite count: %v", err)
	}
	return nil
}

// Delete an invite
func (h *Invite) Delete(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(h.name, "Unauthorized request")
	}
	if account.Issuer != defaultNamespace {
		return errors.Unauthorized(h.name, "Unauthorized request")
	}

	// soft delete by marking as deleted. Note, assumes email was present, doesn't error in case it was never created
	b, _ := json.Marshal(invite{Email: req.Email, Deleted: true})
	return mstore.Write(&mstore.Record{
		Key:   req.Email,
		Value: b,
	})
}

// Validate an invite
func (h *Invite) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	// check if the email exists in the store
	values, err := mstore.Read(req.Email)
	if err == mstore.ErrNotFound {
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
