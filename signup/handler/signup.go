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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/config"
	merrors "github.com/micro/go-micro/v2/errors"
	logger "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/sethvargo/go-diceware/diceware"

	signup "github.com/micro/services/signup/proto/signup"

	inviteproto "github.com/micro/services/account/invite/proto"
	k8sproto "github.com/micro/services/kubernetes/service/proto"
	paymentsproto "github.com/micro/services/payments/provider/proto"
)

const (
	storePrefixAccountSecrets = "secrets/"
	storePrefixNamesapce      = "namespaces/"
	expiryDuration            = 5 * time.Minute
)

type tokenToEmail struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type Signup struct {
	paymentService     paymentsproto.ProviderService
	inviteService      inviteproto.InviteService
	k8sService         k8sproto.KubernetesService
	store              store.Store
	auth               auth.Auth
	sendgridTemplateID string
	sendgridAPIKey     string
	planID             string
	emailFrom          string
	testMode           bool
}

func NewSignup(paymentService paymentsproto.ProviderService,
	inviteService inviteproto.InviteService,
	k8sService k8sproto.KubernetesService,
	store store.Store,
	config config.Config,
	auth auth.Auth) *Signup {

	apiKey := config.Get("micro", "signup", "sendgrid", "api_key").String("")
	templateID := config.Get("micro", "signup", "sendgrid", "template_id").String("")
	planID := config.Get("micro", "signup", "plan_id").String("")
	emailFrom := config.Get("micro", "signup", "email_from").String("Micro Team <support@micro.mu>")
	testMode := config.Get("micro", "signup", "test_env").Bool(false)

	if len(apiKey) == 0 {
		logger.Error("No sendgrid API key provided")
	}
	if len(templateID) == 0 {
		logger.Error("No sendgrid template ID provided")
	}
	if len(planID) == 0 {
		logger.Error("No stripe plan id")
	}
	return &Signup{
		paymentService:     paymentService,
		inviteService:      inviteService,
		k8sService:         k8sService,
		store:              store,
		auth:               auth,
		sendgridAPIKey:     apiKey,
		sendgridTemplateID: templateID,
		planID:             planID,
		emailFrom:          emailFrom,
		testMode:           testMode,
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

	if !e.isAllowedToSignup(ctx, req.Email) {
		return merrors.Forbidden("go.micro.service.signup.notallowed", "user has not been invited to sign up")
	}

	k := randStringBytesMaskImprSrc(8)
	tok := &tokenToEmail{
		Token: k,
		Email: req.Email,
	}

	bytes, err := json.Marshal(tok)
	if err != nil {
		return err
	}

	if err := e.store.Write(&store.Record{
		Key:   req.Email,
		Value: bytes}, store.WriteExpiry(time.Now().Add(expiryDuration))); err != nil {
		return err
	}

	if e.testMode {
		logger.Infof("Sending verification token '%v'", k)
	}

	// Send email
	// @todo send different emails based on if the account already exists
	// ie. registration vs login email.
	err = e.sendEmail(req.Email, k)
	if err != nil {
		return err
	}

	return nil
}

func (e *Signup) isAllowedToSignup(ctx context.Context, email string) bool {
	// for now we're checking the invite service before allowing signup
	// TODO check for a valid invite code rather than just the email
	_, err := e.inviteService.Validate(ctx, &inviteproto.ValidateRequest{Email: email})
	return err == nil
}

// Lifted  from the invite service https://github.com/micro/services/blob/master/projects/invite/handler/invite.go#L187
// sendEmailInvite sends an email invite via the sendgrid API using the
// predesigned email template. Docs: https://bit.ly/2VYPQD1
func (e *Signup) sendEmail(email, token string) error {
	logger.Infof("Sending email to address '%v'", email)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"template_id": e.sendgridTemplateID,
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
	})

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+e.sendgridAPIKey)
	req.Header.Set("Content-Type", "application/json")
	rsp, err := new(http.Client).Do(req)
	if err != nil {
		logger.Infof("Could not send email to %v, error: %v", email, err)
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			logger.Errorf("Could not send email to %v, error: %v", email, err.Error())
			return err
		}
		logger.Errorf("Could not send email to %v, error: %v", email, string(bytes))
		return merrors.InternalServerError("go.micro.service.signup.sendemail", "error sending email")
	}
	return nil
}

func (e *Signup) Verify(ctx context.Context,
	req *signup.VerifyRequest,
	rsp *signup.VerifyResponse) error {
	logger.Info("Received Signup.Verify request")

	recs, err := e.store.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.New("can't verify: record not found")
	} else if err != nil {
		return fmt.Errorf("email verification error: %v", err)
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}

	if tok.Token != req.Token {
		return errors.New("Invalid token")
	}

	secret, err := e.getAccountSecret(req.Email)
	if err != store.ErrNotFound && err != nil {
		return fmt.Errorf("can't get account secret: %v", err)
	}
	// If the user has a secret it means the account is ready
	// to be used, so we log them in.
	if len(secret) > 0 {
		token, err := e.auth.Token(auth.WithCredentials(req.Email, secret))
		if err != nil {
			return err
		}
		rsp.AuthToken = &signup.AuthToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry.Unix(),
			Created:      token.Created.Unix(),
		}
		ns, err := e.getNamespace(req.Email)
		if err != store.ErrNotFound {
			return err
		}
		// @todo what to do if namespace is not found?
		rsp.Namespace = ns
		return nil
	}

	// Otherwisewe just return without an error but with no token
	_, err = e.paymentService.CreateCustomer(ctx, &paymentsproto.CreateCustomerRequest{
		Customer: &paymentsproto.Customer{
			Id:   req.Email,
			Type: "user",
		},
	})
	return err
}

func (e *Signup) getNamespace(email string) (string, error) {
	key := storePrefixNamesapce + email
	recs, err := e.store.Read(key)
	if err != nil {
		return "", err
	}
	return string(recs[0].Value), nil
}

func (e *Signup) saveNamespace(email, namespace string) error {
	key := storePrefixNamesapce + email
	return e.store.Write(&store.Record{Key: key, Value: []byte(namespace)})
}

func (e *Signup) CompleteSignup(ctx context.Context,
	req *signup.CompleteSignupRequest,
	rsp *signup.CompleteSignupResponse) error {
	logger.Info("Received Signup.CompleteSignup request")

	recs, err := e.store.Read(req.Email)
	if err == store.ErrNotFound {
		return errors.New("can't verify: record not found")
	} else if err != nil {
		return err
	}

	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}
	if tok.Token != req.Token {
		return errors.New("invalid token")
	}

	_, err = e.paymentService.CreatePaymentMethod(ctx, &paymentsproto.CreatePaymentMethodRequest{
		CustomerId:   req.Email,
		CustomerType: "user",
		Id:           req.PaymentMethodID,
	})
	if err != nil {
		return err
	}

	_, err = e.paymentService.SetDefaultPaymentMethod(ctx, &paymentsproto.SetDefaultPaymentMethodRequest{
		CustomerId:      req.Email,
		CustomerType:    "user",
		PaymentMethodId: req.PaymentMethodID,
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
	ns, err := e.createNamespace(ctx)
	if err != nil {
		return err
	}
	err = e.saveNamespace(req.Email, ns)
	if err != nil {
		return err
	}
	rsp.Namespace = ns

	_, err = e.auth.Generate(req.Email, auth.WithSecret(secret), auth.WithIssuer(ns))
	if err != nil {
		return err
	}
	t, err := e.auth.Token(auth.WithCredentials(req.Email, secret))
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

// lifted from https://github.com/micro/services/blob/550220a6eff2604b3e6d58d09db2b4489967019c/account/web/handler/handler.go#L114
func (e *Signup) setAccountSecret(id, secret string) error {
	key := storePrefixAccountSecrets + id
	return e.store.Write(&store.Record{Key: key, Value: []byte(secret)})
}

func (e *Signup) getAccountSecret(id string) (string, error) {
	key := storePrefixAccountSecrets + id
	recs, err := e.store.Read(key)
	if err != nil {
		return "", err
	}
	return string(recs[0].Value), nil
}

func (e *Signup) createNamespace(ctx context.Context) (string, error) {
	list, err := diceware.Generate(3)
	if err != nil {
		return "", err
	}
	ns := strings.Join(list, "-")
	if !e.testMode {
		_, err = e.k8sService.CreateNamespace(ctx, &k8sproto.CreateNamespaceRequest{
			Name: ns,
		})
	}
	return ns, err
}
