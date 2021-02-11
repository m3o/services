package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service/client"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
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
				logger.Errorf("Error processing request event %s", err)
				continue
			}
		default:
			logger.Infof("Unrecognised event %+v", ve)

		}
		ev.Ack()
	}
}

func (q *Quota) processAPIKeyCreated(ac *v1api.APIKeyCreateEvent) error {
	// register the user for quotas to pick up any new scopes that may have been assigned
	// TODO - remove registration here to somewhere more appropriate like when we get a front end
	recs, err := store.Read(fmt.Sprintf("%s:", prefixQuotaID), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		logger.Errorf("Error looking up quotas %s", err)
		return err
	}
	candidatePaths := []string{}
	for _, s := range ac.Scopes {
		// if the scope is something like location:write then we just want the first part
		parts := strings.Split(s, ":")
		candidatePaths = append(candidatePaths, fmt.Sprintf("/%s/", parts[0]))
	}
	quotasToAdd := []string{}
	for _, r := range recs {
		quot := &quota{}
		if err := json.Unmarshal(r.Value, quot); err != nil {
			logger.Errorf("Error unmarshalling quota %s", err)
			return err
		}
		for _, cand := range candidatePaths {
			if strings.HasPrefix(quot.Path, cand) {
				quotasToAdd = append(quotasToAdd, quot.ID)
				break
			}
		}
	}
	logger.Infof("Adding quotas for user %s:%s %v", ac.UserId, ac.Namespace, quotasToAdd)

	if err := q.registerUser(ac.UserId, ac.Namespace, quotasToAdd); err != nil {
		logger.Errorf("Error registering user %s", err)
		return err
	}

	// Keys start in blocked status, so work out which paths should be unblocked for the key to operate
	recs, err = store.Read(fmt.Sprintf("%s:%s:%s", prefixMapping, ac.Namespace, ac.UserId), store.ReadPrefix())
	if err != nil {
		return err
	}

	allowList := []string{}
	for _, r := range recs {
		m := &mapping{}
		if err := json.Unmarshal(r.Value, m); err != nil {
			return err
		}
		breach, p, err := q.hasBreachedLimit(m.Namespace, m.UserID, m.QuotaID)
		if err != nil {
			return err
		}
		if breach {
			continue
		}

		allowList = append(allowList, p)
	}
	// update the key to unblock it
	if _, err := q.v1Svc.UpdateAllowedPaths(context.TODO(), &v1api.UpdateAllowedPathsRequest{
		UserId:    ac.UserId,
		Namespace: ac.Namespace,
		Allowed:   allowList,
		KeyId:     ac.ApiKeyId,
	}, client.WithAuthToken()); err != nil {
		logger.Errorf("Error updating allowed paths %s", err)
		return err
	}

	return nil
}

func (q *Quota) hasBreachedLimit(namespace, userID, quotaID string) (bool, string, error) {
	qrecs, err := store.Read(fmt.Sprintf("%s:%s", prefixQuotaID, quotaID))
	if err != nil {
		return false, "", err
	}
	quot := &quota{}
	if err := json.Unmarshal(qrecs[0].Value, quot); err != nil {
		return false, "", err
	}
	// check the current count
	count, err := q.c.read(namespace, userID, quot.Path)
	if err != nil && err != redis.Nil {
		return false, "", err
	}
	return quot.Limit > 0 && count >= quot.Limit, quot.Path, nil

}

func (q *Quota) processRequest(rqe *v1api.RequestEvent) error {

	// count the request
	// count is coarse granularity - we just care which service they've called so /v1/blah
	if !strings.HasPrefix(rqe.Url, "/v1/") {
		logger.Warnf("Discarding unrecognised URL Path %s", rqe.Url)
		return nil
	}
	parts := strings.Split(rqe.Url[1:], "/")
	if len(parts) < 2 {
		logger.Warnf("Discarding unrecognised URL Path %s", rqe.Url)
		return nil
	}
	reqPath := fmt.Sprintf("/%s/", parts[1])

	curr, err := q.c.incr(rqe.Namespace, rqe.UserId, reqPath)
	if err != nil {
		return err
	}
	logger.Infof("Current count is %s %s %s %d", rqe.Namespace, rqe.UserId, reqPath, curr)

	recs, err := store.Read(fmt.Sprintf("%s:%s:%s", prefixMapping, rqe.Namespace, rqe.UserId), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		logger.Errorf("Error getting quotas %s", err)
		return err
	}

	for _, r := range recs {
		m := &mapping{}
		if err := json.Unmarshal(r.Value, m); err != nil {
			logger.Errorf("Error unmarshalling mapping %s", err)
			return err
		}

		breach, p, err := q.hasBreachedLimit(m.Namespace, m.UserID, m.QuotaID)
		if err != nil {
			logger.Errorf("Error calculating breach %s", err)
			return err
		}
		if p != reqPath || !breach {
			continue
		}

		// update the user to block requests
		if _, err := q.v1Svc.UpdateAllowedPaths(context.TODO(), &v1api.UpdateAllowedPathsRequest{
			UserId:    m.UserID,
			Namespace: m.Namespace,
			Blocked:   []string{reqPath},
		}, client.WithAuthToken()); err != nil {
			logger.Errorf("Error updating allowed paths %s", err)
			return err
		}
		break
	}
	return nil
}
