package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/micro/micro/v3/service"

	alertpb "github.com/m3o/services/alert/proto/alert"
	custpb "github.com/m3o/services/customers/proto"
	endtoend "github.com/m3o/services/endtoend/proto"
	"github.com/micro/micro/v3/service/client"
	mconfig "github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	mstore "github.com/micro/micro/v3/service/store"

	"github.com/google/uuid"
)

const (
	signupFrom          = "Micro Team <support@m3o.com>"
	signupSubject       = "Welcome to the M3O Platform"
	microInstallScript  = "https://install.m3o.com/micro"
	signupSuccessString = "Signup complete"
	keyOtp              = "otp"
	keyCheckResult      = "checkResult"

	// how fresh does a check need to be? cron runs every 5 mins and check takes just over 1min.
	checkBuffer = 7 * time.Minute

	microBinary = "/root/bin/micro"
)

var (
	otpRegexp = regexp.MustCompile("please copy and paste this one time token into your terminal:\\s*([a-zA-Z]*)\\s*This token expires")
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
		return errors.InternalServerError("endtoend.check.store", "Failed to load last result %s", err)
	}
	if len(recs) == 0 {
		return errors.InternalServerError("endtoend.check.noresults", "Failed to load last result, no results found")
	}
	cr := checkResult{}
	if err := json.Unmarshal(recs[0].Value, &cr); err != nil {
		return errors.InternalServerError("endtoend.check.unmarshal", "Failed to unmarshal last result %s", err)
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
	return errors.New("endtoend.chack.failed", response.Body, response.StatusCode)

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
		merr, ok := err.(*errors.Error)
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
	if rsp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("error completing signup %s %s", rsp.Status, string(b))
	}

	var custErr error
	loopStart = time.Now()
	for time.Now().Sub(loopStart) < 2*time.Minute {
		time.Sleep(10 * time.Second)
		rsp, err := e.custSvc.Read(context.TODO(), &custpb.ReadRequest{Email: e.email}, client.WithAuthToken())
		if err != nil {
			custErr = fmt.Errorf("error checking customer status %s", err)
			continue
		}
		if rsp.Customer.Status != "active" {
			custErr = fmt.Errorf("customer status should be active but is %s", rsp.Customer.Status)
			continue
		}
		custErr = nil
		break
	}

	if custErr != nil {
		return custErr
	}

	// TODO check balance for intro credit
	// TODO run some apis
	// TODO add credit via stripe

	return nil
}
