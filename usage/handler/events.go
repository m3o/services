package handler

import (
	"encoding/json"
	"fmt"
	"time"

	v1api "github.com/m3o/services/v1api/proto"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (p *UsageSvc) consumeEvents() {
	processTopic := func(topic string, handler func(ch <-chan mevents.Event)) {
		var evs <-chan mevents.Event
		start := time.Now()
		for {
			var err error
			evs, err = mevents.Consume(topic,
				mevents.WithAutoAck(false, 30*time.Second),
				mevents.WithRetryLimit(10), // 10 retries * 30 secs ackWait gives us 5 mins of tolerance for issues
				mevents.WithGroup(fmt.Sprintf("%s-%s", "usage", topic)))
			if err == nil {
				handler(evs)
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
	go processTopic("v1api", p.processV1apiEvents)
}

func (p *UsageSvc) processV1apiEvents(ch <-chan mevents.Event) {
	logger.Infof("Starting to process v1api events")
	for ev := range ch {
		ve := &v1api.Event{}
		if err := json.Unmarshal(ev.Payload, ve); err != nil {
			ev.Nack()
			logger.Errorf("Error unmarshalling v1api event: $s", err)
			continue
		}
		switch ve.Type {
		case "Request":
			if err := p.processRequest(ve.Request, ev.Timestamp); err != nil {
				ev.Nack()
				logger.Errorf("Error processing request event %s", err)
				continue
			}
		default:
			logger.Infof("Unrecognised event %+v", ve)
		}
		ev.Ack()
	}
}

func (p *UsageSvc) processRequest(event *v1api.RequestEvent, t time.Time) error {
	_, err := p.c.incr(event.UserId, event.ApiName, 1, t)
	p.c.incr(event.UserId, fmt.Sprintf("%s$%s", event.ApiName, event.EndpointName), 1, t)
	// incr total counts for the API and individual endpoint
	p.c.incr("total", event.ApiName, 1, t)
	p.c.incr("total", fmt.Sprintf("%s$%s", event.ApiName, event.EndpointName), 1, t)
	return err
}
