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
		email:   email,
		custSvc: custpb.NewCustomersService("customers", srv.Client()),
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
	return errors.New("endtoend.chack.failed", response.Body, response.StatusCode)

}

func (e *Endtoend) RunCheck(ctx context.Context, request *endtoend.Request, response *endtoend.Response) error {
	go func() error {
		if err := installMicro(); err != nil {
			log.Errorf("Error installing micro %s", err)
			return err
		}
		if err := e.signup(); err != nil {
			log.Errorf("Error during signup %s", err)
			return err
		}
		return nil
	}()
	return nil
}

func installMicro() error {
	// setup
	os.Remove("/tmp/micro")

	cmd := exec.Command("wget", microInstallScript)
	cmd.Dir = "/tmp"
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to get install script %s", err)
	}
	cmd = exec.Command("/bin/bash", "micro")
	cmd.Dir = "/tmp"
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to install micro %s", err)
	}
	return nil
}

func (e *Endtoend) signup() error {
	// reset, delete any existing customers
	cust, err := e.custSvc.Read(context.TODO(), &custpb.ReadRequest{Email: e.email}, client.WithAuthToken())
	if err != nil {
		return err
	}
	_, err = e.custSvc.Delete(context.TODO(), &custpb.DeleteRequest{Id: cust.Customer.Id}, client.WithAuthToken())
	if err != nil {
		return err
	}

	start := time.Now()
	cmd := exec.Command("/root/bin/micro", "signup", "--password", uuid.New().String())
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	chErr := make(chan error)
	go func() {
		defer close(chErr)
		outp, err := cmd.CombinedOutput()
		if err != nil {
			chErr <- err
		}
		if !strings.Contains(string(outp), "Finishing signup for") {
			chErr <- fmt.Errorf("Output does not contain success %s", string(outp))
		}
	}()
	go func() {
		time.Sleep(180 * time.Second)
		cmd.Process.Kill()
	}()

	time.Sleep(1 * time.Second)
	_, err = io.WriteString(stdin, e.email+"\n")
	if err != nil {
		return err
	}

	code := ""

	for i := 0; i < 10; i++ {
		time.Sleep(15 * time.Second)
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
		return fmt.Errorf("no OTP code found")
	}

	_, err = io.WriteString(stdin, code+"\n")
	if err != nil {
		return err
	}

	err = <-chErr
	var custErr error
	for i := 0; i < 5; i++ {
		time.Sleep(15 * time.Second)
		rsp, err := e.custSvc.Read(context.TODO(), &custpb.ReadRequest{Email: e.email}, client.WithAuthToken())
		if err != nil {
			custErr = err
			continue
		}
		if rsp.Customer.Status != "active" {
			custErr = fmt.Errorf("customer status is %s", rsp.Customer.Status)
			continue
		}
		custErr = nil
		break
	}
	result := checkResult{
		Time:   time.Now().Unix(),
		Passed: custErr == nil,
	}
	if custErr != nil {
		result.Error = custErr.Error()
	}
	b, _ := json.Marshal(result)

	mstore.Write(&mstore.Record{
		Key:   keyCheckResult,
		Value: b,
	})
	log.Infof("Signup took %d to complete", time.Now().Sub(start))
	return custErr
}
