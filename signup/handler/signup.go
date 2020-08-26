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
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/client"
	merrors "github.com/micro/go-micro/v3/errors"
	logger "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/store"
	mconfig "github.com/micro/micro/v3/service/config"
	mstore "github.com/micro/micro/v3/service/store"
	"github.com/sethvargo/go-diceware/diceware"

	signup "github.com/m3o/services/signup/proto/signup"

	inviteproto "github.com/m3o/services/invite/proto"
	paymentsproto "github.com/m3o/services/payments/provider/proto"
	plproto "github.com/m3o/services/platform/proto"
)

const (
	storePrefixAccountSecrets = "secrets/"
	storePrefixNamesapce      = "namespaces/"
	storePrefixOwner          = "owner/"
	expiryDuration            = 5 * time.Minute
)

type tokenToEmail struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type Signup struct {
	paymentService     paymentsproto.ProviderService
	inviteService      inviteproto.InviteService
	platformService    plproto.PlatformService
	auth               auth.Auth
	sendgridTemplateID string
	sendgridAPIKey     string
	// M3O Platform Subscription plan id
	planID string
	// M3O Addition Users plan id
	additionUsersPriceID string
	emailFrom            string
	paymentMessage       string
	testMode             bool
}

var (
	// TODO: move this message to a better location
	// Message is a predefined message returned during signup
	Message = "Please complete signup at https://m3o.com/subscribe?email=%s and enter the generated token ID: "
)

func NewSignup(paymentService paymentsproto.ProviderService,
	inviteService inviteproto.InviteService,
	platformService plproto.PlatformService, auth auth.Auth) *Signup {

	apiKey := mconfig.Get("micro", "signup", "sendgrid", "api_key").String("")
	templateID := mconfig.Get("micro", "signup", "sendgrid", "template_id").String("")
	planID := mconfig.Get("micro", "signup", "plan_id").String("")
	additionUsersPriceID := mconfig.Get("micro", "signup", "additional_users_price_id").String("")
	emailFrom := mconfig.Get("micro", "signup", "email_from").String("Micro Team <support@micro.mu>")
	testMode := mconfig.Get("micro", "signup", "test_env").Bool(false)
	paymentMessage := mconfig.Get("micro", "signup", "message").String(Message)

	if len(apiKey) == 0 {
		logger.Error("No sendgrid API key provided")
	}
	if len(templateID) == 0 {
		logger.Error("No sendgrid template ID provided")
	}
	if len(planID) == 0 {
		logger.Error("No stripe plan id")
	}
	if len(additionUsersPriceID) == 0 {
		logger.Error("No addition user plan id")
	}
	return &Signup{
		paymentService:       paymentService,
		inviteService:        inviteService,
		platformService:      platformService,
		auth:                 auth,
		sendgridAPIKey:       apiKey,
		sendgridTemplateID:   templateID,
		planID:               planID,
		additionUsersPriceID: additionUsersPriceID,
		emailFrom:            emailFrom,
		testMode:             testMode,
		paymentMessage:       paymentMessage,
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
// A stripe customers and a verification token will be created and an email sent.
func (e *Signup) SendVerificationEmail(ctx context.Context,
	req *signup.SendVerificationEmailRequest,
	rsp *signup.SendVerificationEmailResponse) error {
	logger.Info("Received Signup.SendVerificationEmail request")

	_, isAllowed := e.isAllowedToSignup(ctx, req.Email)
	if !isAllowed {
		return merrors.Forbidden("signup.notallowed", "user has not been invited to sign up")
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

	if err := mstore.Write(&store.Record{
		Key:    req.Email,
		Value:  bytes,
		Expiry: expiryDuration,
	}); err != nil {
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

	if tok.Token != req.Token {
		return errors.New("Invalid token")
	}

	// set the response message
	rsp.Message = fmt.Sprintf(e.paymentMessage, req.Email)
	// we require payment for any signup
	// if not set the CLI will try complete signup without payment id
	rsp.PaymentRequired = true

	// Otherwisewe just return without an error but with no token
	_, err = e.paymentService.CreateCustomer(ctx, &paymentsproto.CreateCustomerRequest{
		Customer: &paymentsproto.Customer{
			Id:   req.Email,
			Type: "user",
			Metadata: map[string]string{
				"email": req.Email,
			},
		},
	}, client.WithAuthToken())
	if err != nil {
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

func (e *Signup) getNamespace(email string) (string, error) {
	key := storePrefixNamesapce + email
	recs, err := mstore.Read(key)
	if err != nil {
		return "", err
	}
	return string(recs[0].Value), nil
}

func (e *Signup) saveNamespace(email, namespace string) error {
	key := storePrefixNamesapce + email
	return mstore.Write(&store.Record{Key: key, Value: []byte(namespace)})
}

func (e *Signup) getOwner(namespace string) (string, error) {
	key := storePrefixOwner + namespace
	recs, err := mstore.Read(key)
	if err != nil {
		return "", err
	}
	return string(recs[0].Value), nil
}

func (e *Signup) saveOwner(email, namespace string) error {
	key := storePrefixOwner + namespace
	return mstore.Write(&store.Record{Key: key, Value: []byte(email)})
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
	if tok.Token != req.Token {
		return errors.New("invalid token")
	}

	if isJoining {
		ownerEmail, err := e.getOwner(ns)
		if err != nil {
			return err
		}
		subs, err := e.paymentService.ListSubscriptions(ctx, &paymentsproto.ListSubscriptionsRequest{
			CustomerId:   ownerEmail,
			CustomerType: "user",
		}, client.WithAuthToken())
		if err != nil {
			return merrors.InternalServerError("signup", "Error listing subscriptions: %v", err)
		}
		var subscription *paymentsproto.Subscription
		for _, sub := range subs.Subscriptions {
			if sub.Plan.Id == e.additionUsersPriceID {
				subscription = sub
			}
		}
		if subscription == nil {
			logger.Info("Creating subscription with quantity 1")
			_, err = e.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
				CustomerId:   ownerEmail,
				CustomerType: "user",
				PriceId:      e.additionUsersPriceID,
				Quantity:     1,
			}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
		} else {
			logger.Info("Increasing subscription quantity")
			_, err = e.paymentService.UpdateSubscription(ctx, &paymentsproto.UpdateSubscriptionRequest{
				SubscriptionId: subscription.Id,
				CustomerId:     ownerEmail,
				CustomerType:   "user",
				PriceId:        e.additionUsersPriceID,
				Quantity:       subscription.Quantity + 1,
			}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
		}
		if err != nil {
			return merrors.InternalServerError("signup", "Error increasing additional user quantity: %v", err)
		}
	} else {
		_, err = e.paymentService.CreatePaymentMethod(ctx, &paymentsproto.CreatePaymentMethodRequest{
			CustomerId:   req.Email,
			CustomerType: "user",
			Id:           req.PaymentMethodID,
		}, client.WithAuthToken())
		if err != nil {
			return err
		}

		_, err = e.paymentService.SetDefaultPaymentMethod(ctx, &paymentsproto.SetDefaultPaymentMethodRequest{
			CustomerId:      req.Email,
			CustomerType:    "user",
			PaymentMethodId: req.PaymentMethodID,
		}, client.WithAuthToken())
		if err != nil {
			return err
		}

		_, err = e.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
			Quantity:     1,
			CustomerId:   req.Email,
			CustomerType: "user",
			PlanId:       e.planID,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
		if err != nil {
			return err
		}

		var err error
		ns, err = e.createNamespace(ctx)
		if err != nil {
			return err
		}
		err = e.saveOwner(req.Email, ns)
		if err != nil {
			return err
		}
	}
	err = e.saveNamespace(req.Email, ns)
	if err != nil {
		return err
	}

	rsp.Namespace = ns

	// take secret from the request
	secret := req.Secret

	// generate a random secret
	if len(req.Secret) == 0 {
		secret = uuid.New().String()
	}
	_, err = e.auth.Generate(req.Email, auth.WithSecret(secret), auth.WithIssuer(ns))
	if err != nil {
		return err
	}

	t, err := e.auth.Token(auth.WithCredentials(req.Email, secret), auth.WithTokenIssuer(ns))
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

func (e *Signup) createNamespace(ctx context.Context) (string, error) {
	list, err := diceware.Generate(3)
	if err != nil {
		return "", err
	}
	ns := strings.Join(list, "-")
	if !e.testMode {
		_, err = e.platformService.CreateNamespace(ctx, &plproto.CreateNamespaceRequest{
			Name: ns,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	}
	return ns, err
}
