// +build m3o

package signup

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro/micro/v3/test"
	"github.com/stripe/stripe-go/v71"
	stripe_client "github.com/stripe/stripe-go/v71/client"
)

const (
	retryCount          = 1
	signupSuccessString = "Signup complete"
)

func TestM3oSignupFlow(t *testing.T) {
	test.TrySuite(t, testM3oSignupFlow, retryCount)
}

func testM3oSignupFlow(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	envToConfigKey := map[string]string{
		"MICRO_STRIPE_API_KEY":       "micro.payments.stripe.api_key",
		"MICRO_SENDGRID_API_KEY":     "micro.signup.sendgrid.api_key",
		"MICRO_SENDGRID_TEMPLATE_ID": "micro.signup.sendgrid.template_id",
		"MICRO_STRIPE_PLAN_ID":       "micro.signup.plan_id",
		"MICRO_EMAIL_FROM":           "micro.signup.email_from",
		"MICRO_TEST_ENV":             "micro.signup.test_env",
	}

	for envKey, configKey := range envToConfigKey {
		val := os.Getenv(envKey)
		if len(val) == 0 {
			t.Fatalf("'%v' flag is missing", envKey)
		}
		outp, err := serv.Command().Exec("config", "set", configKey, val)
		if err != nil {
			t.Fatal(string(outp))
		}
	}

	outp, err := serv.Command().Exec("run", getSrcString("M3O_INVITE_SVC", "../../../invite"))
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	outp, err = serv.Command().Exec("run", getSrcString("M3O_SIGNUP_SVC", "../../../signup"))
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	outp, err = serv.Command().Exec("run", getSrcString("M3O_STRIPE_SVC", "../../../payments/provider/stripe"))
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	if err := test.Try("Find signup and stripe in list", t, func() ([]byte, error) {
		outp, err := serv.Command().Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "stripe") || !strings.Contains(string(outp), "signup") || !strings.Contains(string(outp), "invite") {
			return outp, errors.New("Can't find signup or stripe or invite in list")
		}
		return outp, err
	}, 180*time.Second); err != nil {
		return
	}

	// flags
	envFlag := "-e=" + serv.Env()
	confFlag := "-c=" + serv.Command().Config

	email := "dobronszki@gmail.com"

	time.Sleep(5 * time.Second)

	cmd := exec.Command("micro", envFlag, confFlag, "signup")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		outp, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("Expected an error for login but got none")
		} else if !strings.Contains(string(outp), "signup.notallowed") {
			t.Fatal(string(outp))
		}
		wg.Done()
	}()
	go func() {
		time.Sleep(20 * time.Second)
		cmd.Process.Kill()
	}()
	_, err = io.WriteString(stdin, email+"\n")
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	if t.Failed() {
		return
	}

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "create", "--email="+email)
	}, 5*time.Second)

	// Adjust rules before we signup into a non admin account
	outp, err = serv.Command().Exec("auth", "create", "rule", "--access=granted", "--scope=''", "--resource=\"service:invite:create\"", "invitecreatepublic")
	if err != nil {
		t.Fatalf("Error setting up rules: %v", err)
		return
	}

	password := "PassWord1@"
	signup(serv, t, email, password, false, false)
	if t.Failed() {
		return
	}
	outp, err = serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	ns := strings.TrimSpace(string(outp))

	if strings.Count(ns, "-") != 2 {
		t.Fatalf("Expected 2 dashes in namespace but namespace is: %v", ns)
		return
	}

	t.T().Logf("Namespace set is %v", ns)

	test.Try("Find account", t, func() ([]byte, error) {
		outp, err = serv.Command().Exec("auth", "list", "accounts")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), email) {
			return outp, errors.New("Account not found")
		}
		if strings.Contains(string(outp), "default") {
			return outp, errors.New("Default account should not be present in the namespace")
		}
		return outp, nil
	}, 5*time.Second)

	newEmail := "dobronszki+1@gmail.com"

	test.Login(serv, t, email, password)

	if err := test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "create", "--email="+newEmail, "--namespace="+ns)
	}, 7*time.Second); err != nil {
		t.Fatal(err)
		return
	}

	// Log out and switch namespace back to micro
	outp, err = serv.Command().Exec("user", "config", "set", "micro.auth."+serv.Env())
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	outp, err = serv.Command().Exec("user", "config", "set", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	// @todo: only needed because of logging endpoint etc not being open by default.
	// should create open rules instead
	err = test.Login(serv, t, "admin", "micro")
	if err != nil {
		t.Fatal(err)
		return
	}
	signup(serv, t, newEmail, password, true, true)
	if t.Failed() {
		return
	}
	outp, err = serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	newNs := strings.TrimSpace(string(outp))
	if newNs != ns {
		t.Fatalf("Namespaces should match, old: %v, new: %v", ns, newNs)
		return
	}

	t.T().Logf("Namespace joined: %v", string(outp))
}

func signup(serv test.Server, t *test.T, email, password string, isInvited, shouldJoin bool) {
	envFlag := "-e=" + serv.Env()
	confFlag := "-c=" + serv.Command().Config

	cmd := exec.Command("micro", envFlag, confFlag, "signup", "--password", password)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		outp, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(email+": "+string(outp), err)
			return
		}
		if !strings.Contains(string(outp), signupSuccessString) {
			t.Fatal(string(outp))
			return
		}
	}()
	go func() {
		time.Sleep(25 * time.Second)
		cmd.Process.Kill()
	}()

	time.Sleep(1 * time.Second)
	_, err = io.WriteString(stdin, email+"\n")
	if err != nil {
		t.Fatal(err)
		return
	}

	code := ""
	// careful: there might be multiple codes in the logs
	codes := []string{}
	time.Sleep(2 * time.Second)

	t.Log("looking for code now", email)
	if err := test.Try("Find latest verification token in logs", t, func() ([]byte, error) {
		outp, err := serv.Command().Exec("logs", "-n", "300", "signup")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), email) {
			return outp, errors.New("Output does not contain email")
		}
		if !strings.Contains(string(outp), "Sending verification token") {
			return outp, errors.New("Output does not contain expected")
		}
		for _, line := range strings.Split(string(outp), "\n") {
			if strings.Contains(line, "Sending verification token") {
				codes = append(codes, strings.Split(line, "'")[1])
			}
		}
		return outp, nil
	}, 15*time.Second); err != nil {
		return
	}

	if len(codes) == 0 {
		t.Fatal("No code found")
		return
	}
	code = codes[len(codes)-1]

	t.Log("Code is ", code, " for email ", email)
	if code == "" {
		t.Fatal("Code not found")
		return
	}
	_, err = io.WriteString(stdin, code+"\n")
	if err != nil {
		t.Fatal(err)
		return
	}

	if isInvited {
		time.Sleep(3 * time.Second)
		answer := "own"
		if shouldJoin {
			t.Log("Joining a namespace now")
			answer = "join"
		}
		_, err = io.WriteString(stdin, answer+"\n")
		if err != nil {
			t.Fatal(err)
			return
		}
	}

	time.Sleep(5 * time.Second)

	sc := stripe_client.New(os.Getenv("MICRO_STRIPE_API_KEY"), nil)
	pm, err := sc.PaymentMethods.New(
		&stripe.PaymentMethodParams{
			Card: &stripe.PaymentMethodCardParams{
				Number:   stripe.String("4242424242424242"),
				ExpMonth: stripe.String("7"),
				ExpYear:  stripe.String("2021"),
				CVC:      stripe.String("314"),
			},
			Type: stripe.String("card"),
		})
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = io.WriteString(stdin, pm.ID+"\n")
	if err != nil {
		t.Fatal(err)
	}

	// Don't wait if a test is already failed, this is a quirk of the
	// test framework @todo fix this quirk
	if t.Failed() {
		return
	}
	wg.Wait()
}

func getSrcString(envvar, dflt string) string {
	if env := os.Getenv(envvar); env != "" {
		return env
	}
	return dflt
}
