package handler

import (
	"context"
	"encoding/json"

	cpb "github.com/m3o/services/customers/proto"
	pevents "github.com/m3o/services/pkg/events"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (v *V1) consumeEvents() {
	go pevents.ProcessTopic("customers", "v1api", v.processCustomerEvents)
}

func (v1 *V1) processCustomerEvents(ev mevents.Event) error {
	ctx := context.Background()
	ce := &cpb.Event{}
	if err := json.Unmarshal(ev.Payload, ce); err != nil {
		logger.Errorf("Error unmarshalling customer event: $s", err)
		return nil
	}
	switch ce.Type {
	case cpb.EventType_EventTypeDeleted:
		if err := v1.processCustomerDelete(ctx, ce); err != nil {
			logger.Errorf("Error processing request event %s", err)
			return err
		}
	default:
		logger.Infof("Skipping event %+v", ce)
	}
	return nil

}

func (v1 *V1) processCustomerDelete(ctx context.Context, event *cpb.Event) error {
	return v1.deleteCustomer(ctx, event.Customer.Id)
}
