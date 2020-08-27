package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	pb "github.com/m3o/services/invite/proto"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/errors"
	merrors "github.com/micro/go-micro/v3/errors"
	logger "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service"
	mconfig "github.com/micro/micro/v3/service/config"
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

// New returns an initialised handler
func New(srv *service.Service) *Invite {
	templateID := mconfig.Get("micro", "invite", "sendgrid", "invite_template_id").String("")
	apiKey := mconfig.Get("micro", "invite", "sendgrid", "api_key").String("")
	emailFrom := mconfig.Get("micro", "invite", "email_from").String("Micro Team <support@micro.mu>")
	testMode := mconfig.Get("micro", "invite", "test_env").Bool(false)

	return &Invite{
		name:             srv.Name(),
		inviteTemplateID: templateID,
		sendgridAPIKey:   apiKey,
		emailFrom:        emailFrom,
		testMode:         testMode,
	}
}

// Invite implements the invite service inteface
type Invite struct {
	name             string
	inviteTemplateID string
	sendgridAPIKey   string
	emailFrom        string
	testMode         bool
}

// Invite a user
// Some cases to think about with this function:
// - a micro admin invites someone to enable signup
// - a user invites a user without sharing namespace ie "hey join micro"
// - a user invites a user to share a namespace ie "hey join my namespace on micro"
func (h *Invite) User(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(h.name, "Unauthorized request")
	}

	namespaces := []string{}
	// When admins invite from "micro", we don't save
	// the namespace because that would enable users to join the
	// micro (admin) namespace which  we do not want.
	if len(req.Namespace) > 0 {
		if account.Issuer == defaultNamespace || account.Issuer == req.Namespace {
			namespaces = append(namespaces, account.Issuer)
		} else {
			return errors.Unauthorized(h.name, "Unauthorized request")
		}
	}
	if account.Issuer != defaultNamespace {
		err := h.canInvite(account.ID, namespaces)
		if err != nil {
			return err
		}
	}
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

	err = h.sendEmail(req.Email, h.inviteTemplateID)
	if err != nil {
		return errors.InternalServerError(h.name, "Failed to send email: %v", err)
	}

	if account.Issuer != defaultNamespace {
		return h.increaseInviteCount(account.ID, namespaces, req.Email)
	}
	return nil
}

func (e *Invite) sendEmail(email, token string) error {
	if e.testMode {
		logger.Infof("Test mode enabled, not sending email to address '%v' ", email)
		return nil
	}
	logger.Infof("Sending email to address '%v'", email)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"template_id": e.inviteTemplateID,
		"from": map[string]string{
			"email": e.emailFrom,
		},
		"personalizations": []interface{}{
			map[string]interface{}{
				"to": []map[string]string{
					{
						"email": email,
					},
				},
				"dynamic_template_data": map[string]string{
					"token": token,
				},
			},
		},
		"mail_settings": map[string]interface{}{
			"sandbox_mode": map[string]bool{
				"enable": e.testMode,
			},
		},
	})

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+e.sendgridAPIKey)
	req.Header.Set("Content-Type", "application/json")
	rsp, err := new(http.Client).Do(req)
	if err != nil {
		logger.Infof("Could not send email, error: %v", err)
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			logger.Errorf("Could not send email, error: %v", err.Error())
			return err
		}
		logger.Errorf("Could not send email, error: %v", string(bytes))
		return merrors.InternalServerError("signup.sendemail", "error sending email")
	}
	return nil
}

// has user invited more than 5 invites sent out already
// || does namespace have more than 5 invite
// -> { forbidden }
func (h *Invite) canInvite(userID string, namespaces []string) error {
	userCounts, err := mstore.Read(path.Join(userCountPrefix, userID), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		return errors.InternalServerError(h.name, "can't read user invite count")
	}
	if len(userCounts) >= maxUserInvites {
		return errors.BadRequest(h.name, "user invite limit reached")
	}

	if len(namespaces) == 0 {
		return nil
	}

	namespaceCounts, err := mstore.Read(path.Join(namespaceCountPrefix, userID), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		return errors.BadRequest(h.name, "can''t read namespace invite count")
	}
	if len(namespaceCounts) >= maxNamespaceInvites {
		return errors.BadRequest(h.name, "user invite limit reached")
	}

	return nil
}

func (h *Invite) increaseInviteCount(userID string, namespaces []string, emailToBeInvited string) error {
	err := mstore.Write(&store.Record{
		Key:   path.Join(userCountPrefix, userID, emailToBeInvited),
		Value: nil,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "can't increase user invite count: %v", err)
	}

	if len(namespaces) == 0 {
		return nil
	}

	err = mstore.Write(&store.Record{
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
