package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/client"

	v1api "github.com/m3o/services/v1api/proto"

	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func (q *Quota) consumeEvents() {
	processTopic := func(topic string, handler func(ch <-chan mevents.Event)) {
		var evs <-chan mevents.Event
		start := time.Now()
		for {
			var err error
			evs, err = mevents.Consume(topic,
				mevents.WithAutoAck(false, 30*time.Second),
				mevents.WithRetryLimit(10)) // 10 retries * 30 secs ackWait gives us 5 mins of tolerance for issues
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
	go processTopic("v1api", q.processV1apiEvents)

}

func (q *Quota) processV1apiEvents(ch <-chan mevents.Event) {
	logger.Infof("Starting to process v1api events")
	for ev := range ch {
		ve := &v1api.Event{}
		if err := json.Unmarshal(ev.Payload, ve); err != nil {
			ev.Nack()
			logger.Errorf("Error unmarshalling v1api event: $s", err)
			continue
		}
		switch ve.Type {
		case "APIKeyCreate":
			// update the key with any allowed paths
			if err := q.processAPIKeyCreated(ve.ApiKeyCreate); err != nil {
				ev.Nack()
				logger.Errorf("Error processing API key created event")
				continue
			}
		case "Request":
			// update the count
			if err := q.processRequest(ve.Request); err != nil {
				ev.Nack()
				logger.Errorf("Error processing request event")
				continue
			}
		default:
			logger.Infof("Unrecognised event %+v", ve)

		}
		ev.Ack()
	}
}

func (q *Quota) processAPIKeyCreated(ac *v1api.APIKeyCreateEvent) error {
	// TODO register this for quotas

	// update the key to unblock it
	if _, err := q.v1Svc.UpdateAllowedPaths(context.TODO(), &v1api.UpdateAllowedPathsRequest{
		UserId:    ac.UserId,
		Namespace: ac.Namespace,
		Allowed:   []string{"/v1/"}, // TODO this unblocks all currently
		Blocked:   nil,
		KeyId:     ac.ApiKeyId,
	}, client.WithAuthToken()); err != nil {
		logger.Errorf("Error updating allowed paths %s", err)
		return err
	}

	return nil
}

func (q *Quota) processRequest(rqe *v1api.RequestEvent) error {

	// count the request
	// count is coarse granularity - we just care which service they've called so /v1/blah
	if !strings.HasPrefix(rqe.Url, "/v1/") {
		logger.Warnf("Discarding unrecognised URL path %s", rqe.Url)
		return nil
	}
	parts := strings.Split(rqe.Url[1:], "/")
	if len(parts) < 2 {
		logger.Warnf("Discarding unrecognised URL path %s", rqe.Url)
		return nil
	}

	curr, err := q.c.incr(fmt.Sprintf("%s:%s:%s", rqe.Namespace, rqe.UserId, parts[1]))
	if err != nil {
		return err
	}
	// TODO do post processing - do we need to block the user because of an exhausted quota?
	logger.Infof("Current count is %d", curr)
	return nil
}
