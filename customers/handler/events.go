package handler

import (
	"context"
	"encoding/json"
	"time"

	mevents "github.com/micro/micro/v3/service/events"

	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/logger"
)

type CustomerEvent struct {
	Type     string
	Customer CustomerModel
}

type SubscriptionEvent struct {
	Type         string
	Subscription SubscriptionModel
}

// TODO copy pasted from subscription service, need centralised events repo
type SubscriptionModel struct {
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

func (c *Customers) consumeEvents() {
	var evs <-chan events.Event
	start := time.Now()
	for {
		var err error
		evs, err = mevents.Subscribe("subscriptions",
			events.WithAutoAck(false, 30*time.Second),
			events.WithRetryLimit(10)) // 10 retries * 30 secs ackWait gives us 5 mins of tolerance for issues
		if err == nil {
			c.processSubscriptionEvents(evs)
			start = time.Now()
			continue // if for some reason evs closes we loop and try subscribing again
		}
		// TODO fix me
		if time.Since(start) > 2*time.Minute {
			logger.Fatalf("Failed to subscribe to subscriptions topic %s", err)
		}
		logger.Warnf("Unable to subscribe to evs %s. Will retry in 20 secs", err)
		time.Sleep(20 * time.Second)
	}

}

func (c *Customers) processSubscriptionEvents(ch <-chan events.Event) {
	for ev := range ch {
		sub := &SubscriptionEvent{}
		if err := json.Unmarshal(ev.Payload, sub); err != nil {
			ev.Nack()
			logger.Errorf("Error unmarshalling subscription event: $s", err)
			continue
		}
		switch sub.Type {
		case "subscriptions.created":
			if _, err := updateCustomerStatusByID(sub.Subscription.CustomerID, statusActive); err != nil {
				ev.Nack()
				logger.Errorf("Error updating customers status for customer %s. %s", sub.Subscription.CustomerID, err)
				continue
			}
			logger.Infof("Updated customer status to active from subscriptions.created event %+v", sub)
		case "subscriptions.cancelled":
			if err := c.processCancelledSubscription(&sub.Subscription); err != nil {
				ev.Nack()
				logger.Errorf("Error processing subscription cancel for customer %s. %s", sub.Subscription.CustomerID, err)
				continue
			}
			logger.Infof("Processed subscriptions.cancelled event %+v", sub)
		}
		ev.Ack()
	}
}

func (c *Customers) processCancelledSubscription(sub *SubscriptionModel) error {
	return c.deleteCustomer(context.Background(), sub.CustomerID)

}
