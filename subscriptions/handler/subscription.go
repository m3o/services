package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/micro/go-micro/v3/auth"

	"github.com/google/uuid"

	paymentsproto "github.com/m3o/services/payments/provider/proto"
	subscription "github.com/m3o/services/subscriptions/proto"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	merrors "github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/store"
	mconfig "github.com/micro/micro/v3/service/config"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	mstore "github.com/micro/micro/v3/service/store"
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
	planID         string
}

type SubscriptionType struct {
	PlanID  string
	PriceID string
}

func New(paySvc paymentsproto.ProviderService) *Subscriptions {
	additionalUsersPriceID = mconfig.Get("micro", "subscriptions", "additional_users_price_id").String("")
	planID = mconfig.Get("micro", "subscriptions", "plan_id").String("")
	if len(planID) == 0 {
		logger.Error("No stripe plan id")
	}
	if len(additionalUsersPriceID) == 0 {
		logger.Error("No addition user plan id")
	}

	return &Subscriptions{
		paymentService: paySvc,
		planID:         planID,
	}
}

type Subscription struct {
	// ID in this service
	ID string
	// ID in the payment service
	PaymentSubscriptionID string
	CustomerID            string
	// developer, additional
	Type    string
	Created int64
	// If this sub has been cancelled this represents the end date
	Expires int64
	// If this subscription is paid for by another subscription, this is populated with the ID of the paying sybscription
	ParentSubscriptionID string
}

func objToProto(sub *Subscription) *subscription.Subscription {
	return &subscription.Subscription{
		CustomerID1: sub.CustomerID,
		Created:     sub.Created,
		Expires:     sub.Expires,
		Id:          sub.ID,
		Type:        sub.Type,
	}
}

func (s Subscriptions) Create(ctx context.Context, request *subscription.CreateRequest, response *subscription.CreateResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	customerID := request.CustomerID1
	_, err := s.paymentService.CreateCustomer(ctx, &paymentsproto.CreateCustomerRequest{
		Customer: &paymentsproto.Customer{
			Id:   customerID,
			Type: "user",
			Metadata: map[string]string{
				"email": request.Email,
			},
		},
	}, client.WithAuthToken())
	if err != nil {
		return err
	}
	// TODO The above call might take a while to complete
	_, err = s.paymentService.CreatePaymentMethod(ctx, &paymentsproto.CreatePaymentMethodRequest{
		CustomerId:   customerID,
		CustomerType: "user",
		Id:           request.PaymentMethodID,
	}, client.WithAuthToken())
	if err != nil {
		return err
	}

	_, err = s.paymentService.SetDefaultPaymentMethod(ctx, &paymentsproto.SetDefaultPaymentMethodRequest{
		CustomerId:      customerID,
		CustomerType:    "user",
		PaymentMethodId: request.PaymentMethodID,
	}, client.WithAuthToken())
	if err != nil {
		return err
	}

	rsp, err := s.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
		CustomerId:   customerID,
		CustomerType: "user",
		PlanId:       planID,
		Quantity:     1,
	}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	if err != nil {
		return err
	}
	sub := &Subscription{
		Type:                  "developer", // TODO we'll end up supporting more that one sub type so we'll use request.Type,
		CustomerID:            customerID,
		Created:               time.Now().Unix(),
		ID:                    uuid.New().String(),
		PaymentSubscriptionID: rsp.Subscription.Id,
	}
	if err := s.writeSubscription(sub); err != nil {
		return err
	}
	response.Subscription = objToProto(sub)
	ev := SubscriptionEvent{Subscription: *sub, Type: "subscriptions.created"}
	if err := mevents.Publish(subscriptionTopic, ev); err != nil {
		logger.Errorf("Error publishing subscriptions.created for event %+v", ev)
	}
	return nil
}

func (s Subscriptions) writeSubscription(sub *Subscription) error {
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
		Key:   prefixCustomer + sub.CustomerID + "/" + sub.ID,
		Value: b,
	}); err != nil {
		return err
	}
	return nil
}

func (s Subscriptions) Cancel(ctx context.Context, request *subscription.CancelRequest, response *subscription.CancelResponse) error {
	return errors.InternalServerError("notimplemented", "not implemented")
}

func (s Subscriptions) AddUser(ctx context.Context, request *subscription.AddUserRequest, response *subscription.AddUserResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	subs, err := s.paymentService.ListSubscriptions(ctx, &paymentsproto.ListSubscriptionsRequest{
		CustomerId:   request.OwnerID1,
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
			CustomerId:   request.OwnerID1,
			CustomerType: "user",
			PriceId:      additionalUsersPriceID,
			Quantity:     1,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	} else {
		logger.Info("Increasing sub quantity")
		_, err = s.paymentService.UpdateSubscription(ctx, &paymentsproto.UpdateSubscriptionRequest{
			SubscriptionId: sub.Id,
			CustomerId:     request.OwnerID1,
			CustomerType:   "user",
			PriceId:        additionalUsersPriceID,
			Quantity:       sub.Quantity + 1,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	}
	if err != nil {
		return merrors.InternalServerError("signup", "Error increasing additional user quantity: %v", err)
	}

	recs, err := mstore.Read(prefixCustomer+request.OwnerID1+"/", store.ReadPrefix())
	if err != nil {
		return err
	}

	// we assume a user only has one subscription right now
	parentSub := &Subscription{}
	if err := json.Unmarshal(recs[0].Value, parentSub); err != nil {
		return err
	}

	subscription := &Subscription{
		Type:                  "additional",
		CustomerID:            request.NewUserID1,
		Created:               time.Now().Unix(),
		ID:                    uuid.New().String(),
		ParentSubscriptionID:  parentSub.ID,
		PaymentSubscriptionID: parentSub.PaymentSubscriptionID,
	}

	if err := s.writeSubscription(subscription); err != nil {
		return err
	}
	ev := SubscriptionEvent{Subscription: *subscription, Type: "subscriptions.created"}
	if err := mevents.Publish(subscriptionTopic, ev,
		events.WithMetadata(map[string]string{"user": request.NewUserID1}),
	); err != nil {
		logger.Errorf("Error publishing subscriptions.created for user %s event %+v", request.NewUserID1, ev)
	}
	return nil

}

func (s Subscriptions) Update(ctx context.Context, request *subscription.UpdateRequest, response *subscription.UpdateResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	subs, err := s.paymentService.ListSubscriptions(ctx, &paymentsproto.ListSubscriptionsRequest{
		CustomerId:   request.OwnerID1,
		CustomerType: "user",
		PriceId:      additionalUsersPriceID,
	}, client.WithAuthToken())
	if err != nil {
		return merrors.InternalServerError("subscriptions.update.read", "Error finding sub: %v", err)
	}
	var sub *paymentsproto.Subscription
	if len(subs.Subscriptions) > 0 {
		for _, su := range subs.Subscriptions {
			// plan and price ids are both store in s.Plan.Id
			if su.Plan.Id == request.PriceID {
				sub = su
				break
			}
		}
	}

	if sub == nil {
		if request.Quantity == 0 {
			return errors.InternalServerError("subscriptions.Update", "Something is wrong, trying to create subscription with 0 value")
		}
		logger.Infof("Creating sub with quantity %d", request.Quantity)
		_, err = s.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
			CustomerId:   request.OwnerID1,
			CustomerType: "user",
			PriceId:      request.PriceID,
			Quantity:     request.Quantity,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
		if err != nil {
			return merrors.InternalServerError("signup", "Error creating subscription: %v", err)
		}
	} else {
		logger.Info("Increasing subscription quantity")
		_, err = s.paymentService.UpdateSubscription(ctx, &paymentsproto.UpdateSubscriptionRequest{
			SubscriptionId: sub.Id,
			CustomerId:     request.OwnerID1,
			CustomerType:   "user",
			PriceId:        request.PriceID,
			Quantity:       request.Quantity,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
		if err != nil {
			return merrors.InternalServerError("signup", "Error updating subscription '%v': %v", sub.Id, err)
		}
	}

	return nil
}

func authorizeCall(ctx context.Context) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("subscriptions", "Unauthorized request")
	}
	if account.Issuer != "micro" {
		return errors.Unauthorized("subscriptions", "Unauthorized request")
	}
	return nil
}
