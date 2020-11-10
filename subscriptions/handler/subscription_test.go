package handler

import (
	"encoding/json"
	"testing"
	"time"

	mt "github.com/m3o/services/internal/test"
	"github.com/m3o/services/internal/test/fakes"
	mprovpb "github.com/m3o/services/payments/proto"
	mprov "github.com/m3o/services/payments/proto/protofakes"
	pb "github.com/m3o/services/subscriptions/proto"
	mevents "github.com/micro/micro/v3/service/events"
	mstore "github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"

	. "github.com/onsi/gomega"
)

func mockedSubscription() *Subscriptions {
	ppsvc := &mprov.FakeProviderService{}
	ppsvc.CreateSubscriptionReturns(&mprovpb.CreateSubscriptionResponse{
		Subscription: &mprovpb.Subscription{Id: "5678"},
	}, nil)
	return &Subscriptions{
		config: config{
			AdditionalUsersPriceID: "aupid",
			PlanID:                 "pid",
		},
		paymentService: ppsvc,
	}
}

func TestMain(m *testing.M) {
	mevents.DefaultStream = &fakes.FakeStream{}
	mstore.DefaultStore = memory.NewStore()
	m.Run()
}

func TestSubCreateAndCancel(t *testing.T) {
	g := NewWithT(t)
	subSvc := mockedSubscription()
	ppsvc := subSvc.paymentService.(*mprov.FakeProviderService)
	fstream := mevents.DefaultStream.(*fakes.FakeStream)
	adminCtx := mt.ContextWithAccount("micro", "foo")

	// creation
	cRsp := &pb.CreateResponse{}
	err := subSvc.Create(adminCtx, &pb.CreateRequest{
		CustomerID:      "1234",
		Type:            "user",
		PaymentMethodID: "pm_1234",
		Email:           "foo@bar.com",
	}, cRsp)
	g.Expect(err).To(BeNil())

	g.Expect(fstream.PublishCallCount()).To(Equal(1))
	g.Expect(ppsvc.CreateCustomerCallCount()).To(Equal(1))
	g.Expect(ppsvc.CreatePaymentMethodCallCount()).To(Equal(1))
	g.Expect(ppsvc.SetDefaultPaymentMethodCallCount()).To(Equal(1))
	g.Expect(ppsvc.CreateSubscriptionCallCount()).To(Equal(1))

	recs, err := mstore.Read("", mstore.Prefix(prefixCustomer+"1234/"))
	g.Expect(err).To(BeNil())
	g.Expect(recs).To(HaveLen(1))

	// add user
	ppsvc.ListSubscriptionsReturnsOnCall(0, &mprovpb.ListSubscriptionsResponse{}, nil)
	err = subSvc.AddUser(adminCtx, &pb.AddUserRequest{
		OwnerID:   "1234",
		NewUserID: "2345",
	}, &pb.AddUserResponse{})
	g.Expect(err).To(BeNil())
	// adding the first user should create a new subscription
	g.Expect(ppsvc.CreateSubscriptionCallCount()).To(Equal(2))
	// check its created with the right arg
	_, req, _ := ppsvc.CreateSubscriptionArgsForCall(1)
	g.Expect(req.PriceId).To(Equal(subSvc.config.AdditionalUsersPriceID))
	g.Expect(fstream.PublishCallCount()).To(Equal(2))

	// add a second user
	ppsvc.ListSubscriptionsReturnsOnCall(1, &mprovpb.ListSubscriptionsResponse{
		Subscriptions: []*mprovpb.Subscription{&mprovpb.Subscription{Id: "sub_1234"}},
	}, nil)
	err = subSvc.AddUser(adminCtx, &pb.AddUserRequest{
		OwnerID:   "1234",
		NewUserID: "3456",
	}, &pb.AddUserResponse{})
	g.Expect(err).To(BeNil())
	// adding the second user should update the subscription
	g.Expect(ppsvc.CreateSubscriptionCallCount()).To(Equal(2))
	g.Expect(ppsvc.UpdateSubscriptionCallCount()).To(Equal(1))
	g.Expect(fstream.PublishCallCount()).To(Equal(3))

	// cancellation
	err = subSvc.Cancel(adminCtx, &pb.CancelRequest{
		CustomerID: "1234",
	}, &pb.CancelResponse{})
	g.Expect(err).To(BeNil())

	recs, err = mstore.Read("", mstore.Prefix(prefixCustomer+"1234/"))
	g.Expect(err).To(BeNil())
	g.Expect(recs).To(HaveLen(1))
	sub := &Subscription{}
	g.Expect(json.Unmarshal(recs[0].Value, sub)).To(BeNil())
	g.Expect(sub.Expires).To(BeNumerically("<", time.Now().Unix()+1))
	g.Expect(ppsvc.DeleteCustomerCallCount()).To(Equal(1))
	// three more publishes; 1 for main sub cancel + 2 more for the child subs
	g.Expect(fstream.PublishCallCount()).To(Equal(6))

}
