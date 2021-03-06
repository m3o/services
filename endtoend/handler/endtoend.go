package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	balancepb "github.com/m3o/services/balance/proto"
	"github.com/micro/micro/v3/service"

	alertpb "github.com/m3o/services/alert/proto/alert"
	custpb "github.com/m3o/services/customers/proto"
	endtoend "github.com/m3o/services/endtoend/proto"
	onpb "github.com/m3o/services/onboarding/proto"
	pubpb "github.com/m3o/services/publicapi/proto"
	"github.com/micro/micro/v3/service/client"
	mconfig "github.com/micro/micro/v3/service/config"
	merrors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	mstore "github.com/micro/micro/v3/service/store"

	"github.com/google/uuid"
)

const (
	signupFrom     = "Micro Team <support@m3o.com>"
	signupSubject  = "M3O Platform - Email Verification"
	keyOtp         = "otp"
	keyCheckResult = "checkResult"

	// how fresh does a check need to be? cron runs every 5 mins and check takes just over 1min.
	checkBuffer = 7 * time.Minute
)

var (
	otpRegexp = regexp.MustCompile("please copy and paste this verification code into the signup form:\\s*([a-zA-Z]*)\\s*This code expires")
)

func NewEndToEnd(srv *service.Service) *Endtoend {
	val, err := mconfig.Get("micro.endtoend.email")
	if err != nil {
		log.Fatalf("Cannot configure, error finding email: %s", err)
	}
	email := val.String("")
	if len(email) == 0 {
		log.Fatalf("Cannot configure, email not configured")
	}
	return &Endtoend{
		email:    email,
		custSvc:  custpb.NewCustomersService("customers", srv.Client()),
		alertSvc: alertpb.NewAlertService("alert", srv.Client()),
		balSvc:   balancepb.NewBalanceService("balance", srv.Client()),
		pubSvc:   pubpb.NewPublicapiService("publicapi", srv.Client()),
	}
}

func (e *Endtoend) Mailin(ctx context.Context, req *json.RawMessage, rsp *MailinResponse) error {
	log.Infof("Received Endtoend.Mailin request %d", len(*req))
	var inbound mailinMessage

	if err := json.Unmarshal(*req, &inbound); err != nil {
		log.Errorf("Error unmarshalling request %s", err)
		// returning err would make the email bounce
		return nil
	}
	// TODO make this configurable
	if !strings.Contains(inbound.Headers["to"].(string), e.email) ||
		!strings.Contains(inbound.Headers["from"].(string), signupFrom) ||
		!strings.Contains(inbound.Headers["subject"].(string), signupSubject) {
		// skip
		log.Debugf("Skipping inbound %+v", inbound)
		return nil
	}

	tok := otpRegexp.FindStringSubmatch(inbound.Plain)
	if len(tok) != 2 {
		log.Errorf("Couldn't find token in email body: %s", inbound.Plain)
		// returning err would make the email bounce
		return nil
	}
	otp := otp{
		Token: tok[1],
		Time:  time.Now().Unix(),
	}
	b, err := json.Marshal(otp)
	if err != nil {
		log.Errorf("Failed to marshal otp %s", err)
		// returning err would make the email bounce
		return nil
	}
	if err := mstore.Write(&mstore.Record{
		Key:   keyOtp,
		Value: b,
	}); err != nil {
		log.Errorf("Error storing OTP %s", err)
		return nil
	}
	return nil
}

func (e *Endtoend) Check(ctx context.Context, request *endtoend.Request, response *endtoend.Response) error {
	log.Info("Received Endtoend.Check request")
	recs, err := mstore.Read(keyCheckResult)
	if err != nil {
		return merrors.InternalServerError("endtoend.check.store", "Failed to load last result %s", err)
	}
	if len(recs) == 0 {
		return merrors.InternalServerError("endtoend.check.noresults", "Failed to load last result, no results found")
	}
	cr := checkResult{}
	if err := json.Unmarshal(recs[0].Value, &cr); err != nil {
		return merrors.InternalServerError("endtoend.check.unmarshal", "Failed to unmarshal last result %s", err)
	}
	if cr.Passed && time.Now().Add(-checkBuffer).Unix() < cr.Time {
		response.StatusCode = 200
		return nil
	}
	response.StatusCode = 500
	response.Body = cr.Error
	if len(response.Body) == 0 {
		response.Body = "No recent successful check"
	}
	return merrors.New("endtoend.chack.failed", response.Body, response.StatusCode)

}

func (e *Endtoend) RunCheck(ctx context.Context, request *endtoend.Request, response *endtoend.Response) error {
	go e.runCheck()
	return nil
}

func (e *Endtoend) CronCheck() {
	e.runCheck()
}

func (e *Endtoend) runCheck() error {
	var err error
	start := time.Now()
	defer func() {
		// record the result
		result := checkResult{
			Time:   time.Now().Unix(),
			Passed: err == nil,
		}
		if err != nil {
			result.Error = err.Error()
		}
		b, _ := json.Marshal(result)

		mstore.Write(&mstore.Record{
			Key:   keyCheckResult,
			Value: b,
		})
		if err == nil {
			log.Infof("Signup check took %s to complete", time.Now().Sub(start))
			return
		}

		// alert if required
		e.alertSvc.ReportEvent(context.Background(), &alertpb.ReportEventRequest{
			Event: &alertpb.Event{
				Category: "monitoring",
				Action:   "signup",
				Label:    "endtoend",
				Value:    1,
				Metadata: map[string]string{"error": err.Error()},
			},
		}, client.WithAuthToken())
	}()
	if err = e.signup(); err != nil {
		log.Errorf("Error during signup %s", err)
		return err
	}
	return nil
}

func (e *Endtoend) signup() error {
	// reset, delete any existing customers. Try this a few times, we sometimes get timeout
	var delErr error
	for i := 0; i < 3; i++ {
		_, err := e.custSvc.Delete(context.TODO(), &custpb.DeleteRequest{Email: e.email, Force: true}, client.WithAuthToken(), client.WithRequestTimeout(15*time.Second))
		if err == nil {
			delErr = nil
			break
		}
		merr, ok := err.(*merrors.Error)
		if ok && (merr.Code == 404 || strings.Contains(merr.Detail, "not found")) {
			delErr = nil
			break
		}
		delErr = fmt.Errorf("error while cleaning up existing customer %s", err)
	}
	if delErr != nil {
		return delErr
	}

	start := time.Now()
	rsp, err := http.Post("https://api.m3o.com/onboarding/signup/SendVerificationEmail",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"email":"%s"}`, e.email)))
	if err != nil {
		return fmt.Errorf("error requesting verification email %s", err)
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("error requesting verification email %s %s", rsp.Status, string(b))
	}

	code := ""

	loopStart := time.Now()
	for time.Now().Sub(loopStart) < 2*time.Minute {
		time.Sleep(10 * time.Second)
		log.Infof("Checking for otp")
		recs, err := mstore.Read(keyOtp)
		if err != nil {
			log.Errorf("Error reading otp from store %s", err)
			continue
		}
		if len(recs) == 0 {
			log.Infof("No recs found")
			continue
		}
		otp := otp{}
		if err := json.Unmarshal(recs[0].Value, &otp); err != nil {
			log.Errorf("Error unmarshalling otp from store %s", err)
			continue
		}
		if otp.Time < start.Unix() {
			log.Infof("Otp is old")
			// old token
			continue
		}
		log.Infof("Found otp")
		code = otp.Token
		break
	}
	if len(code) == 0 {
		return fmt.Errorf("no OTP code received by email")
	}
	rsp, err = http.Post("https://api.m3o.com/onboarding/signup/CompleteSignup",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"email":"%s", "token":"%s", "secret":"%s"}`, e.email, code, uuid.New().String())))

	if err != nil {
		return fmt.Errorf("error completing signup %s", err)
	}
	defer rsp.Body.Close()
	b, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		return fmt.Errorf("error completing signup %s %s", rsp.Status, string(b))
	}
	var srsp onpb.CompleteSignupResponse
	json.Unmarshal(b, &srsp)

	var custErr error
	loopStart = time.Now()
	for time.Now().Sub(loopStart) < 2*time.Minute {
		time.Sleep(10 * time.Second)
		rsp, err := e.custSvc.Read(context.TODO(), &custpb.ReadRequest{Email: e.email}, client.WithAuthToken())
		if err != nil {
			custErr = fmt.Errorf("error checking customer status %s", err)
			continue
		}
		if rsp.Customer.Status != "verified" {
			custErr = fmt.Errorf("customer status should be verified but is %s", rsp.Customer.Status)
			continue
		}
		custErr = nil
		break
	}

	if custErr != nil {
		return custErr
	}

	loopStart = time.Now()
	currBal := int64(0)
	for time.Now().Sub(loopStart) < 1*time.Minute {
		// check balance for intro credit
		balRsp, err := e.balSvc.Current(context.Background(), &balancepb.CurrentRequest{CustomerId: srsp.CustomerID}, client.WithAuthToken())
		if err == nil && balRsp.CurrentBalance > 0 {
			currBal = balRsp.CurrentBalance
			break
		}
		log.Errorf("balance response %s", err)
	}
	log.Infof("Balance is %d", currBal)
	if currBal == 0 {
		return fmt.Errorf("no intro credit was applied to customer")
	}

	// generate a key
	req, err := http.NewRequest("POST", "https://api.m3o.com/v1/api/keys/generate", strings.NewReader(`{"description":"test key", "scopes": ["*"]}`))
	if err != nil {
		return fmt.Errorf("error generating key request %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+srsp.AuthToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	rsp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error generating key %s", err)
	}
	defer rsp.Body.Close()
	b, _ = ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		return fmt.Errorf("error generating key %s %s", rsp.Status, string(b))
	}
	keyRsp := struct {
		ApiKey string `json:"api_key"`
	}{}
	if err := json.Unmarshal(b, &keyRsp); err != nil {
		return fmt.Errorf("error generating key, failed to unmarshal response %s", err)
	}

	// sleep to allow the api key to be enabled - async processing bla bla bla
	time.Sleep(10 * time.Second)

	// run some apis
	pubrsp, err := e.pubSvc.List(context.Background(), &pubpb.ListRequest{}, client.WithAuthToken())
	if err != nil {
		return fmt.Errorf("error listing apis %s", err)
	}

	exampleErrs := []string{}
	for _, api := range pubrsp.Apis {
		var examples apiExamples
		if len(api.ExamplesJson) == 0 {
			log.Infof("No examples, skipping %s", api.Name)
			continue
		}

		keys, err := objectKeys(api.ExamplesJson)
		if err != nil {
			log.Errorf("Failed to find keys for %s %s", api.Name, err)
			exampleErrs = append(exampleErrs, fmt.Sprintf("Failed to find keys for %s %s", api.Name, err))
			continue
		}
		if err := json.Unmarshal([]byte(api.ExamplesJson), &examples); err != nil {
			log.Errorf("Failed to unmarshal examples for %s %s", api.Name, err)
			exampleErrs = append(exampleErrs, fmt.Sprintf("Failed to unmarshal examples for %s %s", api.Name, err))
			continue
		}

		// process the examples in json order to respect implicit dependencies
		for _, endpointName := range keys {
			exs := examples[endpointName]
			for _, ex := range exs {
				b, err := json.Marshal(ex.Request)
				if err != nil {
					log.Errorf("Failed to marshal example request for %s %s %s %s", api.Name, endpointName, ex.Title, err)
					exampleErrs = append(exampleErrs, fmt.Sprintf("Failed to marshal example request for %s %s %s %s", api.Name, endpointName, ex.Title, err))
					continue
				}

				req, err := http.NewRequest("POST", fmt.Sprintf("https://api.m3o.com/v1/%s/%s", api.Name, endpointName), bytes.NewReader(b))
				if err != nil {
					log.Errorf("Error generating example request %s %s %s %s", api.Name, endpointName, ex.Title, err)
					exampleErrs = append(exampleErrs, fmt.Sprintf("Error generating example request %s %s %s %s", api.Name, endpointName, ex.Title, err))
					continue
				}
				req.Header.Set("Authorization", "Bearer "+keyRsp.ApiKey)
				req.Header.Set("Content-Type", "application/json")
				rsp, err = http.DefaultClient.Do(req)
				if err != nil {
					log.Errorf("Error running example %s %s %s %s", api.Name, endpointName, ex.Title, err)
					exampleErrs = append(exampleErrs, fmt.Sprintf("Error running example %s %s %s %s", api.Name, endpointName, ex.Title, err))
					continue
				}
				defer rsp.Body.Close()
				b, _ = ioutil.ReadAll(rsp.Body)
				if rsp.StatusCode != 200 {
					log.Errorf("Error running example %s %s %s %s %s", api.Name, endpointName, ex.Title, rsp.Status, string(b))
					exampleErrs = append(exampleErrs, fmt.Sprintf("Error running example %s %s %s %s %s", api.Name, endpointName, ex.Title, rsp.Status, string(b)))
					continue
				}
				log.Infof("API response for example %s %s %s %s %s", api.Name, endpointName, ex.Title, rsp.Status, string(b))
				// TODO test response against expected response with fuzzy matching

			}
		}
	}
	// TODO add credit via stripe
	if len(exampleErrs) > 0 {
		return fmt.Errorf("Errors running api examples. %s", strings.Join(exampleErrs, ". "))
	}
	return nil
}

func objectKeys(s string) ([]string, error) {
	d := json.NewDecoder(strings.NewReader(s))
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	if t != json.Delim('{') {
		return nil, errors.New("expected start of object")
	}
	var keys []string
	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		if t == json.Delim('}') {
			return keys, nil
		}
		keys = append(keys, t.(string))
		if err := skipValue(d); err != nil {
			return nil, err
		}
	}
}
func skipValue(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	switch t {
	case json.Delim('['), json.Delim('{'):
		for {
			if err := skipValue(d); err != nil {
				if err == end {
					break
				}
				return err
			}
		}
	case json.Delim(']'), json.Delim('}'):
		return end
	}
	return nil
}

var end = errors.New("invalid end of array or object")
