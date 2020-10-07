package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/micro/micro/v3/service/auth"

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

	prefixSubscription = "subscription/" // subscription/<subID>
	prefixCustomer     = "customer/"     // customer/<customerID>/<subID>
	prefixParentSub    = "parentSub/"    // parentSub/<parentSubID>/<childSubID>
)

type Subscriptions struct {
	config         config
	paymentService paymentsproto.ProviderService
}

type SubscriptionType struct {
	PlanID  string
	PriceID string
}

type config struct {
	AdditionalUsersPriceID string `json:"additional_users_price_id"`
	PlanID                 string `json:"plan_id"`
}

func New(paySvc paymentsproto.ProviderService) *Subscriptions {
	conf := config{}
	values, err := mconfig.Get("micro.subscriptions")
	if err != nil {
		logger.Warn(err)
	}
	err = values.Scan(&conf)
	if err != nil {
		logger.Warn(err)
	}

	if len(conf.PlanID) == 0 {
		logger.Error("No stripe plan id")
	}
	if len(conf.AdditionalUsersPriceID) == 0 {
		logger.Error("No addition user plan id")
	}

	return &Subscriptions{
		config:         conf,
		paymentService: paySvc,
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
	if err := authorizeAdminCall(ctx); err != nil {
		return err
	}
	customerID := request.CustomerID
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
		PlanId:       s.config.PlanID,
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
	if len(sub.ParentSubscriptionID) > 0 {
		if err := mstore.Write(&store.Record{
			Key:   prefixParentSub + sub.ParentSubscriptionID + "/" + sub.ID,
			Value: b,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s Subscriptions) Cancel(ctx context.Context, request *subscription.CancelRequest, response *subscription.CancelResponse) error {
	if err := authorizeCall(ctx, request.CustomerID); err != nil {
		return err
	}
	if len(request.CustomerID) == 0 {
		return errors.BadRequest("subscriptions.cancel.validation", "Customer ID is required")
	}
	// lookup the subscriptions for this customer
	// doing a prefix lookup so if request.SubscriptionID is blank we just look up all the customer's subs.
	// If they only have one then this will do the right thing
	recs, err := mstore.Read("", mstore.Prefix(prefixCustomer+request.CustomerID+"/"+request.SubscriptionID))
	if err != nil {
		return err
	}
	if len(recs) != 1 {
		return errors.BadRequest("subscriptions.cancel", "Found %d subscriptions for this user. Please specify a valid subscription ID", len(recs))
	}
	sub := &Subscription{}
	if err := json.Unmarshal(recs[0].Value, sub); err != nil {
		logger.Errorf("Error unmarshalling subscription %s %s", recs[0].Key, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support")
	}
	if len(sub.ParentSubscriptionID) > 0 {
		return s.cancelChildSubscription(ctx, sub)
	}

	return s.cancelSubscription(ctx, sub)
}

func (s Subscriptions) cancelSubscription(ctx context.Context, sub *Subscription) error {
	// clean up stripe.
	// deleting the customer will cancel all subscriptions, including the ones for additional services etc
	_, err := s.paymentService.DeleteCustomer(ctx, &paymentsproto.DeleteCustomerRequest{CustomerType: "user", CustomerId: sub.CustomerID})
	if ignoreDeleteError(err) != nil {
		logger.Errorf("Error cancelling subscription with stripe %s %s", sub.ID, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support.")
	}

	// update local obj
	sub.Expires = time.Now().Unix()
	if err := s.writeSubscription(sub); err != nil {
		logger.Errorf("Error persisting subscription cancellation %s %s", sub.ID, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support.")
	}

	// clean up any local subscription objects (additional users)
	recs, err := mstore.Read("", mstore.Prefix(prefixParentSub+sub.ID))
	if err != nil && err != mstore.ErrNotFound {
		logger.Errorf("Error looking up child subscriptions for customer %s subscription ID %s %s", sub.CustomerID, sub.ID, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support.")
	}
	for _, r := range recs {
		var sub *Subscription
		if err := json.Unmarshal(r.Value, sub); err != nil {
			logger.Errorf("Error unmarshalling subscription %s %s", r.Key, err)
			return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support")
		}
		sub.Expires = time.Now().Unix()
		if err := s.writeSubscription(sub); err != nil {
			logger.Errorf("Error updating subscription object for cancellation %s %s", sub.ID, err)
			return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support")
		}
		ev := SubscriptionEvent{Subscription: *sub, Type: "subscriptions.cancelled"}
		if err := mevents.Publish(subscriptionTopic, ev); err != nil {
			logger.Errorf("Error publishing subscriptions.cancelled for event %+v", ev)
		}
	}
	ev := SubscriptionEvent{Subscription: *sub, Type: "subscriptions.cancelled"}
	if err := mevents.Publish(subscriptionTopic, ev); err != nil {
		logger.Errorf("Error publishing subscriptions.cancelled for event %+v", ev)
	}

	return nil
}

func (s Subscriptions) cancelChildSubscription(ctx context.Context, sub *Subscription) error {
	// we should only decrement the additional user's subscription on the parent subscription, no other clean up required
	recs, err := mstore.Read(prefixSubscription + sub.ParentSubscriptionID)
	if err != nil {
		logger.Errorf("Error looking up parent subscription for cancellation parent %s, child %s, err %s", sub.ParentSubscriptionID, sub.ID, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support")
	}
	parentSub := &Subscription{}
	if err := json.Unmarshal(recs[0].Value, parentSub); err != nil {
		logger.Errorf("Error unmarshalling parent subscription for cancellation parent %s, child %s, err %s", sub.ParentSubscriptionID, sub.ID, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support")
	}

	if err := s.updatePaymentSubscription(ctx, parentSub.CustomerID, s.config.AdditionalUsersPriceID, -1, true); err != nil {
		logger.Errorf("Error updating subscription quantity from delete %s %s", sub.ID, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support")
	}
	sub.Expires = time.Now().Unix()
	if err := s.writeSubscription(sub); err != nil {
		logger.Errorf("Error updating subscription object for cancellation %s %s", sub.ID, err)
		return errors.InternalServerError("subscriptions.cancel", "Error cancelling subscription. Please contact support")
	}
	ev := SubscriptionEvent{Subscription: *sub, Type: "subscriptions.cancelled"}
	if err := mevents.Publish(subscriptionTopic, ev); err != nil {
		logger.Errorf("Error publishing subscriptions.cancelled for event %+v", ev)
	}

	return nil
}

// ignoreDeleteError will ignore any 400 or 404 errors returned, useful for idempotent deletes
func ignoreDeleteError(err error) error {
	if err != nil {
		merr, ok := err.(*errors.Error)
		if !ok {
			return err
		}
		if merr.Code == 400 || merr.Code == 404 {
			return nil
		}
		return err
	}
	return nil
}

func (s Subscriptions) AddUser(ctx context.Context, request *subscription.AddUserRequest, response *subscription.AddUserResponse) error {
	if err := authorizeAdminCall(ctx); err != nil {
		return err
	}
	subs, err := s.paymentService.ListSubscriptions(ctx, &paymentsproto.ListSubscriptionsRequest{
		CustomerId:   request.OwnerID,
		CustomerType: "user",
		PriceId:      s.config.AdditionalUsersPriceID,
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
			PriceId:      s.config.AdditionalUsersPriceID,
			Quantity:     1,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	} else {
		logger.Info("Increasing sub quantity")
		_, err = s.paymentService.UpdateSubscription(ctx, &paymentsproto.UpdateSubscriptionRequest{
			SubscriptionId: sub.Id,
			CustomerId:     request.OwnerID,
			CustomerType:   "user",
			PriceId:        s.config.AdditionalUsersPriceID,
			Quantity:       sub.Quantity + 1,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	}
	if err != nil {
		return merrors.InternalServerError("signup", "Error increasing additional user quantity: %v", err)
	}

	recs, err := mstore.Read("", mstore.Prefix(prefixCustomer+request.OwnerID+"/"))
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
	ev := SubscriptionEvent{Subscription: *subscription, Type: "subscriptions.created"}
	if err := mevents.Publish(subscriptionTopic, ev,
		events.WithMetadata(map[string]string{"user": request.NewUserID}),
	); err != nil {
		logger.Errorf("Error publishing subscriptions.created for user %s event %+v", request.NewUserID, ev)
	}
	return nil

}

func (s Subscriptions) Update(ctx context.Context, request *subscription.UpdateRequest, response *subscription.UpdateResponse) error {
	if err := authorizeAdminCall(ctx); err != nil {
		return err
	}
	if err := s.updatePaymentSubscription(ctx, request.OwnerID, request.PriceID, request.Quantity, false); err != nil {
		return err
	}
	return nil
}

// updatePaymentSubscription updates the given subscription with a new quantity. If the subscription doesn't yet exist it will create it.
// if qtyIsDelta quantity is treated as a delta and added to the existing quantity (pass a negative quantity to decrease the quantity).
func (s Subscriptions) updatePaymentSubscription(ctx context.Context, customerID, priceID string, quantity int64, qtyIsDelta bool) error {
	subs, err := s.paymentService.ListSubscriptions(ctx, &paymentsproto.ListSubscriptionsRequest{
		CustomerId:   customerID,
		CustomerType: "user",
		PriceId:      priceID,
	}, client.WithAuthToken())
	if err != nil {
		return merrors.NotFound("subscriptions.update.read", "Error finding sub: %v", err)
	}
	var sub *paymentsproto.Subscription
	if len(subs.Subscriptions) > 0 {
		for _, su := range subs.Subscriptions {
			// plan and price ids are both store in s.Plan.Id
			if su.Plan.Id == priceID {
				sub = su
				break
			}
		}
	}

	if sub == nil {
		if quantity == 0 {
			return errors.InternalServerError("subscriptions.Update", "Something is wrong, trying to create subscription with 0 value")
		}
		logger.Infof("Creating sub with quantity %d", quantity)
		_, err = s.paymentService.CreateSubscription(ctx, &paymentsproto.CreateSubscriptionRequest{
			CustomerId:   customerID,
			CustomerType: "user",
			PriceId:      priceID,
			Quantity:     quantity,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
		if err != nil {
			return merrors.InternalServerError("signup", "Error creating subscription: %v", err)
		}
	} else {
		if qtyIsDelta {
			quantity = sub.Quantity + quantity
			if quantity < 0 {
				return errors.InternalServerError("subscriptions.Update", "Something is wrong, trying to create subscription with negative value")
			}
		}
		logger.Info("Increasing subscription quantity")
		_, err = s.paymentService.UpdateSubscription(ctx, &paymentsproto.UpdateSubscriptionRequest{
			SubscriptionId: sub.Id,
			CustomerId:     customerID,
			CustomerType:   "user",
			PriceId:        priceID,
			Quantity:       quantity,
		}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
		if err != nil {
			return merrors.InternalServerError("signup", "Error updating subscription '%v': %v", sub.Id, err)
		}
	}

	return nil
}

// authorizeAdminCall checks that the context contains an admin token
func authorizeAdminCall(ctx context.Context) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok || account.Issuer != "micro" {
		return errors.Unauthorized("subscriptions", "Unauthorized request")
	}
	return nil
}

// authorizeCall checks that the context contains a token for the customer or is an admin
func authorizeCall(ctx context.Context, customerID string) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok || (account.Issuer != "micro" && account.ID != customerID) {
		return errors.Unauthorized("subscriptions", "Unauthorized request")
	}
	return nil
}
