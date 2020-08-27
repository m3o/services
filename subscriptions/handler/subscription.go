package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	paymentsproto "github.com/m3o/services/payments/provider/proto"
	subscription "github.com/m3o/services/subscriptions/proto"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	merrors "github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/store"
	mconfig "github.com/micro/micro/v3/service/config"
	mcontext "github.com/micro/micro/v3/service/context"
	eventsproto "github.com/micro/micro/v3/service/events/proto"
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
	streamService  eventsproto.StreamService
}

type SubscriptionType struct {
	PlanID  string
	PriceID string
}

func New(paySvc paymentsproto.ProviderService, streamService eventsproto.StreamService) *Subscriptions {
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
		streamService:  streamService,
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
	sub := &Subscription{
		Type:                  request.Type,
		CustomerID:            email,
		Created:               time.Now().Unix(),
		ID:                    uuid.New().String(),
		PaymentSubscriptionID: rsp.Subscription.Id,
	}
	if err := s.writeSubscription(sub); err != nil {
		return err
	}
	response.Subscription = objToProto(sub)
	return s.eventPublish(subscriptionTopic, SubscriptionEvent{Subscription: *sub, Type: "subscriptions.created"})
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

	recs, err := mstore.Read(prefixCustomer+request.OwnerID+"/", store.ReadPrefix())
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
		CustomerID:            request.NewUserID,
		Created:               time.Now().Unix(),
		ID:                    uuid.New().String(),
		ParentSubscriptionID:  parentSub.ID,
		PaymentSubscriptionID: parentSub.PaymentSubscriptionID,
	}

	if err := s.writeSubscription(subscription); err != nil {
		return err
	}

	return s.eventPublish(subscriptionTopic,
		SubscriptionEvent{Subscription: *subscription, Type: "subscriptions.created"},
		events.WithMetadata(map[string]string{"user": request.NewUserID}),
	)

}

// TODO remove this and replace with publish from micro/micro
func (s Subscriptions) eventPublish(topic string, msg interface{}, opts ...events.PublishOption) error {
	// parse the options
	options := events.PublishOptions{
		Timestamp: time.Now(),
	}
	for _, o := range opts {
		o(&options)
	}

	// encode the message if it's not already encoded
	var payload []byte
	if p, ok := msg.([]byte); ok {
		payload = p
	} else {
		p, err := json.Marshal(msg)
		if err != nil {
			return events.ErrEncodingMessage
		}
		payload = p
	}

	// execute the RPC
	_, err := s.streamService.Publish(mcontext.DefaultContext, &eventsproto.PublishRequest{
		Topic:     topic,
		Payload:   payload,
		Metadata:  options.Metadata,
		Timestamp: options.Timestamp.Unix(),
	}, client.WithAuthToken())

	return err
}
