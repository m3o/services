package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v3/metadata"
	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/store"
)

// Gitops is the handler
type Gitops struct {
	Branch     string // e.g. production
	Repository string // e.g. m3o/services
}

// Webhook handles webhooks from GitHub. We use an interface as the request type to ensure no data
// is lost unmarshaling into a struct, as the full message is required in order to verify the hmac
func (g *Gitops) Webhook(ctx context.Context, req json.RawMessage, rsp *WebhookResponse) error {
	// unmarshal the request in a webhookRequest object
	var payload WebhookRequest
	if err := json.Unmarshal(req, &payload); err != nil {
		return errors.InternalServerError("gitops.Webhook", "Error unmarshaling request: %v", err)
	}

	// lookup the secret for the repo
	recs, err := store.Read(payload.Repository.Name)
	if err != nil {
		return errors.InternalServerError("gitops.Webhook", "Error reading webhook from store: %v", err)
	}
	secret := recs[0].Value

	// get the hmac from the request
	reqMacHeader, ok := metadata.Get(ctx, "X-Hub-Signature")
	if !ok {
		return errors.Unauthorized("gitops.Webhook", "Missing required header: X-Hub-Signature")
	}

	// split the header and decode the string
	parts := strings.Split(reqMacHeader, "=")
	if len(parts) < 2 {
		return errors.Unauthorized("gitops.Webhook", "Invalid header: X-Hub-Signature")
	}

	// compare the hmacs
	sha, _ := hex.DecodeString(parts[1])
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(req)
	expect := mac.Sum(nil)
	equals := hmac.Equal(sha, expect)
	if !equals {
		return errors.Unauthorized("gitops.Webhook", "Invalid request signature")
	}

	// check the branch matches
	if payload.Reference != fmt.Sprintf("refs/heads/%v", g.Branch) {
		logger.Infof("Update %v was not on the %v branch", payload.After, g.Branch)
		return nil
	}

	// create any new services and delete any removed ones
	changes := determineChanges(payload.Commits)
	for dir, cType := range changes {
		srv := &gorun.Service{
			Name:    filepath.Base(dir),
			Version: "latest",
			Source:  fmt.Sprintf("github.com/%v/%v", g.Repository, dir),
		}

		switch cType {
		case created:
			if err := runtime.Create(srv, gorun.CreateNamespace("micro")); err != nil && err != gorun.ErrAlreadyExists {
				logger.Errorf("Error creating service %v: %v", dir, err)
			} else {
				logger.Infof("Created service %v", srv.Name)
			}
		case deleted:
			if err := runtime.Delete(srv, gorun.DeleteNamespace("micro")); err != nil {
				logger.Errorf("Error deleting service %v: %v", srv.Name, err)
			} else {
				logger.Infof("Deleted service %v", srv.Name)
			}
		}
	}

	// update all other services
	srvs, err := runtime.Read(gorun.ReadNamespace("micro"))
	if err != nil {
		logger.Errorf("Error reading services: %v", err)
		return nil
	}

	for _, srv := range srvs {
		// don't update a service which was just created
		var alreadyAmended bool
		for dir := range changes {
			if filepath.Base(dir) == srv.Name {
				alreadyAmended = true
				break
			}
		}
		if alreadyAmended {
			continue
		}

		if err := runtime.Update(srv, gorun.UpdateNamespace("micro")); err != nil {
			logger.Errorf("Error updating service %v: %v", srv.Name, err)
		} else {
			logger.Infof("Updated service %v", srv.Name)
		}
	}

	return nil
}
