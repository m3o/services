package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	cproto "github.com/m3o/services/customers/proto"
	inviteproto "github.com/m3o/services/invite/proto"
	nproto "github.com/m3o/services/namespaces/proto"
	signup "github.com/m3o/services/signup/proto/signup"
	sproto "github.com/m3o/services/subscriptions/proto"
	"github.com/patrickmn/go-cache"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/client"
	merrors "github.com/micro/go-micro/v3/errors"
	logger "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/store"
	mconfig "github.com/micro/micro/v3/service/config"
	mstore "github.com/micro/micro/v3/service/store"
)

const (
	expiryDuration = 5 * time.Minute
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
	auth                auth.Auth
	sendgridTemplateID  string
	recoverTemplateID   string
	sendgridAPIKey      string
	emailFrom           string
	paymentMessage      string
	testMode            bool
	cache               *cache.Cache
}

var (
	// TODO: move this message to a better location
	// Message is a predefined message returned during signup
	Message = "Please complete signup at https://m3o.com/subscribe?email=%s and enter the generated token ID: "
)

func NewSignup(inviteService inviteproto.InviteService,
	customerService cproto.CustomersService,
	namespaceService nproto.NamespacesService,
	subscriptionService sproto.SubscriptionsService,
	auth auth.Auth) *Signup {

	apiKey := mconfig.Get("micro", "signup", "sendgrid", "api_key").String("")
	templateID := mconfig.Get("micro", "signup", "sendgrid", "template_id").String("")
	recoverTemplateID := mconfig.Get("micro", "signup", "sendgrid", "recovery_template_id").String("")
	emailFrom := mconfig.Get("micro", "signup", "email_from").String("Micro Team <support@micro.mu>")
	testMode := mconfig.Get("micro", "signup", "test_env").Bool(false)
	paymentMessage := mconfig.Get("micro", "signup", "message").String(Message)

	if len(apiKey) == 0 {
		logger.Error("No sendgrid API key provided")
	}
	if len(templateID) == 0 {
		logger.Error("No sendgrid template ID provided")
	}
	return &Signup{
		inviteService:       inviteService,
		customerService:     customerService,
		namespaceService:    namespaceService,
		subscriptionService: subscriptionService,
		auth:                auth,
		sendgridAPIKey:      apiKey,
		sendgridTemplateID:  templateID,
		emailFrom:           emailFrom,
		testMode:            testMode,
		paymentMessage:      paymentMessage,
		recoverTemplateID:   recoverTemplateID,
		cache:               cache.New(1*time.Minute, 5*time.Minute),
	}
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

// SendVerificationEmail is the first step in the signup flow.SendVerificationEmail
// A stripe customer and a verification token will be created and an email sent.
func (e *Signup) SendVerificationEmail(ctx context.Context,
	req *signup.SendVerificationEmailRequest,
	rsp *signup.SendVerificationEmailResponse) error {
	logger.Info("Received Signup.SendVerificationEmail request")

	_, isAllowed := e.isAllowedToSignup(ctx, req.Email)
	if !isAllowed {
		return merrors.Forbidden("signup.notallowed", "user has not been invited to sign up")
	}

	custResp, err := e.customerService.Create(ctx, &cproto.CreateRequest{
		Email: req.Email,
	}, client.WithAuthToken())
	if err != nil {
		return err
	}

	k := randStringBytesMaskImprSrc(8)
	tok := &tokenToEmail{
		Token:      k,
		Email:      req.Email,
		Created:    time.Now().Unix(),
		CustomerID: custResp.Customer.Id,
	}

	bytes, err := json.Marshal(tok)
	if err != nil {
		return err
	}

	if err := mstore.Write(&store.Record{
		Key:   req.Email,
		Value: bytes,
	}); err != nil {
		return err
	}

	if e.testMode {
		logger.Infof("Sending verification token '%v'", k)
	}

	// Send email
	// @todo send different emails based on if the account already exists
	// ie. registration vs login email.

	err = e.sendEmail(req.Email, e.sendgridTemplateID, map[string]interface{}{
		"token": k,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *Signup) isAllowedToSignup(ctx context.Context, email string) ([]string, bool) {
	// for now we're checking the invite service before allowing signup
	// TODO check for a valid invite code rather than just the email
	rsp, err := e.inviteService.Validate(ctx, &inviteproto.ValidateRequest{Email: email}, client.WithAuthToken())
	if err != nil {
		return nil, false
	}
	return rsp.Namespaces, true
}

// Lifted  from the invite service https://github.com/m3o/services/blob/master/projects/invite/handler/invite.go#L187
// sendEmailInvite sends an email invite via the sendgrid API using the
// predesigned email template. Docs: https://bit.ly/2VYPQD1
func (e *Signup) sendEmail(email, templateID string, templateData map[string]interface{}) error {
	if e.testMode {
		logger.Infof("Test mode enabled, not sending email to address '%v' ", email)
		return nil
	}
	reqBody, _ := json.Marshal(map[string]interface{}{
		"template_id": templateID,
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
				"dynamic_template_data": templateData,
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

func (e *Signup) Verify(ctx context.Context, req *signup.VerifyRequest, rsp *signup.VerifyResponse) error {
	logger.Info("Received Signup.Verify request")

	recs, err := mstore.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.New("can't verify: record not found")
	} else if err != nil {
		return fmt.Errorf("email verification error: %v", err)
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}

	if tok.Token != req.Token || time.Since(time.Unix(tok.Created, 0)) > expiryDuration {
		return errors.New("Invalid token")
	}

	// set the response message
	rsp.Message = fmt.Sprintf(e.paymentMessage, req.Email)
	// we require payment for any signup
	// if not set the CLI will try complete signup without payment id
	rsp.PaymentRequired = true

	if _, err := e.customerService.MarkVerified(ctx, &cproto.MarkVerifiedRequest{
		Email: req.Email,
	}, client.WithAuthToken()); err != nil {
		return err
	}

	// At this point the user should be allowed, only making this call to return namespaces
	namespaces, isAllowed := e.isAllowedToSignup(ctx, req.Email)
	if !isAllowed {
		return merrors.Forbidden("signup.notallowed", "user has not been invited to sign up")
	}
	rsp.Namespaces = namespaces
	return nil
}

func (e *Signup) CompleteSignup(ctx context.Context, req *signup.CompleteSignupRequest, rsp *signup.CompleteSignupResponse) error {
	logger.Info("Received Signup.CompleteSignup request")

	namespaces, isAllowed := e.isAllowedToSignup(ctx, req.Email)
	if !isAllowed {
		return merrors.Forbidden("signup.notallowed", "user has not been invited to sign up")
	}
	ns := ""
	isJoining := len(namespaces) > 0 && len(req.Namespace) > 0 && namespaces[0] == req.Namespace
	if isJoining {
		ns = namespaces[0]
	}

	recs, err := mstore.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.New("can't verify: record not found")
	} else if err != nil {
		return err
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}
	if tok.Token != req.Token { // not checking expiry here because we've already checked it during Verify() step
		return errors.New("invalid token")
	}

	if isJoining {
		if err := e.joinNamespace(ctx, tok.CustomerID, ns); err != nil {
			return err
		}
	} else {
		newNs, err := e.signupWithNewNamespace(ctx, tok.CustomerID, tok.Email, req.PaymentMethodID)
		if err != nil {
			return err
		}
		ns = newNs
	}

	rsp.Namespace = ns

	// take secret from the request
	secret := req.Secret

	// generate a random secret
	if len(req.Secret) == 0 {
		secret = uuid.New().String()
	}
	_, err = e.auth.Generate(tok.CustomerID, auth.WithSecret(secret), auth.WithIssuer(ns), auth.WithMetadata(map[string]string{"username": req.Email}))
	if err != nil {
		return err
	}

	t, err := e.auth.Token(auth.WithCredentials(tok.CustomerID, secret), auth.WithTokenIssuer(ns))
	if err != nil {
		return err
	}
	rsp.AuthToken = &signup.AuthToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry.Unix(),
		Created:      t.Created.Unix(),
	}
	return nil
}

func (e *Signup) Recover(ctx context.Context, req *signup.RecoverRequest, rsp *signup.RecoverResponse) error {
	logger.Info("Received Signup.Recover request")
	_, found := e.cache.Get(req.Email)
	if found {
		return merrors.BadRequest("signup.recover", "We have issued a recovery email recently. Please check that.")
	}

	custResp, err := e.customerService.Read(ctx, &cproto.ReadRequest{Email: req.Email})
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == 404 { // not found
			return merrors.NotFound("signup.recover", "Could not find account with that email address")
		}
		return merrors.InternalServerError("signup.recover", "Error looking up account")
	}

	listRsp, err := e.namespaceService.List(ctx, &nproto.ListRequest{
		User: custResp.Customer.Id,
	}, client.WithAuthToken())
	if err != nil {
		return merrors.InternalServerError("signup.recover", "Error calling namespace service: %v", err)
	}
	if len(listRsp.Namespaces) == 0 {
		return merrors.BadRequest("signup.recover", "We don't recognize this account")
	}

	// Sendgrid wants objects in a list not string
	namespaces := []map[string]string{}
	for _, v := range listRsp.Namespaces {
		namespaces = append(namespaces, map[string]string{
			"id": v.Id,
		})
	}

	logger.Infof("Sending email with data %v", namespaces)
	err = e.sendEmail(req.Email, e.recoverTemplateID, map[string]interface{}{
		"namespaces": namespaces,
	})
	if err == nil {
		e.cache.Set(req.Email, true, cache.DefaultExpiration)
	}
	return err
}

func (e *Signup) signupWithNewNamespace(ctx context.Context, customerID, email, paymentMethodID string) (string, error) {
	// TODO fix type to be more than just developer
	_, err := e.subscriptionService.Create(ctx, &sproto.CreateRequest{CustomerID: customerID, Type: "developer", PaymentMethodID: paymentMethodID, Email: email}, client.WithAuthToken())
	if err != nil {
		return "", err
	}
	nsRsp, err := e.namespaceService.Create(ctx, &nproto.CreateRequest{Owners: []string{customerID}}, client.WithAuthToken())
	if err != nil {
		return "", err
	}
	return nsRsp.Namespace.Id, nil
}

func (e *Signup) joinNamespace(ctx context.Context, customerID, ns string) error {
	rsp, err := e.namespaceService.Read(ctx, &nproto.ReadRequest{
		Id: ns,
	}, client.WithAuthToken())
	if err != nil {
		return err
	}
	ownerID := rsp.Namespace.Owners[0]

	_, err = e.subscriptionService.AddUser(ctx, &sproto.AddUserRequest{OwnerID: ownerID, NewUserID: customerID}, client.WithAuthToken())
	if err != nil {
		return merrors.InternalServerError("signup.join.subscription", "Error adding user to subscription %s", err)
	}

	_, err = e.namespaceService.AddUser(ctx, &nproto.AddUserRequest{Namespace: ns, User: customerID}, client.WithAuthToken())
	if err != nil {
		return merrors.InternalServerError("signup.join.namespace", "Error adding user to namespace %s", err)
	}

	return nil
}
