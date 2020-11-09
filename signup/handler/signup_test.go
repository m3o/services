package handler

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	mt "github.com/m3o/services/internal/test"

	malert "github.com/m3o/services/alert/proto/alert/alertfakes"
	mcustpb "github.com/m3o/services/customers/proto"
	mcust "github.com/m3o/services/customers/proto/protofakes"
	memail "github.com/m3o/services/emails/proto/protofakes"
	minvitepb "github.com/m3o/services/invite/proto"
	minvite "github.com/m3o/services/invite/proto/protofakes"
	mnspb "github.com/m3o/services/namespaces/proto"
	mns "github.com/m3o/services/namespaces/proto/protofakes"
	mpay "github.com/m3o/services/payments/proto/protofakes"
	pb "github.com/m3o/services/signup/proto/signup"
	msub "github.com/m3o/services/subscriptions/proto/protofakes"
	mauth "github.com/micro/micro/v3/service/auth/noop"
	"github.com/micro/micro/v3/service/errors"
	mstore "github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"

	. "github.com/onsi/gomega"

	"github.com/patrickmn/go-cache"
)

func TestMain(m *testing.M) {
	mstore.DefaultStore = memory.NewStore()
	m.Run()
}

func mockedSignup() *Signup {
	return &Signup{
		inviteService:       &minvite.FakeInviteService{},
		customerService:     &mcust.FakeCustomersService{},
		namespaceService:    &mns.FakeNamespacesService{},
		subscriptionService: &msub.FakeSubscriptionsService{},
		paymentService:      &mpay.FakeProviderService{},
		emailService:        &memail.FakeEmailsService{},
		auth:                mauth.NewAuth(),
		config:              conf{},
		cache:               cache.New(1*time.Minute, 5*time.Minute),
		alertService:        &malert.FakeAlertService{},
	}

}

func TestSendVerificationEmail(t *testing.T) {
	g := NewWithT(t)
	signupSvc := mockedSignup()

	invite := minvite.FakeInviteService{}
	invite.ValidateReturns(nil, errors.InternalServerError("error", "Error"))
	// set up so second call is successful
	invite.ValidateReturnsOnCall(1, &minvitepb.ValidateResponse{}, nil)

	signupSvc.inviteService = &invite

	cust := mcust.FakeCustomersService{}
	cust.CreateReturns(&mcustpb.CreateResponse{Customer: &mcustpb.Customer{Id: "1234"}}, nil)
	signupSvc.customerService = &cust

	// should have an error
	err := signupSvc.sendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{Email: "foo@bar1.com"}, &pb.SendVerificationEmailResponse{})
	g.Expect(err).To(Equal(errors.Forbidden("signup.notallowed", notInvitedErrorMsg)))

	// should succeed
	err = signupSvc.sendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{Email: "foo@bar.com"}, &pb.SendVerificationEmailResponse{})
	g.Expect(err).To(BeNil())

}

func TestSignup(t *testing.T) {
	t.Run("PaidTier", func(t *testing.T) {
		signupSvc := mockedSignup()
		signupSvc.config = conf{
			NoPayment: false,
		}
		invite := minvite.FakeInviteService{}
		invite.ValidateReturns(&minvitepb.ValidateResponse{}, nil)
		signupSvc.inviteService = &invite
		testSignup(t, signupSvc, false, func(g *WithT, verRsp *pb.VerifyResponse) {
			g.Expect(verRsp.Namespaces).To(BeEmpty())
			g.Expect(verRsp.PaymentRequired).To(BeTrue())
			g.Expect(verRsp.CustomerID).To(Equal("1234"))
			g.Expect(verRsp.Message).ToNot(BeEmpty())
		})
	})
	t.Run("FreeTier", func(t *testing.T) {
		signupSvc := mockedSignup()
		signupSvc.config = conf{
			NoPayment: true,
		}
		invite := minvite.FakeInviteService{}
		invite.ValidateReturns(&minvitepb.ValidateResponse{}, nil)
		signupSvc.inviteService = &invite
		testSignup(t, signupSvc, false, func(g *WithT, verRsp *pb.VerifyResponse) {
			g.Expect(verRsp.Namespaces).To(BeEmpty())
			g.Expect(verRsp.PaymentRequired).To(BeFalse())
			g.Expect(verRsp.CustomerID).To(Equal("1234"))
			g.Expect(verRsp.Message).ToNot(BeEmpty())
		})
	})
	t.Run("PaidTierJoinNamespace", func(t *testing.T) {
		signupSvc := mockedSignup()
		signupSvc.config = conf{
			NoPayment: false,
		}
		invite := minvite.FakeInviteService{}
		invite.ValidateReturns(&minvitepb.ValidateResponse{
			Namespaces: []string{"foo"},
		}, nil)
		signupSvc.inviteService = &invite

		testSignup(t, signupSvc, true, func(g *WithT, verRsp *pb.VerifyResponse) {
			g.Expect(verRsp.Namespaces).ToNot(BeEmpty())
			g.Expect(verRsp.Namespaces[0]).To(Equal("foo"))
			g.Expect(verRsp.PaymentRequired).To(BeTrue())
			g.Expect(verRsp.CustomerID).To(Equal("1234"))
			g.Expect(verRsp.Message).ToNot(BeEmpty())
		})
	})
	t.Run("FreeTierJoinNamespace", func(t *testing.T) {
		signupSvc := mockedSignup()
		signupSvc.config = conf{
			NoPayment: true,
		}
		invite := minvite.FakeInviteService{}
		invite.ValidateReturns(&minvitepb.ValidateResponse{
			Namespaces: []string{"foo"},
		}, nil)
		signupSvc.inviteService = &invite
		testSignup(t, signupSvc, true, func(g *WithT, verRsp *pb.VerifyResponse) {
			g.Expect(verRsp.Namespaces).ToNot(BeEmpty())
			g.Expect(verRsp.Namespaces[0]).To(Equal("foo"))
			g.Expect(verRsp.PaymentRequired).To(BeFalse())
			g.Expect(verRsp.CustomerID).To(Equal("1234"))
			g.Expect(verRsp.Message).ToNot(BeEmpty())
		})
	})
	t.Run("PaidTierInvitedButDeclineJoinNamespace", func(t *testing.T) {
		signupSvc := mockedSignup()
		signupSvc.config = conf{
			NoPayment: false,
		}
		invite := minvite.FakeInviteService{}
		invite.ValidateReturns(&minvitepb.ValidateResponse{
			Namespaces: []string{"foo"},
		}, nil)
		signupSvc.inviteService = &invite
		testSignup(t, signupSvc, false, func(g *WithT, verRsp *pb.VerifyResponse) {
			g.Expect(verRsp.Namespaces).ToNot(BeEmpty())
			g.Expect(verRsp.Namespaces[0]).To(Equal("foo"))
			g.Expect(verRsp.PaymentRequired).To(BeTrue())
			g.Expect(verRsp.CustomerID).To(Equal("1234"))
			g.Expect(verRsp.Message).ToNot(BeEmpty())
		})
	})
	t.Run("FreeTierInvitedButDeclineJoinNamespace", func(t *testing.T) {
		signupSvc := mockedSignup()
		signupSvc.config = conf{
			NoPayment: true,
		}
		invite := minvite.FakeInviteService{}
		invite.ValidateReturns(&minvitepb.ValidateResponse{
			Namespaces: []string{"foo"},
		}, nil)
		signupSvc.inviteService = &invite
		testSignup(t, signupSvc, false, func(g *WithT, verRsp *pb.VerifyResponse) {
			g.Expect(verRsp.Namespaces).ToNot(BeEmpty())
			g.Expect(verRsp.Namespaces[0]).To(Equal("foo"))
			g.Expect(verRsp.PaymentRequired).To(BeFalse())
			g.Expect(verRsp.CustomerID).To(Equal("1234"))
			g.Expect(verRsp.Message).ToNot(BeEmpty())
		})
	})

}

func testSignup(t *testing.T, signupSvc *Signup, join bool, checkVerifyRsp func(g *WithT, response *pb.VerifyResponse)) {
	g := NewWithT(t)

	emails := memail.FakeEmailsService{}
	signupSvc.emailService = &emails

	cust := mcust.FakeCustomersService{}
	cust.CreateReturns(&mcustpb.CreateResponse{Customer: &mcustpb.Customer{Id: "1234"}}, nil)
	signupSvc.customerService = &cust

	nsSvc := mns.FakeNamespacesService{}
	nsSvc.ReadReturns(&mnspb.ReadResponse{
		Namespace: &mnspb.Namespace{
			Id:     "foo-bar-baz",
			Owners: []string{"6789"},
		},
	}, nil)
	nsSvc.CreateReturns(&mnspb.CreateResponse{
		Namespace: &mnspb.Namespace{
			Id:     "foo-bar-baz",
			Owners: []string{"6789"},
		},
	}, nil)
	signupSvc.namespaceService = &nsSvc

	email := mt.TestEmail()

	// send verification email
	err := signupSvc.sendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{Email: email}, &pb.SendVerificationEmailResponse{})
	g.Expect(err).To(BeNil())

	// grab token
	_, req, _ := emails.SendArgsForCall(0)
	dat := map[string]interface{}{}
	g.Expect(json.Unmarshal(req.TemplateData, &dat)).To(BeNil())
	tok, ok := dat["token"].(string)
	g.Expect(ok).To(BeTrue())
	g.Expect(tok).To(HaveLen(8))

	// test incorrect tok
	err = signupSvc.Verify(context.TODO(), &pb.VerifyRequest{
		Email: email,
		Token: "aslkdja",
	}, &pb.VerifyResponse{})
	g.Expect(err).To(HaveOccurred())

	// test correct tok
	verRsp := &pb.VerifyResponse{}
	err = signupSvc.Verify(context.TODO(), &pb.VerifyRequest{
		Email: email,
		Token: tok,
	}, verRsp)
	g.Expect(err).To(BeNil())

	checkVerifyRsp(g, verRsp)

	g.Expect(cust.MarkVerifiedCallCount()).To(Equal(1))

	if !signupSvc.config.NoPayment {
		err := signupSvc.SetPaymentMethod(context.TODO(), &pb.SetPaymentMethodRequest{
			Email:         email,
			PaymentMethod: "pm_12345",
		}, &pb.SetPaymentMethodResponse{})
		g.Expect(err).To(BeNil())
	}

	cmpRsp := &pb.CompleteSignupResponse{}
	ns := ""
	if join {
		ns = verRsp.Namespaces[0]
	}
	err = signupSvc.CompleteSignup(context.TODO(), &pb.CompleteSignupRequest{
		Email:     email,
		Token:     tok,
		Secret:    "password",
		Namespace: ns,
	}, cmpRsp)
	g.Expect(err).To(BeNil())
	if join {
		g.Expect(cmpRsp.Namespace).To(Equal(ns))
	} else {
		g.Expect(cmpRsp.Namespace).NotTo(Equal(ns))
		g.Expect(strings.Count(cmpRsp.Namespace, "-")).To(Equal(2))
	}

	// TODO check sub service addUser is called

}
