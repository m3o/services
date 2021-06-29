package handler

import (
	"context"
	"encoding/json"

	balance "github.com/m3o/services/balance/proto"
	customers "github.com/m3o/services/customers/proto"
	"github.com/m3o/services/pkg/events"
	v1 "github.com/m3o/services/v1/proto"
	"github.com/micro/micro/v3/service/client"
	merrors "github.com/micro/micro/v3/service/errors"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (b *Mixpanel) consumeEvents() {
	go events.ProcessTopic("v1api", "mixpanel", b.processV1APIEvent)
	go events.ProcessTopic("balance", "mixpanel", b.processBalanceEvent)
	go events.ProcessTopic("customers", "mixpanel", b.processCustomerEvent)
}

func (b *Mixpanel) processV1APIEvent(ev mevents.Event) error {
	ve := &v1.Event{}

	if err := json.Unmarshal(ev.Payload, ve); err != nil {
		logger.Errorf("Error unmarshalling v1 event, discarding: $s", err)
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
		logger.Infof("Event type for v1 not supported %s", ve.Type)
		return nil
	}

	ignore, err := b.ignoreCustomer(customerID)
	if ignore {
		return nil
	}
	if err != nil {
		return err
	}

	mEv := b.client.newMixpanelEvent(ev.Topic, ve.Type, customerID, ev.ID, ev.Timestamp.Unix(), ve)
	if err := b.client.Track(mEv); err != nil {
		logger.Errorf("Error tracking event %s", err)
		return err
	}
	return nil

}

func (b *Mixpanel) ignoreCustomer(customerID string) (bool, error) {
	// should we ignore?
	rsp, err := b.custSvc.Read(context.Background(), &customers.ReadRequest{Id: customerID}, client.WithAuthToken())
	if err != nil {
		merr, ok := err.(*merrors.Error)
		if ok {
			if merr.Code == 404 {
				logger.Errorf("Failed to find customer %s", customerID)
				return true, nil
			}
		}
		logger.Errorf("Failed to read customer %s %s", customerID, err)
		return false, err
	}
	for _, i := range b.ignoreList {
		if rsp.Customer.Email == i {
			return true, nil
		}
	}
	return false, nil

}

func (b *Mixpanel) processBalanceEvent(ev mevents.Event) error {
	ve := &balance.Event{}
	if err := json.Unmarshal(ev.Payload, ve); err != nil {
		logger.Errorf("Error unmarshalling balance event, discarding: $s", err)
		return nil
	}
	customerID := ve.CustomerId
	ignore, err := b.ignoreCustomer(customerID)
	if ignore {
		return nil
	}
	if err != nil {
		return err
	}

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
	for _, i := range b.ignoreList {
		if ve.Customer.Email == i {
			return nil
		}
	}

	mEv := b.client.newMixpanelEvent(ev.Topic, ve.Type.String(), customerID, ev.ID, ev.Timestamp.Unix(), ve)
	if err := b.client.Track(mEv); err != nil {
		logger.Errorf("Error tracking event %s", err)
		return err
	}

	return nil
}
