package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/micro/go-micro/v3/auth"

	"github.com/micro/micro/v3/service/errors"

	emails "github.com/m3o/services/emails/proto"
	"github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"
)

const (
	defaultIssuer = "micro"
)

func NewEmailsHandler() *Emails {
	enabled := config.Get("micro", "emails", "enabled").Bool(false)
	emailFrom := config.Get("micro", "emails", "email_from").String("Micro Team <support@micro.mu>")
	apiKey := ""
	if enabled {
		apiKey = config.Get("micro", "emails", "sendgrid", "api_key").String("")
		if len(apiKey) == 0 {
			log.Fatalf("Sendgrid API key not configured")
		}
	}
	return &Emails{
		sendingEnabled: enabled,
		apiKey:         apiKey,
		defaultFrom:    emailFrom,
	}
}

type Emails struct {
	sendingEnabled bool
	apiKey         string
	defaultFrom    string
}

func (e *Emails) Send(ctx context.Context, request *emails.SendRequest, response *emails.SendResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok || acc.Issuer != defaultIssuer { // currently only for m3o internal use
		return errors.Unauthorized("emails.send.unauthorized", "Not authorized")
	}
	if len(request.TemplateId) == 0 {
		return errors.BadRequest("emails.send.validation", "Missing template ID")
	}
	if len(request.To) == 0 {
		return errors.BadRequest("emails.send.validation", "Missing to address")
	}

	templateData := map[string]interface{}{}
	if len(request.TemplateData) > 0 {
		if err := json.Unmarshal(request.TemplateData, templateData); err != nil {
			return errors.BadRequest("emails.send.templatedata", "Unable to unmarshal template data")
		}
	}
	if err := e.sendEmail(request.From, request.To, request.TemplateId, templateData); err != nil {
		return errors.InternalServerError("emails.sendemail", "Error sending email")
	}
	return nil
}

// sendEmail sends an email invite via the sendgrid API using the
// pre-designed email template. Docs: https://bit.ly/2VYPQD1
func (e *Emails) sendEmail(from, to, templateID string, templateData map[string]interface{}) error {
	if !e.sendingEnabled {
		masked := to
		if len(to) > 4 {
			masked = masked[:4] + strings.Repeat("*", len(masked[4:]))
		} else {
			masked = strings.Repeat("*", len(masked))
		}
		log.Infof("Sending disabled. Skipping send to %s, template ID %s", masked, templateID)
		return nil
	}
	emailFrom := from
	if len(emailFrom) == 0 {
		emailFrom = e.defaultFrom // TODO only works while this is an internal M3O service
	}
	reqBody, _ := json.Marshal(map[string]interface{}{
		"template_id": templateID,
		"from": map[string]string{
			"email": emailFrom,
		},
		"personalizations": []interface{}{
			map[string]interface{}{
				"to": []map[string]string{
					{
						"email": to,
					},
				},
				"dynamic_template_data": templateData,
			},
		},
		"mail_settings": map[string]interface{}{
			"sandbox_mode": map[string]bool{
				"enable": !e.sendingEnabled,
			},
		},
	})

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")
	rsp, err := new(http.Client).Do(req)
	if err != nil {
		return fmt.Errorf("could not send email, error: %v", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("could not send email, error: %v", string(bytes))
	}
	return nil
}
