package handler

import (
	"encoding/json"
	"time"

	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/logger"
	mevents "github.com/micro/micro/v3/service/events"
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

func ConsumeEvents() {
	go func() {
		var events <-chan events.Event
		start := time.Now()
		for {
			var err error
			events, err = mevents.Subscribe("subscriptions")
			if err == nil {
				break
			}
			// TODO fix me
			if time.Since(start) > 2*time.Minute {
				logger.Fatalf("Failed to subscribe to subscriptions topic %s", err) // TODO should be fatal
			}
			logger.Warnf("Unable to subscribe to events %s. Will retru in 20 secs", err)
			time.Sleep(20 * time.Second)
		}
		go processSubscriptionEvents(events)

	}()

}

func processSubscriptionEvents(ch <-chan events.Event) {
	// TODO need a mechanism to return the message to the queue for retry
	for ev := range ch {
		sub := &SubscriptionEvent{}
		if err := json.Unmarshal(ev.Payload, sub); err != nil {
			logger.Errorf("Error unmarshalling subscription event: $s", err)
			continue
		}
		switch sub.Type {
		case "subscriptions.created":
			if _, err := updateCustomerStatus(sub.Subscription.CustomerID, statusActive); err != nil {
				logger.Errorf("Error updating customers status for customers %s. %s", sub.Subscription.CustomerID, err)
				continue
			}
			logger.Infof("Updated customer status to active from subscriptions.created event %+v", sub)
		}

	}
	// TODO what do you do if the channel closes
}
