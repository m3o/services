package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/events"
	mstore "github.com/micro/micro/v3/service/store"

	"github.com/micro/go-micro/v3/errors"

	paymentsproto "github.com/m3o/services/payments/provider/proto"
	subscription "github.com/m3o/services/subscriptions/proto"
	"github.com/micro/go-micro/v3/client"
	log "github.com/micro/go-micro/v3/logger"
	mconfig "github.com/micro/micro/v3/service/config"
)

const (
	subscriptionTopic = "subscriptions"

	prefixSubscription = "subscription/"
	prefixCustomer     = "customer/"
)

type Subscriptions struct {
	paymentService paymentsproto.ProviderService
	typeMap        TypeMap
}

type SubscriptionType struct {
	PlanID string
}

type TypeMap map[string]SubscriptionType

func New(paySvc paymentsproto.ProviderService) *Subscriptions {
	typeMap := TypeMap{}
	if err := mconfig.Get("micro", "subscriptions", "types").Scan(&typeMap); err != nil {
		log.Fatalf("Failed to retrieve subscription type map %s", err)
	}

	return &Subscriptions{
		paymentService: paySvc,
		typeMap:        typeMap,
	}
}

type SubscriptionModel struct {
	ID         string
	CustomerID string
	Type       string
	Created    int64
	Expires    int64
}

func objToProto(sub *SubscriptionModel) *subscription.Subscription {
	return &subscription.Subscription{
		CustomerID: sub.CustomerID,
		Created:    sub.Created,
		Expires:    sub.Expires,
		Id:         sub.ID,
		Type:       sub.Type,
	}
}

func (s Subscriptions) Create(ctx context.Context, request *subscription.CreateRequest, response *subscription.CreateResponse) error {
	email := request.CustomerID
	_, err := s.paymentService.CreateCustomer(ctx, &paymentsproto.CreateCustomerRequest{
		Customer: &paymentsproto.Customer{
			Id:   email,
			Type: "user",
			Metadata: map[string]string{
				"email": email,
			},
		},
	}, client.WithAuthToken())
	if err != nil {
		return err
	}
	// TODO The above call might take a while to complete
	_, err = s.paymentService.CreatePaymentMethod(ctx, &paymentsproto.CreatePaymentMethodRequest{
		CustomerId:   email,
		CustomerType: "user",
		Id:           request.PaymentMethodID,
	}, client.WithAuthToken())
	if err != nil {
		return err
	}

	_, err = s.paymentService.SetDefaultPaymentMethod(ctx, &paymentsproto.SetDefaultPaymentMethodRequest{
		CustomerId:      email,
		CustomerType:    "user",
		PaymentMethodId: request.PaymentMethodID,
	}, client.WithAuthToken())
	if err != nil {
		return err
	}

	subscriptionType, ok := s.typeMap[request.Type]
	if !ok {
		return errors.BadRequest("subscriptions.create.subtype", "Subscription type not recognised")
	}
	rsp, err := s.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
		CustomerId:   email,
		CustomerType: "user",
		PlanId:       subscriptionType.PlanID,
	}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	if err != nil {
		return err
	}
	sub := &SubscriptionModel{
		Type:       request.Type,
		CustomerID: email,
		Created:    time.Now().Unix(),
		ID:         rsp.Subscription.Id,
	}
	b, err := json.Marshal(sub)
	if err != nil {
		return err
	}
	if err := mstore.Write(&store.Record{
		Key:   prefixSubscription + sub.ID,
		Value: b,
	}); err != nil {
		return err
	}
	if err := mstore.Write(&store.Record{
		Key:   prefixCustomer + sub.CustomerID + "-" + sub.ID,
		Value: b,
	}); err != nil {
		return err
	}
	response.Subscription = objToProto(sub)
	return events.Publish(subscriptionTopic, SubscriptionEvent{Subscription: *sub, Type: "subscriptions.created"})
}

func (s Subscriptions) Cancel(ctx context.Context, request *subscription.CancelRequest, response *subscription.CancelResponse) error {
	return errors.InternalServerError("notimplemented", "not implemented")
}

func (s Subscriptions) List(ctx context.Context, request *subscription.ListRequest, response *subscription.ListResponse) error {
	return errors.InternalServerError("notimplemented", "not implemented")
}
