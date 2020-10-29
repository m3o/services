// +build m3o

package signup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// generates test emails
func testEmail(nth int) string {
	uid := randStringRunes(8)
	if nth == 0 {
		return fmt.Sprintf("platform+citest+%v@m3o.com", uid)
	}
	return fmt.Sprintf("platform+citest+%v+%v@m3o.com", nth, uid)
}

func TestSignupFlow(t *testing.T) {
	test.TrySuite(t, testSignupFlow, retryCount)
}

func setupM3Tests(serv test.Server, t *test.T) {
	setupM3TestsImpl(serv, t, false)
}

func setupFreeM3Tests(serv test.Server, t *test.T) {
	setupM3TestsImpl(serv, t, true)
}

func setupM3TestsImpl(serv test.Server, t *test.T, freeTier bool) {
	envToConfigKey := map[string][]string{
		"MICRO_STRIPE_API_KEY":                      {"micro.payments.stripe.api_key"},
		"MICRO_SENDGRID_API_KEY":                    {"micro.emails.sendgrid.api_key"},
		"MICRO_SENDGRID_TEMPLATE_ID":                {"micro.signup.sendgrid.template_id"},
		"MICRO_SENDGRID_INVITE_TEMPLATE_ID":         {"micro.invite.sendgrid.invite_template_id"},
		"MICRO_STRIPE_PLAN_ID":                      {"micro.subscriptions.plan_id"},
		"MICRO_STRIPE_ADDITIONAL_USERS_PRICE_ID":    {"micro.subscriptions.additional_users_price_id"},
		"MICRO_EMAIL_FROM":                          {"micro.signup.email_from"},
		"MICRO_TEST_ENV":                            {"micro.signup.test_env"},
		"MICRO_STRIPE_ADDITIONAL_SERVICES_PRICE_ID": {"micro.subscriptions.additional_services_price_id"},
	}

	if err := test.Try("Set up config values", t, func() ([]byte, error) {
		for envKey, configKeys := range envToConfigKey {
			val := os.Getenv(envKey)
			if len(val) == 0 {
				t.Fatalf("'%v' flag is missing", envKey)
			}
			for _, configKey := range configKeys {
				outp, err := serv.Command().Exec("config", "set", configKey, val)
				if err != nil {
					return outp, err
				}
			}
		}
		if freeTier {
			outp, err := serv.Command().Exec("config", "set", "micro.signup.no_payment", "true")
			if err != nil {
				return outp, err
			}
			outp, err = serv.Command().Exec("config", "set", "micro.signup.message", "Finishing signup for %s")
			if err != nil {
				return outp, err
			}

		}
		return serv.Command().Exec("config", "set", "micro.billing.max_included_services", "3")
	}, 10*time.Second); err != nil {
		t.Fatal(err)
		return
	}

	services := []struct {
		envVar string
		deflt  string
	}{
		{envVar: "M3O_INVITE_SVC", deflt: "../../../invite"},
		{envVar: "M3O_SIGNUP_SVC", deflt: "../../../signup"},
		{envVar: "M3O_STRIPE_SVC", deflt: "../../../payments"},
		{envVar: "M3O_CUSTOMERS_SVC", deflt: "../../../customers"},
		{envVar: "M3O_NAMESPACES_SVC", deflt: "../../../namespaces"},
		{envVar: "M3O_SUBSCRIPTIONS_SVC", deflt: "../../../subscriptions"},
		{envVar: "M3O_PLATFORM_SVC", deflt: "../../../platform"},
		{envVar: "M3O_EMAILS_SVC", deflt: "../../../emails"},
	}

	for _, v := range services {
		outp, err := serv.Command().Exec("run", getSrcString(v.envVar, v.deflt))
		if err != nil {
			t.Fatal(string(outp))
			return
		}
	}

	if err := test.Try("Find signup, invite and stripe in list", t, func() ([]byte, error) {
		outp, err := serv.Command().Exec("services")
		if err != nil {
			return outp, err
		}
		list := []string{"payments", "signup", "invite", "emails", "customers"}
		logOutp := []byte{}
		fail := false
		for _, s := range list {
			if !strings.Contains(string(outp), s) {
				o, _ := serv.Command().Exec("logs", s)
				logOutp = append(logOutp, o...)
				fail = true
			}
		}
		if fail {
			return append(outp, logOutp...), errors.New("Can't find required services in list")
		}
		return outp, err
	}, 180*time.Second); err != nil {
		return
	}

	// setup rules

	// Adjust rules before we signup into a non admin account
	outp, err := serv.Command().Exec("auth", "create", "rule", "--access=granted", "--scope=''", "--resource=\"service:invite:*\"", "invite")
	if err != nil {
		t.Fatalf("Error setting up rules: %v", string(outp))
		return
	}

	// Adjust rules before we signup into a non admin account
	outp, err = serv.Command().Exec("auth", "create", "rule", "--access=granted", "--scope=''", "--resource=\"service:signup:*\"", "signup")
	if err != nil {
		t.Fatalf("Error setting up rules: %v", string(outp))
		return
	}

	// Adjust rules before we signup into a non admin account
	outp, err = serv.Command().Exec("auth", "create", "rule", "--access=granted", "--scope=''", "--resource=\"service:auth:*\"", "auth")
	if err != nil {
		t.Fatalf("Error setting up rules: %v", string(outp))
		return
	}

	// copy the config with the admin logged in so we can use it for reading logs
	// we dont want to have an open access rule for logs as it's not how it works in live
	confPath := serv.Command().Config
	outp, err = exec.Command("cp", "-rf", confPath, confPath+".admin").CombinedOutput()
	if err != nil {
		t.Fatalf("Error copying config: %v", outp)
		return
	}
}

func logout(serv test.Server, t *test.T) {
	// Log out and switch namespace back to micro
	outp, err := serv.Command().Exec("user", "config", "set", "micro.auth."+serv.Env())
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	outp, err = serv.Command().Exec("user", "config", "set", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
}

func testSignupFlow(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)

	email := testEmail(0)

	time.Sleep(5 * time.Second)

	// Log out of the admin account to start testing signups
	logout(serv, t)

	password := "PassWord1@"
	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})
	if t.Failed() {
		return
	}
	t.Logf("Signup 1 complete %s", time.Now())
	outp, err := serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
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

	newEmail := testEmail(1)
	newEmail2 := testEmail(2)

	test.Login(serv, t, email, password)

	if err := test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+newEmail, "--namespace="+ns)
	}, 7*time.Second); err != nil {
		t.Fatal(err)
		return
	}
	if err := test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+newEmail2, "--namespace="+ns)
	}, 7*time.Second); err != nil {
		t.Fatal(err)
		return
	}

	logout(serv, t)

	signup(serv, t, newEmail, password, signupOptions{inviterEmail: email, xthInvitee: 1, isInvitedToNamespace: true, shouldJoin: true})
	if t.Failed() {
		return
	}
	t.Logf("Signup 2 complete %s", time.Now())
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

	logout(serv, t)

	signup(serv, t, newEmail2, password, signupOptions{inviterEmail: email, xthInvitee: 2, isInvitedToNamespace: true, shouldJoin: true})
	t.Logf("Signup 3 complete %s", time.Now())
	if t.Failed() {
		return
	}
	outp, err = serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	newNs = strings.TrimSpace(string(outp))
	if newNs != ns {
		t.Fatalf("Namespaces should match, old: %v, new: %v", ns, newNs)
		return
	}

	t.T().Logf("Namespace joined: %v", string(outp))
}

func TestAdminInvites(t *testing.T) {
	test.TrySuite(t, testAdminInvites, retryCount)
}

func testAdminInvites(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	time.Sleep(2 * time.Second)

	logout(serv, t)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	outp, err := serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	ns := strings.TrimSpace(string(outp))
	if ns == "micro" {
		t.Fatal("SECURITY FLAW: invited user ended up in micro namespace")
	}
	if strings.Count(ns, "-") != 2 {
		t.Fatalf("Expected 2 dashes in namespace but namespace is: %v", ns)
		return
	}

	t.T().Logf("Namespace joined: %v", string(outp))
}

func TestAdminInviteNoLimit(t *testing.T) {
	test.TrySuite(t, testAdminInviteNoLimit, retryCount)
}

func testAdminInviteNoLimit(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)

	// Make sure test mod is on otherwise this will spam
	for i := 0; i < 10; i++ {
		test.Try("Send invite", t, func() ([]byte, error) {
			return serv.Command().Exec("invite", "user", "--email="+testEmail(i))
		}, 5*time.Second)
	}
}

func TestUserInviteLimit(t *testing.T) {
	test.TrySuite(t, testUserInviteLimit, retryCount)
}

func testUserInviteLimit(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	logout(serv, t)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	// Make sure test mod is on otherwise this will spam
	for i := 0; i < 5; i++ {
		test.Try("Send invite", t, func() ([]byte, error) {
			return serv.Command().Exec("invite", "user", "--email="+testEmail(i+1))
		}, 5*time.Second)
	}

	outp, err := serv.Command().Exec("invite", "user", "--email="+testEmail(7))
	if err == nil {
		t.Fatalf("Sending 6th invite should fail: %v", outp)
	}
}

func TestUserInviteNoJoin(t *testing.T) {
	test.TrySuite(t, testUserInviteNoJoin, retryCount)
}

func testUserInviteNoJoin(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	logout(serv, t)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	outp, err := serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	ns := strings.TrimSpace(string(outp))
	if strings.Count(ns, "-") != 2 {
		t.Fatalf("Expected 2 dashes in namespace but namespace is: %v", ns)
		return
	}

	newEmail := testEmail(1)

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+newEmail)
	}, 5*time.Second)

	logout(serv, t)

	signup(serv, t, newEmail, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	outp, err = serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	newNs := strings.TrimSpace(string(outp))
	if strings.Count(newNs, "-") != 2 {
		t.Fatalf("Expected 2 dashes in namespace but namespace is: %v", ns)
		return
	}

	if ns == newNs {
		t.Fatal("User should not have joined invitees namespace")
	}
}

func TestUserInviteJoinDecline(t *testing.T) {
	test.TrySuite(t, testUserInviteJoinDecline, retryCount)
}

func testUserInviteJoinDecline(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	logout(serv, t)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	outp, err := serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	ns := strings.TrimSpace(string(outp))
	if strings.Count(ns, "-") != 2 {
		t.Fatalf("Expected 2 dashes in namespace but namespace is: %v", ns)
		return
	}

	newEmail := testEmail(1)

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+newEmail, "--namespace="+ns)
	}, 5*time.Second)

	logout(serv, t)

	signup(serv, t, newEmail, password, signupOptions{isInvitedToNamespace: true, shouldJoin: false})

	outp, err = serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	newNs := strings.TrimSpace(string(outp))
	if strings.Count(newNs, "-") != 2 {
		t.Fatalf("Expected 2 dashes in namespace but namespace is: %v", ns)
		return
	}

	if ns == newNs {
		t.Fatal("User should not have joined invitees namespace")
	}
}

func TestUserInviteToNotOwnedNamespace(t *testing.T) {
	test.TrySuite(t, testUserInviteToNotOwnedNamespace, retryCount)
}

func testUserInviteToNotOwnedNamespace(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	logout(serv, t)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	outp, err := serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	ns := strings.TrimSpace(string(outp))
	if strings.Count(ns, "-") != 2 {
		t.Fatalf("Expected 2 dashes in namespace but namespace is: %v", ns)
		return
	}

	newEmail := testEmail(1)

	outp, err = serv.Command().Exec("invite", "user", "--email="+newEmail, "--namespace=not-my-namespace")
	if err == nil {
		t.Fatalf("Should not be able to invite to an unowned namespace, output: %v", string(outp))
	}

	// Testing for micro namespace just to be sure as it's the worst case
	outp, err = serv.Command().Exec("invite", "user", "--email="+newEmail, "--namespace=micro")
	if err == nil {
		t.Fatalf("Should not be able to invite to an unowned namespace, output: %v", string(outp))
	}
}

func TestServicesSubscription(t *testing.T) {
	test.TrySuite(t, testServicesSubscription, retryCount)
}

func testServicesSubscription(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	serv.Command().Exec("run", "github.com/micro/services/helloworld")
	serv.Command().Exec("run", "github.com/micro/services/blog/posts")
	serv.Command().Exec("run", "github.com/micro/services/blog/tags")
	serv.Command().Exec("run", "github.com/micro/services/test/pubsub")

	test.Try("Wait for services", t, func() ([]byte, error) {
		outp, err := serv.Command().Exec("status")
		if !strings.Contains(string(outp), "helloworld") || !strings.Contains(string(outp), "posts") || !strings.Contains(string(outp), "posts") ||
			!strings.Contains(string(outp), "tags") || !strings.Contains(string(outp), "pubsub") {
			return outp, errors.New("Can't find services")
		}
		return outp, err
	}, 30*time.Second)

	adminConfFlag := "-c=" + serv.Command().Config + ".admin"
	envFlag := "-e=" + serv.Env()
	exec.Command("micro", envFlag, adminConfFlag, "run", "../../../usage").CombinedOutput()
	time.Sleep(4 * time.Second)
	exec.Command("micro", envFlag, adminConfFlag, "run", "../../../billing").CombinedOutput()

	customerId := ""
	test.Try("Get changes", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", envFlag, adminConfFlag, "billing", "updates").CombinedOutput()
		outp1, _ := exec.Command("micro", envFlag, adminConfFlag, "logs", "billing").CombinedOutput()
		fulloutp := append(outp, outp1...)
		if err != nil {
			return fulloutp, err
		}
		updatesRsp := map[string]interface{}{}
		err = json.Unmarshal(outp, &updatesRsp)
		if err != nil {
			return fulloutp, err
		}
		updates, ok := updatesRsp["updates"].([]interface{})
		if !ok {
			return fulloutp, errors.New("Unexpected output")
		}
		if len(updates) == 0 {
			return fulloutp, errors.New("No updates found")
		}
		if updates[0].(map[string]interface{})["quantityTo"].(string) != "1" {
			return fulloutp, errors.New("Quantity should be 1")
		}
		customerId = updates[0].(map[string]interface{})["customerID"].(string)
		if !strings.Contains(string(outp), "Additional services") {
			return fulloutp, errors.New("unexpected output")
		}
		if strings.Contains(string(outp), "Additional users") {
			return fulloutp, errors.New("unexpected output")
		}
		return fulloutp, err
	}, 90*time.Second)

	test.Try("Apply change", t, func() ([]byte, error) {
		return exec.Command("micro", envFlag, adminConfFlag, "billing", "apply", "--customerID="+customerId).CombinedOutput()
	}, 5*time.Second)

	time.Sleep(4 * time.Second)
	subs := getSubscriptions(t, email)
	priceID := os.Getenv("MICRO_STRIPE_ADDITIONAL_SERVICES_PRICE_ID")
	sub, ok := subs[priceID]
	if !ok {
		t.Fatalf("Sub with id %v not found", priceID)
		return
	}
	if sub.Quantity != 1 {
		t.Fatalf("Quantity should be 1, but it's %v", sub.Quantity)
	}

	test.Try("Get changes again", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", envFlag, adminConfFlag, "billing", "updates").CombinedOutput()
		outp1, _ := exec.Command("micro", envFlag, adminConfFlag, "logs", "billing").CombinedOutput()
		fulloutp := append(outp, outp1...)
		if err != nil {
			return fulloutp, err
		}
		updatesRsp := map[string]interface{}{}
		err = json.Unmarshal(outp, &updatesRsp)
		if err != nil {
			return fulloutp, err
		}
		updates, ok := updatesRsp["updates"].([]interface{})
		if ok && len(updates) > 0 {
			return fulloutp, errors.New("Updates found when there should be none")
		}
		return fulloutp, err
	}, 20*time.Second)
}

func TestUsersSubscription(t *testing.T) {
	test.TrySuite(t, testUsersSubscription, retryCount)
}

func testUsersSubscription(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})

	serv.Command().Exec("auth", "create", "account", "create", "hi@there.com")

	adminConfFlag := "-c=" + serv.Command().Config + ".admin"
	envFlag := "-e=" + serv.Env()
	exec.Command("micro", envFlag, adminConfFlag, "kill", "billing").CombinedOutput()
	time.Sleep(2 * time.Second)
	exec.Command("micro", envFlag, adminConfFlag, "run", "../../../usage").CombinedOutput()
	time.Sleep(4 * time.Second)
	exec.Command("micro", envFlag, adminConfFlag, "run", "../../../billing").CombinedOutput()

	customerId := ""
	test.Try("Get changes", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", envFlag, adminConfFlag, "billing", "updates").CombinedOutput()
		outp1, _ := exec.Command("micro", envFlag, adminConfFlag, "logs", "billing").CombinedOutput()
		fulloutp := append(outp, outp1...)
		if err != nil {
			return fulloutp, err
		}
		updatesRsp := map[string]interface{}{}
		err = json.Unmarshal(outp, &updatesRsp)
		if err != nil {
			return outp, err
		}
		updates, ok := updatesRsp["updates"].([]interface{})
		if !ok {
			return outp, errors.New("Unexpected output")
		}
		if len(updates) == 0 {
			return outp, errors.New("No updates found")
		}
		if updates[0].(map[string]interface{})["quantityTo"].(string) != "1" {
			return outp, errors.New("Quantity should be 1")
		}
		customerId = updates[0].(map[string]interface{})["customerID"].(string)
		if !strings.Contains(string(outp), "Additional users") {
			return outp, errors.New("unexpected output")
		}
		if strings.Contains(string(outp), "Additional services") {
			return outp, errors.New("unexpected output")
		}
		return outp, err
	}, 40*time.Second)

	test.Try("Apply change", t, func() ([]byte, error) {
		return exec.Command("micro", envFlag, adminConfFlag, "billing", "apply", "--customerID="+customerId).CombinedOutput()
	}, 5*time.Second)

	time.Sleep(4 * time.Second)
	subs := getSubscriptions(t, email)
	priceID := os.Getenv("MICRO_STRIPE_ADDITIONAL_USERS_PRICE_ID")
	sub, ok := subs[priceID]
	if !ok {
		t.Fatalf("Sub with id %v not found", priceID)
		return
	}
	if sub.Quantity != 1 {
		t.Fatalf("Quantity should be 1, but it's %v", sub.Quantity)
	}
}

// returns map witj plan (price) id -> subscriptions
func getSubscriptions(t *test.T, email string) map[string]*stripe.Subscription {
	sc := stripe_client.New(os.Getenv("MICRO_STRIPE_API_KEY"), nil)
	subListParams := &stripe.SubscriptionListParams{}
	subListParams.Limit = stripe.Int64(20)
	subListParams.AddExpand("data.customer")
	iter := sc.Subscriptions.List(subListParams)
	count := 0
	// email -> plan/price id -> subscription
	plans := map[string]*stripe.Subscription{}
	for iter.Next() {
		if count > 20 {
			break
		}
		count++

		c := iter.Subscription()
		if c.Customer.Email == email {
			if c.Plan != nil {
				plans[c.Plan.ID] = c
			}
		}
	}
	return plans
}

type signupOptions struct {
	isInvitedToNamespace bool
	shouldJoin           bool
	inviterEmail         string
	xthInvitee           int
	freeTier             bool
}

func signup(serv test.Server, t *test.T, email, password string, opts signupOptions) {
	envFlag := "-e=" + serv.Env()
	confFlag := "-c=" + serv.Command().Config
	adminConfFlag := "-c=" + serv.Command().Config + ".admin"

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
			t.Fatal(string(outp), err)
			return
		}
		if !strings.Contains(string(outp), signupSuccessString) {
			t.Fatal(string(outp))
			return
		}
		if !opts.shouldJoin {
			if !opts.freeTier && !strings.Contains(string(outp), "Please complete signup at") {
				t.Fatal(string(outp))
				return
			}
			if opts.freeTier && !strings.Contains(string(outp), "Finishing signup for") {
				t.Fatal(string(outp))
				return
			}
		}
	}()
	go func() {
		time.Sleep(60 * time.Second)
		t.Logf("Killing process")
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
		outp, err := exec.Command("micro", envFlag, adminConfFlag, "logs", "-n", "300", "signup").CombinedOutput()
		if err != nil {
			return outp, err
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

	if opts.isInvitedToNamespace {
		time.Sleep(3 * time.Second)
		answer := "own"
		if opts.shouldJoin {
			t.Log("Joining a namespace now")
			answer = "join"
		}
		_, err = io.WriteString(stdin, answer+"\n")
		if err != nil {
			t.Fatal(err)
			return
		}
	}

	if !opts.shouldJoin && !opts.freeTier {
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
		t.Log("Added a new payment method to Stripe")

		// using a curl here as `call` redirection to micro namespace doesnt work properly, unlike
		// dynamic commands

		curl := func(serv test.Server, path, email, paymentMethod string) (map[string]interface{}, error) {
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%v/%v?email=%v&paymentMethod=%v", serv.APIPort(), path, url.QueryEscape(email), paymentMethod))
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			m := map[string]interface{}{}
			return m, json.Unmarshal(body, &m)
		}

		rsp, err := curl(serv, "signup/setPaymentMethod", email, pm.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(rsp) > 0 {
			t.Fatal(rsp)
		}
		t.Log("Added a new payment method to M3O")

	}

	if !opts.freeTier {

		// Some gotchas for this: while the stripe api documentation online
		// shows prices and plans being separate entities, even v71 version of the
		// go library only has plans. However, it seems like the prices are under plans too.
		test.Try("Check subscription in stripe", t, func() ([]byte, error) {
			sc := stripe_client.New(os.Getenv("MICRO_STRIPE_API_KEY"), nil)
			subListParams := &stripe.SubscriptionListParams{}
			subListParams.Limit = stripe.Int64(20)
			subListParams.AddExpand("data.customer")
			iter := sc.Subscriptions.List(subListParams)
			count := 0
			// email -> plan/price id -> subscription
			userPlans := map[string]*stripe.Subscription{}
			inviterPlans := map[string]*stripe.Subscription{}
			for iter.Next() {
				if count > 20 {
					break
				}
				count++

				c := iter.Subscription()
				if len(opts.inviterEmail) > 0 && c.Customer.Email == opts.inviterEmail {
					if c.Plan != nil {
						inviterPlans[c.Plan.ID] = c
					}
				}
				if c.Customer.Email == email {
					if c.Plan != nil {
						userPlans[c.Plan.ID] = c
					}
				}
			}

			if opts.shouldJoin {
				priceID := os.Getenv("MICRO_STRIPE_ADDITIONAL_USERS_PRICE_ID")
				sub, found := inviterPlans[priceID]
				if !found {
					return nil, fmt.Errorf("Subscription with price ID %v not found", priceID)
				}
				if sub.Quantity != int64(opts.xthInvitee) {
					return nil, fmt.Errorf("Subscription quantity '%v' should match invitee number '%v", sub.Quantity, opts.xthInvitee)
				}
			} else {
				planID := os.Getenv("MICRO_STRIPE_PLAN_ID")
				sub, found := userPlans[planID]
				if !found {
					return nil, fmt.Errorf("Subscription with plan ID %v not found", planID)
				}
				if sub.Quantity != 1 {
					return nil, fmt.Errorf("Subscription quantity should be 1 but is %v", sub.Quantity)
				}
			}
			t.Logf("Subscription checked")
			return nil, nil
		}, 50*time.Second)
	}
	test.Try("Check customer marked active", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", envFlag, adminConfFlag, "customers", "read", "--email="+email).CombinedOutput()
		if err != nil {
			t.Logf("Error checking customer status %s %s", string(outp), err)
			return outp, err
		}
		if !strings.Contains(string(outp), `"status": "active"`) {
			outp, _ = exec.Command("micro", envFlag, adminConfFlag, "logs", "customers").CombinedOutput()
			return outp, fmt.Errorf("Customer status is not active")
		}
		t.Logf("Customer marked active")
		return nil, nil
	}, 60*time.Second)

	// Don't wait if a test is already failed, this is a quirk of the
	// test framework @todo fix this quirk
	if t.Failed() {
		t.Logf("Failed signup")
		return
	}
	t.Logf("Waiting at end of signup")
	wg.Wait()
	t.Logf("Signup complete for %s", email)

}

func getSrcString(envvar, dflt string) string {
	if env := os.Getenv(envvar); env != "" {
		return env
	}
	return dflt
}

func TestDuplicateInvites(t *testing.T) {
	test.TrySuite(t, testDuplicateInvites, retryCount)
}

func testDuplicateInvites(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)
	test.Try("Send invite again", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)
	outp, err := serv.Command().Exec("logs", "invite")
	if err != nil {
		t.Fatalf("Unexpected error retrieving logs %s", err)
	}
	if !strings.Contains(string(outp), "Invite already sent for user "+email) {
		t.Fatalf("Invite was sent multiple times")
	}

	// test a force resend
	email2 := testEmail(1)
	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email2)
	}, 5*time.Second)
	test.Try("Send invite again", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email2, "--resend=true")
	}, 5*time.Second)
	outp, err = serv.Command().Exec("logs", "invite")
	if err != nil {
		t.Fatalf("Unexpected error retrieving logs %s", err)
	}
	if strings.Contains(string(outp), "Invite already sent for user "+email2) {
		t.Fatalf("Invite should have been sent multiple times but was blocked")
	}

}

func TestInviteEmailValidation(t *testing.T) {
	test.TrySuite(t, testInviteEmailValidation, retryCount)
}

func testInviteEmailValidation(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)

	outp, _ := serv.Command().Exec("invite", "user", "--email=notanemail.com")
	if !strings.Contains(string(outp), "400") {
		t.Fatalf("Expected a 400 bad request error %s", string(outp))
	}

}

func TestSubCancellation(t *testing.T) {
	test.TrySuite(t, testSubCancellation, retryCount)
}

func testSubCancellation(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupM3Tests(serv, t)
	email := testEmail(0)
	password := "PassWord1@"

	test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+email)
	}, 5*time.Second)

	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})
	if t.Failed() {
		return
	}

	outp, err := serv.Command().Exec("user", "config")
	if err != nil {
		t.Fatalf("Error getting user config %s %s", string(outp), err)
	}

	ns := ""
	for _, v := range strings.Split(string(outp), "\n") {
		if !strings.HasPrefix(v, "namespace: ") {
			continue
		}
		ns = strings.TrimPrefix(v, "namespace: ")
	}
	if len(ns) == 0 {
		t.Fatalf("Unable to determine the namespace of the user %s", string(outp))
	}

	adminConfFlag := "-c=" + serv.Command().Config + ".admin"
	envFlag := "-e=" + serv.Env()
	outp, err = exec.Command("micro", envFlag, adminConfFlag, "customers", "read", "--email="+email).CombinedOutput()
	if err != nil {
		t.Fatalf("Error looking up customer ID %s %s ", string(outp), err)
	}
	type cs struct {
		Id string `json:"id"`
	}
	type rsp struct {
		Customer cs `json:"customer"`
	}
	csObj := &rsp{}
	if err := json.Unmarshal(outp, csObj); err != nil {
		t.Fatalf("Error unmarshalling customer %s %s ", string(outp), err)
	}
	outp, err = exec.Command("micro", envFlag, adminConfFlag, "subscriptions", "cancel", "--customerID="+csObj.Customer.Id).CombinedOutput()
	if err != nil {
		t.Fatalf("Error cancelling %+v %s %s ", csObj, string(outp), err)
	}
	// check stripe for sub cancelled
	subs := getSubscriptions(t, email)
	for _, v := range subs {
		if v.Status != stripe.SubscriptionStatusCanceled {
			t.Fatalf("Subscription was not cancelled %+v", v)
		}
	}
	test.Try("Check customer deleted", t, func() ([]byte, error) {
		// check customer is deleted which happens async
		outp, err = exec.Command("micro", envFlag, adminConfFlag, "customers", "read", "--email="+email).CombinedOutput()
		if !strings.Contains(string(outp), "not found") {
			return outp, fmt.Errorf("Customer should not be found %s", err)
		}
		return nil, nil
	}, 60*time.Second)
	// check namespace is gone
	outp, err = exec.Command("micro", envFlag, adminConfFlag, "namespaces", "list", "--user="+csObj.Customer.Id).CombinedOutput()
	if strings.Contains(string(outp), "namespaces") {
		t.Fatalf("Customer should not have any namespaces %s %s", string(outp), err)
	}
	// check auth is gone
	outp, err = exec.Command("micro", envFlag, adminConfFlag, "store", "list", "--table", "auth").CombinedOutput()
	if strings.Contains(string(outp), ns) {
		t.Fatalf("Customer should not have any auth remnants %s %s", string(outp), err)
	}

	// try the signup again with this email. should succeed if we've deleted everything properly
	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false})
	if t.Failed() {
		return
	}

	outp, err = serv.Command().Exec("user", "config")
	if err != nil {
		t.Fatalf("Error getting user config %s %s", string(outp), err)
	}
	// check if we've correctly given a new namespace, would be bad if we gave them the same
	newNS := ""
	for _, v := range strings.Split(string(outp), "\n") {
		if !strings.HasPrefix(v, "namespace: ") {
			continue
		}
		newNS = strings.TrimPrefix(v, "namespace: ")
	}
	if len(newNS) == 0 {
		t.Fatalf("Unable to determine the namespace of the user %s", string(outp))
	}
	if newNS == ns {
		t.Fatalf("Error, we've reassigned an old namespace %s", ns)
	}
	// check we have a new customer ID
	outp, err = exec.Command("micro", envFlag, adminConfFlag, "customers", "read", "--email="+email).CombinedOutput()
	if err != nil {
		t.Fatalf("Error looking up customer ID %s %s ", string(outp), err)
	}
	newCsObj := &rsp{}
	if err := json.Unmarshal(outp, newCsObj); err != nil {
		t.Fatalf("Error unmarshalling customer %s %s ", string(outp), err)
	}
	if newCsObj.Customer.Id == csObj.Customer.Id {
		t.Fatalf("Error, we've reassigned an old customerID %s", ns)
	}

}

func TestFreeSignupFlow(t *testing.T) {
	test.TrySuite(t, testFreeSignupFlow, retryCount)
}

func testFreeSignupFlow(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupFreeM3Tests(serv, t)

	email := testEmail(0)

	time.Sleep(5 * time.Second)

	// Log out of the admin account to start testing signups
	logout(serv, t)

	password := "PassWord1@"
	signup(serv, t, email, password, signupOptions{isInvitedToNamespace: false, shouldJoin: false, freeTier: true})
	if t.Failed() {
		return
	}
	t.Logf("Signup 1 complete %s", time.Now())
	outp, err := serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
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

	newEmail := testEmail(1)
	newEmail2 := testEmail(2)

	test.Login(serv, t, email, password)

	if err := test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+newEmail, "--namespace="+ns)
	}, 7*time.Second); err != nil {
		t.Fatal(err)
		return
	}
	if err := test.Try("Send invite", t, func() ([]byte, error) {
		return serv.Command().Exec("invite", "user", "--email="+newEmail2, "--namespace="+ns)
	}, 7*time.Second); err != nil {
		t.Fatal(err)
		return
	}

	logout(serv, t)

	signup(serv, t, newEmail, password, signupOptions{inviterEmail: email, xthInvitee: 1, isInvitedToNamespace: true, shouldJoin: true, freeTier: true})
	if t.Failed() {
		return
	}
	t.Logf("Signup 2 complete %s", time.Now())
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

	logout(serv, t)

	signup(serv, t, newEmail2, password, signupOptions{inviterEmail: email, xthInvitee: 2, isInvitedToNamespace: true, shouldJoin: true, freeTier: true})
	t.Logf("Signup 3 complete %s", time.Now())
	if t.Failed() {
		return
	}
	outp, err = serv.Command().Exec("user", "config", "get", "namespaces."+serv.Env()+".current")
	if err != nil {
		t.Fatalf("Error getting namespace: %v", err)
		return
	}
	newNs = strings.TrimSpace(string(outp))
	if newNs != ns {
		t.Fatalf("Namespaces should match, old: %v, new: %v", ns, newNs)
		return
	}

	t.T().Logf("Namespace joined: %v", string(outp))
}
