package dao

import (
	"time"

	"github.com/google/uuid"
)

// Notification is the notification object.
type Notification struct {
	ID           string
	Message      string
	ResourceID   string
	ResourceType string
	Created      time.Time
	Read         time.Time
	SubscriberID string
}

// NewNotification returns a new notification object
func NewNotification(msg, resourceID, resourceType, subscriberID string) *Notification {
	notif := &Notification{
		ID:           uuid.New().String(),
		Message:      msg,
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Created:      time.Now(),
		Read:         time.Time{},
		SubscriberID: subscriberID,
	}
	return notif
}

// Subscription is the subscription object.
type Subscription struct {
	ID           string
	SubscriberID string
	ResourceType string
	ResourceID   string
	Created      time.Time
}

// NewSubscription returns a new subscription object
func NewSubscription(subscriberID, resourceType, resourceID string) *Subscription {
	sub := &Subscription{
		ID:           uuid.New().String(),
		SubscriberID: subscriberID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}
	return sub
}
