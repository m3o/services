package handler

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/m3o/services/tests/fakes"
	mevents "github.com/micro/micro/v3/service/events"
	mstore "github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"

	mt "github.com/m3o/services/internal/test"
	mprovpb "github.com/m3o/services/payments/proto"
	mprov "github.com/m3o/services/payments/proto/protofakes"
	pb "github.com/m3o/services/subscriptions/proto"

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

func TestSubCancellation(t *testing.T) {
	g := NewWithT(t)
	subSvc := mockedSubscription()
	fstream := &fakes.FakeStream{}
	mevents.DefaultStream = fstream
	mstore.DefaultStore = memory.NewStore()
	adminCtx := mt.ContextWithAccount("micro", "foo")
	cRsp := &pb.CreateResponse{}
	err := subSvc.Create(adminCtx, &pb.CreateRequest{
		CustomerID:      "1234",
		Type:            "user",
		PaymentMethodID: "pm_1234",
		Email:           "foo@bar.com",
	}, cRsp)
	g.Expect(err).To(BeNil())
	g.Expect(fstream.PublishCallCount()).To(Equal(1))

	recs, err := mstore.Read("", mstore.Prefix(prefixCustomer+"1234/"))
	g.Expect(err).To(BeNil())
	g.Expect(recs).To(HaveLen(1))

	err = subSvc.Cancel(adminCtx, &pb.CancelRequest{
		CustomerID: "1234",
	}, &pb.CancelResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(fstream.PublishCallCount()).To(Equal(2))

	recs, err = mstore.Read("", mstore.Prefix(prefixCustomer+"1234/"))
	g.Expect(err).To(BeNil())
	g.Expect(recs).To(HaveLen(1))
	sub := &Subscription{}
	g.Expect(json.Unmarshal(recs[0].Value, sub)).To(BeNil())
	g.Expect(sub.Expires).To(BeNumerically("<", time.Now().Unix()+1))
	ppsvc := subSvc.paymentService.(*mprov.FakeProviderService)
	g.Expect(ppsvc.DeleteCustomerCallCount()).To(Equal(1))
}
