package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/config"
	logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"

	onboarding "github.com/micro/services/onboarding/proto/onboarding"

	paymentsproto "github.com/micro/services/payments/provider/proto"
)

type tokenToEmail struct {
	Email      string `json:"email"`
	Token      string `json:"token"`
	IsVerified bool   `json:"is_verified"`
}

type Onboarding struct {
	paymentService     paymentsproto.ProviderService
	store              store.Store
	sendgridTemplateID string
	sendgridAPIKey     string
}

func NewOnboarding(paymentService paymentsproto.ProviderService,
	store store.Store,
	config config.Config) *Onboarding {
	apiKey := config.Get("micro", "onboarding", "sendgrid", "api_key").String("")
	templateID := config.Get("micro", "onboarding", "sendgrid", "template_id").String("")

	if len(apiKey) == 0 {
		logger.Error("No sendgrid API key provided")
	}
	if len(templateID) == 0 {
		logger.Error("No sendgrid template ID provided")
	}
	return &Onboarding{
		paymentService:     paymentService,
		store:              store,
		sendgridAPIKey:     apiKey,
		sendgridTemplateID: templateID,
	}
}

// SendVerificationEmail is the first step in the onboarding flow.SendVerificationEmail
// A stripe customer and a verification token will be created and an email sent.
func (e *Onboarding) SendVerificationEmail(ctx context.Context,
	req *onboarding.SendVerificationEmailRequest,
	rsp *onboarding.SendVerificationEmailResponse) error {
	logger.Info("Received Onboarding.SendVerificationEmail request")

	// Create a Stripe customer and send a verification email.
	_, err := e.paymentService.CreateCustomer(ctx, &paymentsproto.CreateCustomerRequest{
		Customer: &paymentsproto.Customer{
			Id:   req.Email,
			Type: "user",
		},
	})
	if err != nil {
		return err
	}

	// Save token
	// @todo maybe use something nicer.
	token := uuid.New().String()
	tok := &tokenToEmail{
		Token: token,
		Email: req.Email,
	}
	bytes, err := json.Marshal(tok)
	if err != nil {
		return err
	}

	if err := e.store.Write(&store.Record{Key: req.Email, Value: bytes}); err != nil {
		return err
	}

	// Send email
	err = e.sendEmail(req.Email, token)
	if err != nil {
		return err
	}

	return nil
}

// Lifted  from the invite service https://github.com/micro/services/blob/master/projects/invite/handler/invite.go#L187
// sendEmailInvite sends an email invite via the sendgrid API using the
// predesigned email template. Docs: https://bit.ly/2VYPQD1
func (e *Onboarding) sendEmail(email, token string) error {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"template_id": e.sendgridTemplateID,
		"from": map[string]string{
			"email": "Micro <support@micro.mu>",
		},
		"personalizations": []interface{}{
			map[string]interface{}{
				"to": []map[string]string{
					{
						// @todo consider using proper name
						"name":  strings.Split(email, "@")[0],
						"email": email,
					},
				},
				"dynamic_template_data": map[string]string{
					"token": token,
				},
			},
		},
	})

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+e.sendgridAPIKey)
	req.Header.Set("Content-Type", "application/json")

	if rsp, err := new(http.Client).Do(req); err != nil {
		logger.Info("Could not send email to %v, error: %v", email, err)
	} else if rsp.StatusCode != 202 {
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		logger.Info("Could not send email to %v, error: %v", email, string(bytes))
	}
	return nil
}

func (e *Onboarding) Verify(ctx context.Context,
	req *onboarding.VerifyRequest,
	rsp *onboarding.VerifyResponse) error {
	logger.Info("Received Onboarding.Verify request")
	return nil
}

func (e *Onboarding) FinishOnboarding(ctx context.Context,
	req *onboarding.FinishOnboardingRequest,
	rsp *onboarding.FinishOnboardingResponse) error {
	logger.Info("Received Onboarding.FinishOnboarding request")
	return nil
}
