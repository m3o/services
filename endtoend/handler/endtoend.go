package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	if cr.Passed && time.Now().Add(-5*time.Minute).Unix() < cr.Time {
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
		e.alertSvc.ReportEvent(context.TODO(), &alertpb.ReportEventRequest{
			Event: &alertpb.Event{
				Category: "monitoring",
				Action:   "signup",
				Label:    "endtoend",
				Value:    1,
				Metadata: map[string]string{"error": err.Error()},
			},
		}, client.WithAuthToken())
	}()
	if err = installMicro(); err != nil {
		log.Errorf("Error installing micro %s", err)
		return err
	}
	if err = e.signup(); err != nil {
		log.Errorf("Error during signup %s", err)
		return err
	}
	return nil
}

func installMicro() error {
	// setup
	os.Remove("/tmp/micro")

	cmd := exec.Command("wget", microInstallScript)
	cmd.Dir = "/tmp"
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get install script %s", err)
	}
	cmd = exec.Command("/bin/bash", "micro")
	cmd.Dir = "/tmp"
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install micro %s", err)
	}
	return nil
}

func (e *Endtoend) signup() error {
	// reset, delete any existing customers
	_, err := e.custSvc.Delete(context.TODO(), &custpb.DeleteRequest{Email: e.email, Force: true}, client.WithAuthToken())
	if err != nil {
		return fmt.Errorf("error while cleaning up existing customer %s", err)
	}

	start := time.Now()
	cmd := exec.Command(microBinary, "signup", "--password", uuid.New().String())
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error running micro signup command %s", err)
	}
	chErr := make(chan error)
	go func() {
		defer close(chErr)
		outp, err := cmd.CombinedOutput()
		if err != nil {
			chErr <- err
		}
		if !strings.Contains(string(outp), signupSuccessString) {
			chErr <- fmt.Errorf("micro signup output does not contain success %s", string(outp))
		}
	}()
	go func() {
		time.Sleep(180 * time.Second)
		cmd.Process.Kill()
	}()

	time.Sleep(1 * time.Second)
	_, err = io.WriteString(stdin, e.email+"\n")
	if err != nil {
		return fmt.Errorf("error inputting email to micro command %s", err)
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

	_, err = io.WriteString(stdin, code+"\n")
	if err != nil {
		return fmt.Errorf("error inputting otp code to micro command %s", err)
	}

	err = <-chErr
	if err != nil {
		log.Errorf("Error while completing signup %s", err)
		return fmt.Errorf("error completing signup %s", err)
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

	// run a service
	cmd = exec.Command(microBinary, "run", "github.com/micro/services/helloworld")
	outp, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running helloworld %s, %s", err, string(outp))
	}
	// wait for it to run
	var statusErr error
	loopStart = time.Now()
	for time.Now().Sub(loopStart) < 2*time.Minute {
		time.Sleep(10 * time.Second)
		cmd = exec.Command(microBinary, "status")
		outp, err := cmd.CombinedOutput()
		if err != nil {
			statusErr = fmt.Errorf("error checking helloworld status %s, %s", err, string(outp))
		}
		if strings.Contains(string(outp), "running") {
			statusErr = nil
			break
		}
		statusErr = fmt.Errorf("helloworld not in running state. Output of status command %s", string(outp))
	}
	if statusErr != nil {
		return statusErr
	}
	return nil
}
