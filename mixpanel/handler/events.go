package handler

import (
	"encoding/json"
	"fmt"
	"time"

	balance "github.com/m3o/services/balance/proto"
	customers "github.com/m3o/services/customers/proto"
	v1api "github.com/m3o/services/v1api/proto"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (b *Mixpanel) consumeEvents() {
	go b.processTopic("v1api", b.processV1APIEvent)
	go b.processTopic("balance", b.processBalanceEvent)
	go b.processTopic("customers", b.processCustomerEvent)
}

func (b *Mixpanel) processTopic(topic string, handler func(ev mevents.Event) error) {
	var evs <-chan mevents.Event
	start := time.Now()
	for {
		var err error
		evs, err = mevents.Consume(topic,
			mevents.WithAutoAck(false, 30*time.Second),
			mevents.WithRetryLimit(10),
			mevents.WithGroup(fmt.Sprintf("%s-%s", "mixpanel", topic))) // 10 retries * 30 secs ackWait gives us 5 mins of tolerance for issues
		if err == nil {
			logger.Infof("Starting to process %s events", topic)
			b.processEvents(evs, handler)
			start = time.Now()
			continue // if for some reason evs closes we loop and try subscribing again
		}
		// TODO fix me
		if time.Since(start) > 2*time.Minute {
			logger.Fatalf("Failed to subscribe to topic %s: %s", topic, err)
		}
		logger.Warnf("Unable to subscribe to topic %s. Will retry in 20 secs. %s", topic, err)
		time.Sleep(20 * time.Second)
	}
}

func (b *Mixpanel) processEvents(ch <-chan mevents.Event, handler func(ev mevents.Event) error) {
	for {
		t := time.NewTimer(60 * time.Minute)
		var ev mevents.Event
		select {
		case ev = <-ch:
			t.Stop()
			if len(ev.ID) == 0 {
				// channel closed
				logger.Infof("Channel closed, retrying stream connection")
				return
			}
		case <-t.C:
			// safety net in case we stop receiving messages for some reason
			logger.Infof("No messages received for last 60 minutes retrying connection")
			return
		}
		if err := handler(ev); err != nil {
			ev.Nack()
			continue
		}
		ev.Ack()
	}
}

func (b *Mixpanel) processV1APIEvent(ev mevents.Event) error {
	ve := &v1api.Event{}

	if err := json.Unmarshal(ev.Payload, ve); err != nil {
		logger.Errorf("Error unmarshalling v1api event, discarding: $s", err)
		return nil
	}

	customerID := ""
	switch ve.Type {
	case "Request":
		customerID = ve.Request.UserId
	case "APIKeyCreate":
		customerID = ve.ApiKeyCreate.UserId
	case "APIKeyRevoke":
		customerID = ve.ApiKeyRevoke.UserId
	default:
		logger.Infof("Event type for v1api not supported %s", ve.Type)
		return nil
	}

	mEv := b.client.newMixpanelEvent(ev.Topic, ve.Type, customerID, ev.ID, ev.Timestamp.Unix(), ve)
	if err := b.client.Track(mEv); err != nil {
		logger.Errorf("Error tracking event %s", err)
		return err
	}
	return nil

}

func (b *Mixpanel) processBalanceEvent(ev mevents.Event) error {
	ve := &balance.Event{}
	if err := json.Unmarshal(ev.Payload, ve); err != nil {
		logger.Errorf("Error unmarshalling balance event, discarding: $s", err)
		return nil
	}
	customerID := ve.CustomerId

	mEv := b.client.newMixpanelEvent(ev.Topic, ve.Type.String(), customerID, ev.ID, ev.Timestamp.Unix(), ve)
	if err := b.client.Track(mEv); err != nil {
		logger.Errorf("Error tracking event %s", err)
		return err
	}

	return nil
}

func (b *Mixpanel) processCustomerEvent(ev mevents.Event) error {
	ve := &customers.Event{}
	if err := json.Unmarshal(ev.Payload, ve); err != nil {
		logger.Errorf("Error unmarshalling balance event, discarding: $s", err)
		return nil
	}
	customerID := ve.Customer.Id

	mEv := b.client.newMixpanelEvent(ev.Topic, ve.Type.String(), customerID, ev.ID, ev.Timestamp.Unix(), ve)
	if err := b.client.Track(mEv); err != nil {
		logger.Errorf("Error tracking event %s", err)
		return err
	}

	return nil
}
