package handler

import (
	"context"
	"testing"
	"time"

	malert "github.com/m3o/services/alert/proto/alert/alertfakes"
	mcustpb "github.com/m3o/services/customers/proto"
	mcust "github.com/m3o/services/customers/proto/protofakes"
	memail "github.com/m3o/services/emails/proto/protofakes"
	minvitepb "github.com/m3o/services/invite/proto"
	minvite "github.com/m3o/services/invite/proto/protofakes"
	mns "github.com/m3o/services/namespaces/proto/protofakes"
	mpay "github.com/m3o/services/payments/proto/protofakes"
	pb "github.com/m3o/services/signup/proto/signup"
	fakes "github.com/m3o/services/signup/signupfakes"
	msub "github.com/m3o/services/subscriptions/proto/protofakes"
	mauth "github.com/micro/micro/v3/service/auth/client"
	"github.com/micro/micro/v3/service/errors"
	mstore "github.com/micro/micro/v3/service/store"

	. "github.com/onsi/gomega"

	"github.com/patrickmn/go-cache"
)

func setup() {
	// mock things like store, config, etc so that you can do expect
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
	g := NewGomegaWithT(t)
	signupSvc := mockedSignup()

	mstore.DefaultStore = &fakes.FakeStore{}
	invite := minvite.FakeInviteService{}
	invite.ValidateReturns(nil, errors.InternalServerError("error", "Error"))
	invite.ValidateReturnsOnCall(1, &minvitepb.ValidateResponse{}, nil)

	signupSvc.inviteService = &invite

	cust := mcust.FakeCustomersService{}
	cust.CreateReturns(&mcustpb.CreateResponse{Customer: &mcustpb.Customer{Id: "1234"}}, nil)
	signupSvc.customerService = &cust

	err := signupSvc.sendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{Email: "foo@bar1.com"}, &pb.SendVerificationEmailResponse{})
	g.Expect(err).To(Equal(errors.Forbidden("signup.notallowed", notInvitedErrorMsg)))

	err = signupSvc.sendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{Email: "foo@bar.com"}, &pb.SendVerificationEmailResponse{})
	g.Expect(err).To(BeNil())
}
