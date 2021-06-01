package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/m3o/services/balance/proto"
	stripepb "github.com/m3o/services/stripe/proto"
	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	mevents "github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

const (
	msgInsufficientFunds = "Insufficient funds"
)

func (b *Balance) consumeEvents() {
	processTopic := func(topic string, handler func(ch <-chan mevents.Event)) {
		var evs <-chan mevents.Event
		start := time.Now()
		for {
			var err error
			evs, err = mevents.Consume(topic,
				mevents.WithAutoAck(false, 30*time.Second),
				mevents.WithRetryLimit(10),
				mevents.WithGroup(fmt.Sprintf("%s-%s", "balance", topic))) // 10 retries * 30 secs ackWait gives us 5 mins of tolerance for issues
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
	go processTopic("v1api", b.processV1apiEvents)
	go processTopic("stripe", b.processStripeEvents)
}

func (b *Balance) processV1apiEvents(ch <-chan mevents.Event) {
	logger.Infof("Starting to process v1api events")
	for {
		t := time.NewTimer(600 * time.Minute)
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
			logger.Infof("No messages received for last 2 minutes retrying connection")
			return
		}

		ve := &v1api.Event{}
		if err := json.Unmarshal(ev.Payload, ve); err != nil {
			ev.Nack()
			logger.Errorf("Error unmarshalling v1api event: $s", err)
			continue
		}
		switch ve.Type {
		case "APIKeyCreate":
			if err := b.processAPIKeyCreated(ve.ApiKeyCreate); err != nil {
				ev.Nack()
				logger.Errorf("Error processing API key created event %s", err)
				continue
			}
		case "Request":
			if err := b.processRequest(ve.Request); err != nil {
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

func (b *Balance) processAPIKeyCreated(ac *v1api.APIKeyCreateEvent) error {
	currBal, err := b.c.read(ac.UserId, "$balance$")
	if err != nil {
		return err
	}

	// Keys start in blocked status, so unblock if they have the cash
	if currBal <= 0 {
		if _, err := b.v1Svc.BlockKey(context.Background(), &v1api.BlockKeyRequest{
			UserId:    ac.UserId,
			Namespace: ac.Namespace,
			Message:   msgInsufficientFunds,
		}, client.WithAuthToken()); err != nil {
			if merr, ok := err.(*errors.Error); ok && merr.Code == 404 {
				return nil
			}
			logger.Errorf("Error blocking key %s", err)
			return err
		}
		return nil
	}
	if _, err := b.v1Svc.UnblockKey(context.Background(), &v1api.UnblockKeyRequest{
		UserId:    ac.UserId,
		Namespace: ac.Namespace,
		KeyId:     ac.ApiKeyId,
	}, client.WithAuthToken()); err != nil {
		if merr, ok := err.(*errors.Error); ok && merr.Code == 404 {
			return nil
		}
		logger.Errorf("Error unblocking key %s", err)
		return err
	}
	return nil
}

func (b *Balance) processRequest(rqe *v1api.RequestEvent) error {
	apiName := rqe.ApiName
	rsp, err := b.pubSvc.get(apiName)
	if err != nil {
		logger.Errorf("Error looking up API %s", err)
		return err
	}

	methodName := rqe.EndpointName
	price, ok := rsp.Pricing[methodName]
	if !ok {
		logger.Warnf("Failed to find price for api call %s:%s", apiName, methodName)
		return nil
	}
	// decrement the balance
	currBal, err := b.c.decr(rqe.UserId, "$balance$", price)
	if err != nil {
		return err
	}

	if currBal > 0 {
		return nil
	}

	evt := pb.Event{
		Type:       pb.EventType_EventTypeZeroBalance,
		CustomerId: rqe.UserId,
	}
	if err := events.Publish(pb.EventsTopic, &evt); err != nil {
		logger.Errorf("Error publishing event %+v", evt)
	}

	// no more money, cut them off
	if _, err := b.v1Svc.BlockKey(context.TODO(), &v1api.BlockKeyRequest{
		UserId:    rqe.UserId,
		Namespace: rqe.Namespace,
		Message:   msgInsufficientFunds,
	}, client.WithAuthToken()); err != nil {
		// TODO if we fail here we might double count because the message will be retried
		if merr, ok := err.(*errors.Error); ok && merr.Code == 404 {
			return nil
		}
		logger.Errorf("Error blocking key %s", err)
		return err
	}

	return nil
}

func (b *Balance) processStripeEvents(ch <-chan mevents.Event) {
	logger.Infof("Starting to process stripe events")
	for {
		t := time.NewTimer(600 * time.Minute)
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
			logger.Infof("No messages received for last 2 minutes retrying connection")
			return
		}

		ve := &stripepb.Event{}
		logger.Infof("Processing event %+v", ev)
		if err := json.Unmarshal(ev.Payload, ve); err != nil {
			ev.Nack()
			logger.Errorf("Error unmarshalling stripe event: $s", err)
			continue
		}
		switch ve.Type {
		case "ChargeSucceeded":
			if err := b.processChargeSucceeded(ve.ChargeSucceeded); err != nil {
				ev.Nack()
				logger.Errorf("Error processing charge succeeded event %s", err)
				continue
			}
		default:
			logger.Infof("Unrecognised event %+v", ve)

		}
		ev.Ack()
	}
}

func (b *Balance) processChargeSucceeded(ev *stripepb.ChargeSuceededEvent) error {
	// TODO if we return error and we have already incremented the counter then we double count so make this idempotent
	// safety first
	if ev == nil || ev.Amount == 0 {
		return nil
	}
	// add to balance
	// stripe event is in cents, multiply by 100 to get the fraction that balance represents
	_, err := b.c.incr(ev.CustomerId, "$balance$", ev.Amount*100)
	if err != nil {
		logger.Errorf("Error incrementing balance %s", err)
	}

	srsp, err := b.stripeSvc.GetPayment(context.Background(), &stripepb.GetPaymentRequest{Id: ev.ChargeId}, client.WithAuthToken())
	if err != nil {
		return err
	}

	adj, err := storeAdjustment(ev.CustomerId, ev.Amount*100, ev.CustomerId, "Funds added", true, map[string]string{
		"receipt_url": srsp.Payment.ReceiptUrl,
	})
	if err != nil {
		return err
	}

	evt := pb.Event{
		Type: pb.EventType_EventTypeIncrement,
		Adjustment: &pb.Adjustment{
			Id:        adj.ID,
			Created:   adj.Created.Unix(),
			Delta:     adj.Amount,
			Reference: adj.Reference,
			Meta:      adj.Meta,
		},
		CustomerId: adj.CustomerID,
	}
	if err := events.Publish(pb.EventsTopic, &evt); err != nil {
		logger.Errorf("Error publishing event %+v", evt)
	}

	namespace := microNamespace

	// unblock key
	if _, err := b.v1Svc.UnblockKey(context.Background(), &v1api.UnblockKeyRequest{
		UserId:    ev.CustomerId,
		Namespace: namespace,
	}, client.WithAuthToken()); err != nil {
		// TODO if we fail here we might double count because the message will be retried
		logger.Errorf("Error unblocking key %s", err)
		return err
	}
	return nil
}
