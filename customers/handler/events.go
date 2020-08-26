package handler

import (
	"encoding/json"

	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/logger"
	mevents "github.com/micro/micro/v3/service/events"
)

type CustomerEvent struct {
	Type     string
	Customer CustomerModel
}

func init() {
	events, err := mevents.Subscribe("subscriptions")
	if err != nil {
		logger.Fatalf("Failed to subscribe to payments event stream")
	}
	go processSubscriptionEvents(events)
	// TODO
}

type SubscriptionEvent struct {
	Type       string
	CustomerID string
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
		case "subscription.created":
			if _, err := updateCustomerStatus(sub.CustomerID, statusActive); err != nil {
				logger.Errorf("Error updating customers status for customers %s. %s", sub.CustomerID, err)
				continue
			}
		}

	}
	// TODO what do you do if the channel closes
}
