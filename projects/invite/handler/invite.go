package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/micro/go-micro/v2/errors"

	"github.com/micro/go-micro/v2/auth"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/projects/invite/proto"
	project "github.com/micro/services/projects/service/proto"
	users "github.com/micro/services/users/service/proto"
)

var (
	// how long should an invite live for
	inviteTTL = time.Hour * 24
)

// invite is written to the store
type invite struct {
	Name        string
	Email       string
	ProjectID   string
	ProjectName string
}

// Invites implements the invites service interface
type Invites struct {
	name               string
	store              store.Store
	sendgridAPIKey     string
	sendgridTemplateID string
	users              users.UsersService
	projects           project.ProjectService
}

// New returns an initialised handler
func New(service micro.Service) *Invites {
	sendgridAPIKey := service.Options().Config.Get("sendgrid-api-key").String("")
	if len(sendgridAPIKey) == 0 {
		logger.Warn("Missing required config: 'sendgrid-api-key'")
	}

	sendgridTemplateID := service.Options().Config.Get("sendgrid-template-id").String("")
	if len(sendgridAPIKey) == 0 {
		logger.Warn("Missing required config: 'sendgrid-template-id'")
	}

	return &Invites{
		name:               service.Name(),
		store:              store.DefaultStore,
		sendgridAPIKey:     sendgridAPIKey,
		sendgridTemplateID: sendgridTemplateID,
		users:              users.NewUsersService("go.micro.service.users", service.Client()),
		projects:           project.NewProjectsService("go.micro.service.projects", service.Client()),
	}
}

// Generate an invite to a user. An email will be sent to this
// user containing a code which is valid for 24 hours.
func (i *Invites) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	// validate the request
	if len(req.ProjectId) == 0 {
		return errors.BadRequest(i.name, "Missing project id")
	}
	if len(req.Email) == 0 {
		return errors.BadRequest(i.name, "Missing email")
	}
	if len(req.Name) == 0 {
		return errors.BadRequest(i.name, "Missing name")
	}

	// get the auth account of the current user (the one sending the invite)
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Forbidden(i.name, "Error loading account")
	}

	// lookup the user using the account id (email)
	uRsp, err := i.users.Read(ctx, &users.ReadRequest{Email: acc.ID})
	if err != nil {
		return errors.InternalServerError(i.name, "Error loading user: %v", err)
	}

	// lookup the project (we need the name)
	pRsp, err := i.projects.Read(ctx, &project.ReadRequest{Id: req.ProjectId})
	if err != nil {
		return err
	}

	// generate the invite and write it to the store
	code := uuid.New().String()
	invite := &invite{Name: req.Name, Email: req.Email, ProjectID: req.ProjectId, ProjectName: pRsp.Project.Name}
	bytes, err := json.Marshal(invite)
	if err != nil {
		return errors.InternalServerError(i.name, "Error mashaling json: %v", err)
	}
	rec := &store.Record{Key: code, Value: bytes, Expiry: inviteTTL}
	if err := i.store.Write(rec); err != nil {
		return errors.InternalServerError(i.name, "Error writing to the store: %v", err)
	}

	// send the email invite async
	go i.sendEmailInvite(req.Name, req.Email, code, pRsp.Project.Name, uRsp.User.FirstName)
	return nil
}

// Verify is called to ensure a code is valid, e.g has not expired.
// This rpc should be called when the user opens the link in their
// email before they create a profile.
func (i *Invites) Verify(ctx context.Context, req *pb.VerifyRequest, rsp *pb.VerifyResponse) error {
	// lookup the invite
	recs, err := i.store.Read(req.Code)
	if err == store.ErrNotFound {
		return errors.BadRequest(i.name, "Invalid invite code")
	} else if err != nil {
		return errors.InternalServerError(i.name, "Error reading from the store: %v", err)
	}

	// unmarshal the invite
	var inv *invite
	if err := json.Unmarshal(recs[0].Value, &inv); err != nil {
		return errors.InternalServerError(i.name, "Error unmashaling json: %v", err)
	}
	rsp.ProjectName = inv.ProjectName
	rsp.Email = inv.Email

	return nil
}

// Redeem is used called after user completes signup and has an account.
// Now they have an account we can redeem the invite and add the user
// to the project. Once this rpc is called, the invite code can no longer
// be used. The email address used when generating the invite must match
// the email of the user redeeming the token.
func (i *Invites) Redeem(ctx context.Context, req *pb.RedeemRequest, rsp *pb.RedeemResponse) error {
	// lookup the invite
	recs, err := i.store.Read(req.Code)
	if err == store.ErrNotFound {
		return errors.BadRequest(i.name, "Invalid invite code")
	} else if err != nil {
		return errors.InternalServerError(i.name, "Error reading from the store: %v", err)
	}

	// unmarshal the record
	var inv *invite
	if err := json.Unmarshal(recs[0].Value, &inv); err != nil {
		return errors.InternalServerError(i.name, "Error unmashaling json: %v", err)
	}

	// lookup the user using the id
	// uRsp, err := i.users.Read(ctx, &users.ReadRequest{Id: req.UserId})
	// if err != nil {
	// 	return err
	// }

	// validate the users email matches the one the invite was sent to
	// if inv.Email != uRsp.User.Email {
	// 	return errors.BadRequest(i.name, "The users email does not match the one invited")
	// }

	// add the user as a project member
	_, err = i.projects.AddMember(ctx, &project.AddMemberRequest{
		Role:      project.Role_Collaborator,
		ProjectId: inv.ProjectID,
		Member: &project.Member{
			Type: "user",
			Id:   req.UserId,
		},
	})
	return err
}

// sendEmailInvite sends an email invite via the sendgrid API using the
// predesigned email template. Docs: https://bit.ly/2VYPQD1
func (i *Invites) sendEmailInvite(name, email, code, project, inviter string) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"template_id": i.sendgridTemplateID,
		"from": map[string]string{
			"email": "Micro <support@micro.mu>",
		},
		"personalizations": []interface{}{
			map[string]interface{}{
				"to": []map[string]string{
					{
						"name":  name,
						"email": email,
					},
				},
				"dynamic_template_data": map[string]string{
					"projectName": project,
					"inviteeName": name,
					"inviterName": inviter,
					"code":        code,
				},
			},
		},
	})

	req, _ := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+i.sendgridAPIKey)
	req.Header.Set("Content-Type", "application/json")

	if rsp, err := new(http.Client).Do(req); err != nil {
		logger.Info("Could not send email to %v, error: %v", email, err)
	} else if rsp.StatusCode != 202 {
		bytes, _ := ioutil.ReadAll(rsp.Body)
		logger.Info("Could not send email to %v, error: %v", email, string(bytes))
	}
}
