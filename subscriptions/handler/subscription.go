package handler

import (
	"context"
	"encoding/json"
	"time"

	mconfig "github.com/micro/micro/v3/service/config"

	merrors "github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	mstore "github.com/micro/micro/v3/service/store"

	"github.com/micro/go-micro/v3/errors"

	paymentsproto "github.com/m3o/services/payments/provider/proto"
	subscription "github.com/m3o/services/subscriptions/proto"
	"github.com/micro/go-micro/v3/client"
)

const (
	subscriptionTopic = "subscriptions"

	prefixSubscription = "subscription/"
	prefixCustomer     = "customer/"
)

var (
	additionalUsersPriceID = ""
	planID                 = ""
)

type Subscriptions struct {
	paymentService paymentsproto.ProviderService
}

type SubscriptionType struct {
	PlanID  string
	PriceID string
}

func New(paySvc paymentsproto.ProviderService) *Subscriptions {
	additionalUsersPriceID = mconfig.Get("micro", "signup", "additional_users_price_id").String("")
	planID = mconfig.Get("micro", "signup", "plan_id").String("")
	if len(planID) == 0 {
		logger.Error("No stripe plan id")
	}
	if len(additionalUsersPriceID) == 0 {
		logger.Error("No addition user plan id")
	}

	return &Subscriptions{
		paymentService: paySvc,
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

	rsp, err := s.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
		CustomerId:   email,
		CustomerType: "user",
		PlanId:       planID,
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

func (s Subscriptions) AddUser(ctx context.Context, request *subscription.AddUserRequest, response *subscription.AddUserResponse) error {
	subs, err := s.paymentService.ListSubscriptions(ctx, &paymentsproto.ListSubscriptionsRequest{
		CustomerId:   request.OwnerID,
		CustomerType: "user",
		PriceId:      additionalUsersPriceID,
	}, client.WithAuthToken())
	if err != nil {
		return merrors.InternalServerError("subscriptions.adduser.read", "Error finding sub: %v", err)
	}
	var sub *paymentsproto.Subscription
	if len(subs.Subscriptions) > 0 {
		sub = subs.Subscriptions[0]
	}

	if sub == nil {
		logger.Info("Creating sub with quantity 1")
		_, err = s.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
			CustomerId:   request.OwnerID,
			CustomerType: "user",
			PriceId:      additionalUsersPriceID,
			Quantity:     1,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	} else {
		logger.Info("Increasing sub quantity")
		_, err = s.paymentService.UpdateSubscription(ctx, &paymentsproto.UpdateSubscriptionRequest{
			SubscriptionId: sub.Id,
			CustomerId:     request.OwnerID,
			CustomerType:   "user",
			PriceId:        additionalUsersPriceID,
			Quantity:       sub.Quantity + 1,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	}
	if err != nil {
		return merrors.InternalServerError("signup", "Error increasing additional user quantity: %v", err)
	}
	return nil
}
