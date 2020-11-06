package handler

import (
	"context"
	"testing"
	"time"

	minvitepb "github.com/m3o/services/invite/proto"

	"github.com/stretchr/testify/mock"

	malert "github.com/m3o/services/alert/mocks"
	mcust "github.com/m3o/services/customers/mocks"
	mcustpb "github.com/m3o/services/customers/proto"
	memail "github.com/m3o/services/emails/mocks"
	minvite "github.com/m3o/services/invite/mocks"
	mns "github.com/m3o/services/namespaces/mocks"
	mpay "github.com/m3o/services/payments/mocks"
	pb "github.com/m3o/services/signup/proto/signup"
	msub "github.com/m3o/services/subscriptions/mocks"
	mauth "github.com/micro/micro/v3/service/auth/client"
	"github.com/micro/micro/v3/service/errors"

	"github.com/stretchr/testify/assert"

	"github.com/patrickmn/go-cache"
)

func setup() {
	// mock things like store, config, etc so that you can do expect
}

func mockedSignup() *Signup {
	return &Signup{
		inviteService:       &minvite.InviteService{},
		customerService:     &mcust.CustomersService{},
		namespaceService:    &mns.NamespacesService{},
		subscriptionService: &msub.SubscriptionsService{},
		paymentService:      &mpay.ProviderService{},
		emailService:        &memail.EmailsService{},
		auth:                mauth.NewAuth(),
		config:              conf{},
		cache:               cache.New(1*time.Minute, 5*time.Minute),
		alertService:        &malert.AlertService{},
	}

}

func TestSendVerificationEmail(t *testing.T) {
	signupSvc := mockedSignup()

	invite := minvite.InviteService{}
	invite.On("Validate", mock.Anything, &minvitepb.ValidateRequest{
		Email: "foo@bar.com",
	}, mock.Anything).Return(&minvitepb.ValidateResponse{}, nil)
	invite.On("Validate", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.InternalServerError("error", "Error"))
	signupSvc.inviteService = &invite

	cust := mcust.CustomersService{}
	cust.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(&mcustpb.CreateResponse{Customer: &mcustpb.Customer{Id: "1234"}}, nil)
	signupSvc.customerService = &cust

	err := signupSvc.sendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{Email: "foo@bar1.com"}, &pb.SendVerificationEmailResponse{})
	if assert.Error(t, err) {
		assert.Equal(t, errors.Forbidden("signup.notallowed", notInvitedErrorMsg), err)
	}

	err = signupSvc.sendVerificationEmail(context.TODO(), &pb.SendVerificationEmailRequest{Email: "foo@bar.com"}, &pb.SendVerificationEmailResponse{})
	if assert.Error(t, err) {
		assert.Equal(t, errors.Forbidden("signup.notallowed", notInvitedErrorMsg), err)
	}
}
