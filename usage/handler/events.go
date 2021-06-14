package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cpb "github.com/m3o/services/customers/proto"
	pevents "github.com/m3o/services/pkg/events"
	v1api "github.com/m3o/services/v1api/proto"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (p *UsageSvc) consumeEvents() {
	go pevents.ProcessTopic("v1api", "usage", p.processV1apiEvents)
	go pevents.ProcessTopic("customers", "usage", p.processCustomerEvents)
}

func (p *UsageSvc) processV1apiEvents(ev mevents.Event) error {
	ctx := context.Background()
	ve := &v1api.Event{}
	if err := json.Unmarshal(ev.Payload, ve); err != nil {
		logger.Errorf("Error unmarshalling v1api event: $s", err)
		return nil
	}
	switch ve.Type {
	case "Request":
		if err := p.processRequest(ctx, ve.Request, ev.Timestamp); err != nil {
			logger.Errorf("Error processing request event %s", err)
			return err
		}
	default:
		logger.Infof("Unrecognised event %+v", ve)
	}
	return nil

}

func (p *UsageSvc) processRequest(ctx context.Context, event *v1api.RequestEvent, t time.Time) error {
	_, err := p.c.incr(ctx, event.UserId, event.ApiName, 1, t)
	p.c.incr(ctx, event.UserId, fmt.Sprintf("%s$%s", event.ApiName, event.EndpointName), 1, t)
	// incr total counts for the API and individual endpoint
	p.c.incr(ctx, "total", event.ApiName, 1, t)
	p.c.incr(ctx, "total", fmt.Sprintf("%s$%s", event.ApiName, event.EndpointName), 1, t)
	return err
}

func (p *UsageSvc) processCustomerEvents(ev mevents.Event) error {
	ctx := context.Background()
	ce := &cpb.Event{}
	if err := json.Unmarshal(ev.Payload, ce); err != nil {
		logger.Errorf("Error unmarshalling customer event: $s", err)
		return nil
	}
	switch ce.Type {
	case cpb.EventType_EventTypeDeleted:
		if err := p.processCustomerDelete(ctx, ce); err != nil {
			logger.Errorf("Error processing request event %s", err)
			return err
		}
	default:
		logger.Infof("Skipping event %+v", ce)
	}
	return nil

}

func (p *UsageSvc) processCustomerDelete(ctx context.Context, event *cpb.Event) error {
	// delete all their usage
	return p.deleteUser(ctx, event.Customer.Id)
}
