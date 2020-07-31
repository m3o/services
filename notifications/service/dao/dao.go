package dao

import (
	"encoding/json"
	"fmt"

	"github.com/micro/go-micro/v3/store"
	mstore "github.com/micro/micro/v3/service/store"
)

const (
	prefixNotificationsBySubscriber = "notif:%s:"                            // notif:subscriberID:
	keyNotificationsBySubscriber    = prefixNotificationsBySubscriber + "%s" // notif:subscriberID:notificationID
	prefixSubscriptionsBySubscriber = "sub:%s:"                              // sub:subscriberID:
	keySubscriptionsBySubscriber    = prefixSubscriptionsBySubscriber + "%s" // sub:subscriberID:subscriptionID
	prefixSubscriptionsByResource   = "subres:%s:%s:"                        // subres:resourceType:resourceID:
	keySubscriptionsByResource      = prefixSubscriptionsByResource + "%s"   // subres:resourceType:resourceID:subscriberID

)

// CreateNotification stores the notification
func CreateNotification(notif *Notification) error {
	return writeObj(notif, fmt.Sprintf(keyNotificationsBySubscriber, notif.SubscriberID, notif.ID))
}

func writeObj(obj interface{}, key string) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	if err := mstore.Write(&store.Record{
		Key:   key,
		Value: b,
	}); err != nil {
		return err
	}
	return nil
}

// ListNotificationsForSubscriber returns a list of all notifications for the given subscriber
func ListNotificationsForSubscriber(subscriberID string) ([]*Notification, error) {
	entries, err := mstore.Read(fmt.Sprintf(prefixNotificationsBySubscriber, subscriberID), store.ReadPrefix())
	if err != nil {
		return nil, err
	}
	ret := make([]*Notification, len(entries))
	for i, v := range entries {
		n := &Notification{}
		if err := json.Unmarshal(v.Value, n); err != nil {
			return nil, err
		}
		ret[i] = n
	}
	return ret, nil
}

// UpdateNotification updates the notification
func UpdateNotification(notif *Notification) error {
	return CreateNotification(notif)
}

// ReadNotification returns the notification with the given ID
func ReadNotification(subscriberID, notifID string) (*Notification, error) {
	entries, err := mstore.Read(fmt.Sprintf(keyNotificationsBySubscriber, subscriberID, notifID))
	if err != nil {
		return nil, err
	}
	ret := &Notification{}
	if err := json.Unmarshal(entries[0].Value, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// CreateSubscription stores the subscription
func CreateSubscription(sub *Subscription) error {

	if err := writeObj(sub, fmt.Sprintf(keySubscriptionsBySubscriber, sub.SubscriberID, sub.ID)); err != nil {
		return err
	}
	// TODO should we rollback previous if this one fails?
	return writeObj(sub, fmt.Sprintf(keySubscriptionsByResource, sub.ResourceType, sub.ResourceID, sub.SubscriberID))
}

// ListSubscriptionsForResource returns the subscriptions for a resource
func ListSubscriptionsForResource(resourceType, resourceID string) ([]*Subscription, error) {
	entries, err := mstore.Read(fmt.Sprintf(prefixSubscriptionsByResource, resourceType, resourceID), store.ReadPrefix())
	if err != nil {
		return nil, err
	}
	ret := make([]*Subscription, len(entries))
	for i, v := range entries {
		s := &Subscription{}
		if err := json.Unmarshal(v.Value, s); err != nil {
			return nil, err
		}
		ret[i] = s
	}
	return ret, nil
}

// DeleteSubscription removes the subscription from the store
func DeleteSubscription(subscriptionID, subscriberID, resourceType, resourceID string) error {
	// Normally we'd have a deleted timestamp rather than deleting from the DB but we don't really need an audit trail here
	if err := mstore.Delete(fmt.Sprintf(keySubscriptionsByResource, resourceType, resourceID, subscriberID)); err != nil {
		return err
	}
	return mstore.Delete(fmt.Sprintf(keySubscriptionsBySubscriber, subscriberID, subscriptionID))
}

// ReadSubscription returns a subscription for the given params
func ReadSubscription(subscriberID, resourceType, resourceID string) (*Subscription, error) {
	entries, err := mstore.Read(fmt.Sprintf(keySubscriptionsByResource, resourceType, resourceID, subscriberID))
	if err != nil {
		return nil, err
	}
	ret := &Subscription{}
	if err := json.Unmarshal(entries[0].Value, ret); err != nil {
		return nil, err
	}
	return ret, nil

}

// ListSubscriptionsForSubscriber returns the list of subscriptions for this subscriber
func ListSubscriptionsForSubscriber(subscriberID string) ([]*Subscription, error) {
	entries, err := mstore.Read(fmt.Sprintf(prefixSubscriptionsBySubscriber, subscriberID), store.ReadPrefix())
	if err != nil {
		return nil, err
	}
	ret := make([]*Subscription, len(entries))
	for i, v := range entries {
		sub := &Subscription{}
		if err := json.Unmarshal(v.Value, sub); err != nil {
			return nil, err
		}
		ret[i] = sub
	}
	return ret, nil

}
