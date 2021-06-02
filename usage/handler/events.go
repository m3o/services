package handler

import (
	"encoding/json"
	"fmt"
	"time"

	pevents "github.com/m3o/services/pkg/events"
	v1api "github.com/m3o/services/v1api/proto"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (p *UsageSvc) consumeEvents() {
	go pevents.ProcessTopic("v1api", "usage", p.processV1apiEvents)
}

func (p *UsageSvc) processV1apiEvents(ev mevents.Event) error {
	ve := &v1api.Event{}
	if err := json.Unmarshal(ev.Payload, ve); err != nil {
		logger.Errorf("Error unmarshalling v1api event: $s", err)
		return nil
	}
	switch ve.Type {
	case "Request":
		if err := p.processRequest(ve.Request, ev.Timestamp); err != nil {
			logger.Errorf("Error processing request event %s", err)
			return err
		}
	default:
		logger.Infof("Unrecognised event %+v", ve)
	}
	return nil

}

func (p *UsageSvc) processRequest(event *v1api.RequestEvent, t time.Time) error {
	_, err := p.c.incr(event.UserId, event.ApiName, 1, t)
	p.c.incr(event.UserId, fmt.Sprintf("%s$%s", event.ApiName, event.EndpointName), 1, t)
	// incr total counts for the API and individual endpoint
	p.c.incr("total", event.ApiName, 1, t)
	p.c.incr("total", fmt.Sprintf("%s$%s", event.ApiName, event.EndpointName), 1, t)
	return err
}
