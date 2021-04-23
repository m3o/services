package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	aproto "github.com/m3o/services/alert/proto/alert"
	cproto "github.com/m3o/services/customers/proto"
	eproto "github.com/m3o/services/emails/proto"
	inviteproto "github.com/m3o/services/invite/proto"
	nproto "github.com/m3o/services/namespaces/proto"
	onboarding "github.com/m3o/services/onboarding/proto"
	pproto "github.com/m3o/services/payments/proto"
	sproto "github.com/m3o/services/subscriptions/proto"
	authproto "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	mconfig "github.com/micro/micro/v3/service/config"
	cont "github.com/micro/micro/v3/service/context"
	merrors "github.com/micro/micro/v3/service/errors"
	logger "github.com/micro/micro/v3/service/logger"
	model "github.com/micro/micro/v3/service/model"
	mstore "github.com/micro/micro/v3/service/store"
)

const (
	namespace          = "backend"
	internalErrorMsg   = "An error occurred during onboarding. Contact #m3o-support at slack.m3o.com if the issue persists"
	notInvitedErrorMsg = "You have not been invited to the service. Please request an invite on m3o.com"
)

const (
	expiryDuration      = 5 * time.Minute
	prefixPaymentMethod = "payment-method/"

	onboardingTopic = "onboarding"
)

type tokenToEmail struct {
	Email      string `json:"email"`
	Token      string `json:"token"`
	Created    int64  `json:"created"`
	CustomerID string `json:"customerID"`
}

type Signup struct {
	inviteService       inviteproto.InviteService
	customerService     cproto.CustomersService
	namespaceService    nproto.NamespacesService
	subscriptionService sproto.SubscriptionsService
	alertService        aproto.AlertService
	paymentService      pproto.ProviderService
	emailService        eproto.EmailsService
	auth                auth.Auth
	accounts            authproto.AccountsService
	config              conf
	cache               *cache.Cache
	resetCode           model.Model
}

var (
	// TODO: move this message to a better location
	// Message is a predefined message returned during onboarding
	Message = "Please complete onboarding at https://m3o.com/subscribe?email=%s. This command will now wait for you to finish."
)

type ResetToken struct {
	Created int64
	ID      string
	Token   string
}

type sendgridConf struct {
	TemplateID         string `json:"template_id"`
	RecoveryTemplateID string `json:"recovery_template_id"`
}

type conf struct {
	TestMode       bool         `json:"test_env"`
	PaymentMessage string       `json:"message"`
	Sendgrid       sendgridConf `json:"sendgrid"`
	// using a negative "nopayment" rather than "paymentrequired" because it will default to having to pay if not set
	NoPayment bool `json:"no_payment"`
}

func NewSignup(srv *service.Service, auth auth.Auth) *Signup {
	c := conf{}
	val, err := mconfig.Get("micro.onboarding")
	if err != nil {
		logger.Warnf("Error getting config: %v", err)
	}
	err = val.Scan(&c)
	if err != nil {
		logger.Warnf("Error scanning config: %v", err)
	}

	if len(strings.TrimSpace(c.PaymentMessage)) == 0 {
		c.PaymentMessage = Message
	}
	if !c.TestMode && len(c.Sendgrid.TemplateID) == 0 {
		logger.Warnf("No sendgrid template ID provided")
	}

	s := &Signup{
		inviteService:       inviteproto.NewInviteService("invite", srv.Client()),
		customerService:     cproto.NewCustomersService("customers", srv.Client()),
		namespaceService:    nproto.NewNamespacesService("namespaces", srv.Client()),
		subscriptionService: sproto.NewSubscriptionsService("subscriptions", srv.Client()),
		paymentService:      pproto.NewProviderService("payments", srv.Client()),
		emailService:        eproto.NewEmailsService("emails", srv.Client()),
		auth:                auth,
		accounts:            authproto.NewAccountsService("auth", srv.Client()),
		config:              c,
		cache:               cache.New(1*time.Minute, 5*time.Minute),
		alertService:        aproto.NewAlertService("alert", srv.Client()),
		resetCode:           model.New(ResetToken{}, nil),
	}
	return s
}

// taken from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// SendVerificationEmail is the first step in the onboarding flow.SendVerificationEmail
// A stripe customer and a verification token will be created and an email sent.
func (e *Signup) SendVerificationEmail(ctx context.Context,
	req *onboarding.SendVerificationEmailRequest,
	rsp *onboarding.SendVerificationEmailResponse) error {
	err := e.sendVerificationEmail(ctx, req, rsp)
	if err != nil {
		logger.Warnf("Error during sending verification email: %v", err)
	}
	return err
}

func (e *Signup) sendVerificationEmail(ctx context.Context,
	req *onboarding.SendVerificationEmailRequest,
	rsp *onboarding.SendVerificationEmailResponse) error {
	logger.Info("Received Signup.SendVerificationEmail request")

	k := randStringBytesMaskImprSrc(8)
	tok := &tokenToEmail{
		Token:      k,
		Email:      req.Email,
		Created:    time.Now().Unix(),
		CustomerID: uuid.New().String(),
	}

	bytes, err := json.Marshal(tok)
	if err != nil {
		logger.Error(err)
		return merrors.InternalServerError("onboarding.SendVerificationEmail", internalErrorMsg)
	}

	if err := mstore.Write(&mstore.Record{
		Key:   req.Email,
		Value: bytes,
	}); err != nil {
		logger.Error(err)
		return merrors.InternalServerError("onboarding.SendVerificationEmail", internalErrorMsg)
	}
	// HasPaymentMethod needs to resolve email from token, so we save the
	// same record under a token too
	if err := mstore.Write(&mstore.Record{
		Key:   tok.Token,
		Value: bytes,
	}); err != nil {
		logger.Error(err)
		return merrors.InternalServerError("onboarding.SendVerificationEmail", internalErrorMsg)
	}

	if e.config.TestMode {
		logger.Infof("Sending verification token '%v'", k)
	}

	// Send email
	// @todo send different emails based on if the account already exists
	// ie. registration vs login email.

	err = e.sendEmail(ctx, req.Email, e.config.Sendgrid.TemplateID, map[string]interface{}{
		"token": k,
	})
	if err != nil {
		logger.Errorf("Error when sending email to %v: %v", req.Email, err)
		return merrors.InternalServerError("onboarding.SendVerificationEmail", internalErrorMsg)
	}

	return nil
}

// Lifted  from the invite service https://github.com/m3o/services/blob/master/projects/invite/handler/invite.go#L187
// sendEmailInvite sends an email invite via the sendgrid API using the
// predesigned email template. Docs: https://bit.ly/2VYPQD1
func (e *Signup) sendEmail(ctx context.Context, email, templateID string, templateData map[string]interface{}) error {
	b, _ := json.Marshal(templateData)
	_, err := e.emailService.Send(context.TODO(), &eproto.SendRequest{To: email, TemplateId: templateID, TemplateData: b}, client.WithAuthToken())
	return err
}

func (e *Signup) Verify(ctx context.Context, req *onboarding.VerifyRequest, rsp *onboarding.VerifyResponse) error {
	err := e.verify(ctx, req, rsp)
	if err != nil {
		logger.Warnf("Error during verification: %v", err)
	}
	return err
}

func (e *Signup) verify(ctx context.Context, req *onboarding.VerifyRequest, rsp *onboarding.VerifyResponse) error {
	logger.Info("Received Signup.Verify request")

	recs, err := mstore.Read(req.Email)
	if err == mstore.ErrNotFound {
		logger.Errorf("Can't verify, record for %v is not found", req.Email)
		return merrors.InternalServerError("onboarding.Verify", internalErrorMsg)
	} else if err != nil {
		logger.Errorf("email verification error: %v", err)
		return merrors.InternalServerError("onboarding.Verify", internalErrorMsg)
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}

	if tok.Token != req.Token {
		return merrors.Forbidden("onboarding.Verify", "The token you provided is invalid")
	}

	if time.Since(time.Unix(tok.Created, 0)) > expiryDuration {
		return merrors.Forbidden("onboarding.Verify", "The token you provided has expired")
	}

	// set the response message
	rsp.Message = fmt.Sprintf(e.config.PaymentMessage, req.Email)

	rsp.Namespaces = []string{namespace}
	rsp.CustomerID = tok.CustomerID

	return nil
}

func (e *Signup) CompleteSignup(ctx context.Context, req *onboarding.CompleteSignupRequest, rsp *onboarding.CompleteSignupResponse) error {
	err := e.completeSignup(ctx, req, rsp)
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (e *Signup) completeSignup(ctx context.Context, req *onboarding.CompleteSignupRequest, rsp *onboarding.CompleteSignupResponse) error {
	logger.Info("Received Signup.CompleteSignup request")

	recs, err := mstore.Read(req.Email)
	if err == mstore.ErrNotFound {
		logger.Errorf("Can't verify record for %v: record not found", req.Email)
		return merrors.InternalServerError("onboarding.CompleteSignup", internalErrorMsg)
	} else if err != nil {
		logger.Errorf("Error reading store: err")
		return merrors.InternalServerError("onboarding.CompleteSignup", internalErrorMsg)
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		logger.Errorf("Error when unmarshaling stored token object for %v: %v", req.Email, err)
		return merrors.InternalServerError("onboarding.CompleteSignup", internalErrorMsg)
	}
	if tok.Token != req.Token { // not checking expiry here because we've already checked it during Verify() step
		return merrors.Forbidden("onboarding.CompleteSignup.invalid_token", "The token you provided is incorrect")
	}

	rsp.Namespace = namespace

	// take secret from the request
	secret := req.Secret

	// generate a random secret
	if len(req.Secret) == 0 {
		secret = uuid.New().String()
	}
	_, err = e.auth.Generate(tok.CustomerID, auth.WithSecret(secret), auth.WithIssuer(namespace), auth.WithName(req.Email), auth.WithMetadata(map[string]string{
		"username": req.Username,
	}))
	if err != nil {
		logger.Errorf("Error generating token for %v: %v", tok.CustomerID, err)
		return merrors.InternalServerError("onboarding.CompleteSignup", internalErrorMsg)
	}

	t, err := e.auth.Token(auth.WithCredentials(tok.CustomerID, secret), auth.WithTokenIssuer(namespace))
	if err != nil {
		logger.Errorf("Can't get token for %v: %v", tok.CustomerID, err)
		return merrors.InternalServerError("onboarding.CompleteSignup", internalErrorMsg)
	}
	rsp.AuthToken = &onboarding.AuthToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry.Unix(),
		Created:      t.Created.Unix(),
	}

	return nil
}

func (e *Signup) Recover(ctx context.Context, req *onboarding.RecoverRequest, rsp *onboarding.RecoverResponse) error {
	logger.Info("Received Signup.Recover request")
	_, found := e.cache.Get(req.Email)
	if found {
		return merrors.BadRequest("onboarding.recover", "We have issued a recovery email recently. Please check that.")
	}

	token := uuid.New().String()
	created := time.Now().Unix()
	err := e.resetCode.Create(ResetToken{
		ID:      req.Email,
		Token:   token,
		Created: created,
	})
	logger.Infof("Sent recovery code %v to email %v", token, req.Email)
	if err != nil {
		return merrors.InternalServerError("onboarding.recover", err.Error())
	}

	err = e.sendEmail(ctx, req.Email, e.config.Sendgrid.RecoveryTemplateID, map[string]interface{}{
		"token": token,
	})
	if err == nil {
		e.cache.Set(req.Email, true, cache.DefaultExpiration)
	}

	return err
}

func (e *Signup) ResetPassword(ctx context.Context, req *onboarding.ResetPasswordRequest, rsp *onboarding.ResetPasswordResponse) error {
	m := ResetToken{}
	if req.Email == "" {
		return errors.New("Email is required")
	}
	err := e.resetCode.Read(model.QueryEquals("ID", req.Email), &m)
	if err != nil {
		return err
	}

	if m.ID == "" {
		return errors.New("can't find token")
	}
	if m.Token == "" {
		return errors.New("can't find token")
	}
	if m.Created == 0 {
		return errors.New("expiry can't be calculated")
	}
	if m.Token != req.Token {
		return errors.New("tokens don't match")
	}
	if time.Unix(m.Created, 0).Before(time.Now().Add(-1 * 10 * time.Minute)) {
		return errors.New("token expired")
	}

	_, err = e.accounts.ChangeSecret(cont.DefaultContext, &authproto.ChangeSecretRequest{
		Id:        req.Email,
		NewSecret: req.Password,
		Options: &authproto.Options{
			Namespace: req.Namespace,
		},
	}, client.WithAuthToken())
	if err != nil {
		return err
	}
	e.resetCode.Delete(model.QueryByID(m.ID))
	return err
}
