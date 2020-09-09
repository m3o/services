package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro/go-micro/v3/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

// Gitops is the handler
type Gitops struct {
	Branch     string // e.g. production
	Repository string // e.g. m3o/services
}

// WebhookResponse is the response type of the webhook request
type WebhookResponse struct{}

// WebhookRequest is the payload struct sent from GitHub
type WebhookRequest struct {
	After      string     `json:"after"`  // e.g. 4c4eee3fad645d165817ecbec597be6d24685d54
	Before     string     `json:"before"` // e.g. 1cb75bed2ae11fe6c860e4ec2b73ba70f22210de
	Reference  string     `json:"ref"`    // e.g. refs/heads/master
	Repository Repository `json:"repository"`
}

// Repository object sent from GitHub
type Repository struct {
	Name string `json:"full_name"` // e.g. m3o/services
}

// Webhook handles webhooks from GitHub. We use an interface as the request type to ensure no data
// is lost unmarshaling into a struct, as the full message is required in order to verify the hmac
func (g *Gitops) Webhook(ctx context.Context, req json.RawMessage, rsp *WebhookResponse) error {
	fmt.Println(string(req))

	// unmarshal the request in a WebhookRequest object
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
	reqMac, ok := metadata.Get(ctx, "X-Hub-Signature")
	if !ok {
		return errors.Unauthorized("gitops.Webhook", "Missing required header: X-Hub-Signature")
	}
	reqMac = strings.TrimPrefix(reqMac, "sha1=")

	// compare the hmacs
	mac := hmac.New(sha1.New, secret)
	mac.Write(req)
	expectedMAC := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(reqMac), []byte(expectedMAC)) {
		logger.Warnf("HMAC doesn't match")
		// TODO: Urgent, fix the HMAC check.
		// return errors.Unauthorized("gitops.Webhook", "Invalid hmac")
	}

	// check the branch matches
	if payload.Reference != fmt.Sprintf("refs/heads/%v", g.Branch) {
		logger.Infof("Update %v was not on the %v branch", payload.After, g.Branch)
		return nil
	}

	return nil
}
