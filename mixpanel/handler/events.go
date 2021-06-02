package handler

import (
	"encoding/json"

	balance "github.com/m3o/services/balance/proto"
	customers "github.com/m3o/services/customers/proto"
	"github.com/m3o/services/pkg/events"
	v1api "github.com/m3o/services/v1api/proto"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (b *Mixpanel) consumeEvents() {
	go events.ProcessTopic("v1api", "mixpanel", b.processV1APIEvent)
	go events.ProcessTopic("balance", "mixpanel", b.processBalanceEvent)
	go events.ProcessTopic("customers", "mixpanel", b.processCustomerEvent)
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
