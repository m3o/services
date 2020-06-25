package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/config"
	logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"

	onboarding "github.com/micro/services/onboarding/proto/onboarding"

	paymentsproto "github.com/micro/services/payments/provider/proto"
)

const storePrefixAccountSecrets = "secrets/"

type tokenToEmail struct {
	Email      string `json:"email"`
	Token      string `json:"token"`
	IsVerified bool   `json:"is_verified"`
}

type Onboarding struct {
	paymentService     paymentsproto.ProviderService
	store              store.Store
	auth               auth.Auth
	sendgridTemplateID string
	sendgridAPIKey     string
	planID             string
}

func NewOnboarding(paymentService paymentsproto.ProviderService,
	store store.Store,
	config config.Config,
	auth auth.Auth) *Onboarding {
	apiKey := config.Get("micro", "onboarding", "sendgrid", "api_key").String("")
	templateID := config.Get("micro", "onboarding", "sendgrid", "template_id").String("")
	planID := config.Get("micro", "onboarding", "plan_id").String("")

	if len(apiKey) == 0 {
		logger.Error("No sendgrid API key provided")
	}
	if len(templateID) == 0 {
		logger.Error("No sendgrid template ID provided")
	}
	return &Onboarding{
		paymentService:     paymentService,
		store:              store,
		auth:               auth,
		sendgridAPIKey:     apiKey,
		sendgridTemplateID: templateID,
		planID:             planID,
	}
}

// SendVerificationEmail is the first step in the onboarding flow.SendVerificationEmail
// A stripe customer and a verification token will be created and an email sent.
func (e *Onboarding) SendVerificationEmail(ctx context.Context,
	req *onboarding.SendVerificationEmailRequest,
	rsp *onboarding.SendVerificationEmailResponse) error {
	logger.Info("Received Onboarding.SendVerificationEmail request")

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
	// @todo send different emails based on if the account already exists
	// ie. registration vs login email.
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

	recs, err := e.store.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.New("Can't verify: record not found")
	} else if err != nil {
		return err
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}
	if tok.Token != req.Token {
		return errors.New("Invalid token")
	}

	// If the user is verified, we are going to log her in and the
	// flow stops here. We return the auth token to be used for further calls.
	if tok.IsVerified {
		secret, err := e.getAccountSecret(req.Email)
		if err != nil {
			return err
		}
		token, err := e.auth.Token(auth.WithCredentials(req.Email, secret))
		if err != nil {
			return err
		}
		// @todo Is this correct?
		rsp.AuthToken = token.RefreshToken
		return nil
	}
	_, err = e.paymentService.CreateCustomer(ctx, &paymentsproto.CreateCustomerRequest{
		Customer: &paymentsproto.Customer{
			Id:   req.Email,
			Type: "user",
		},
	})
	return err
}

func (e *Onboarding) FinishOnboarding(ctx context.Context,
	req *onboarding.FinishOnboardingRequest,
	rsp *onboarding.FinishOnboardingResponse) error {
	logger.Info("Received Onboarding.FinishOnboarding request")

	recs, err := e.store.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.New("Can't verify: record not found")
	} else if err != nil {
		return err
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}
	if tok.Token != req.Token {
		return errors.New("Invalid token")
	}

	_, err = e.paymentService.CreatePaymentMethod(ctx, &paymentsproto.CreatePaymentMethodRequest{
		CustomerId:   req.Email,
		CustomerType: "user",
		Id:           req.PaymentMethodId,
	})
	if err != nil {
		return err
	}

	_, err = e.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
		CustomerId:   req.Email,
		CustomerType: "user",
		PlanId:       e.planID,
	})
	if err != nil {
		return err
	}
	secret := uuid.New().String()
	err = e.setAccountSecret(req.Email, secret)
	if err != nil {
		return err
	}
	_, err = e.auth.Generate(req.Email, auth.WithSecret(secret))
	if err != nil {
		return err
	}
	t, err := e.auth.Token(auth.WithCredentials(req.Email, secret))
	if err != nil {
		return err
	}
	// @todo correct thing to use here?
	rsp.AuthToken = t.RefreshToken
	return nil
}

// lifted from https://github.com/micro/services/blob/550220a6eff2604b3e6d58d09db2b4489967019c/account/web/handler/handler.go#L114
func (e *Onboarding) setAccountSecret(id, secret string) error {
	key := storePrefixAccountSecrets + id
	return e.store.Write(&store.Record{Key: key, Value: []byte(secret)})
}

func (e *Onboarding) getAccountSecret(id string) (string, error) {
	key := storePrefixAccountSecrets + id
	recs, err := e.store.Read(key)
	if err != nil {
		return "", err
	}
	return string(recs[0].Value), nil
}
