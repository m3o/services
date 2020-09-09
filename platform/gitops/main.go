package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/m3o/services/platform/gitops/handler"
	gostore "github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

var (
	// Token is a GitHub PAT for use when registering the webhook
	Token string
	// Repository to watch for changes
	Repository string
	// Branch to filter changes using
	Branch string
	// WebhookURL is the url of the webhook endpoint, e.g. https://mydomain.com/gitops/webhook
	WebhookURL string
)

func main() {
	// create the service
	srv := service.New(
		service.Name("gitops"),
		service.Version("latest"),
	)

	// load the configuration
	loadConfig()

	// register the webhook if it doesn't exist
	if webhookExists() {
		logger.Infof("Webhook already exists for %v", Repository)
	} else {
		logger.Infof("Registering webhook for %v.", Repository)
		registerWebhook()
		logger.Infof("Successfully registed webhook for %v.", Repository)
	}

	srv.Handle(&handler.Gitops{
		Repository: Repository,
		Branch:     Branch,
	})

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

// loadConfig will load the configuration. If a required value if not provided the missing key will
// be logged fatally.
func loadConfig() {
	Token = config.Get("micro", "gitops", "token").String("")
	if len(Token) == 0 {
		logger.Fatalf("Missing required config: micro.gitops.token")
	}

	WebhookURL = config.Get("micro", "gitops", "webhook").String("")
	if len(WebhookURL) == 0 {
		logger.Fatalf("Missing required config: micro.gitops.webhook")
	}

	Branch = config.Get("micro", "gitops", "branch").String("master")
	Repository = config.Get("micro", "gitops", "repository").String("m3o/services")
}

// webhookExists returns a boolean indicating if a webhook has already been registered for the
// current repository. If an error occurs it will be logged fatally.
func webhookExists() bool {
	_, err := store.Read(Repository)
	switch err {
	case nil:
		// the record exists
		return true
	case gostore.ErrNotFound:
		// the record does not exit
		return false
	default:
		// an unknown error occured
		logger.Fatalf("Error reading webhook from store: %v", err)
		return false
	}
}

// registerWebhook registers a GitHub webhook for the current repository. If an error occurs it will
// be logged fatally. If the request succeeds, the webhook secret will be written to the store.
func registerWebhook() {
	// secret will be used as the key to generate the HMAC hex digest value in the X-Hub-Signature
	// header on the webhooks
	secret := uuid.New().String()

	// construct the request body
	// docs: https://docs.github.com/en/rest/reference/repos#create-a-repository-webhook
	data, err := json.Marshal(map[string]interface{}{
		"config": map[string]string{
			"url":          WebhookURL,
			"content_type": "json",
			"secret":       secret,
		},
	})

	// construct the request
	url := fmt.Sprintf("https://api.github.com/repos/%v/hooks", Repository)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logger.Fatalf("Error creating the request: %v", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+Token)

	// execute the request.
	client := new(http.Client)
	rsp, err := client.Do(req)
	if err != nil {
		logger.Fatalf("Error executing create webhook request: %v", err)
	}
	if rsp.StatusCode != http.StatusCreated {
		bytes, _ := ioutil.ReadAll(rsp.Body)
		defer rsp.Body.Close()
		logger.Fatalf("Error creating webhook: %v - %v", rsp.Status, string(bytes))
	}

	// write the secret to the store
	record := &gostore.Record{
		Key:   Repository,
		Value: []byte(secret),
	}
	if err := store.Write(record); err != nil {
		logger.Fatalf("Error writing webhook to the store: %v", err)
	}
}
