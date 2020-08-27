package handler

import (
	"encoding/json"
	"time"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service/context"
	mevents "github.com/micro/micro/v3/service/events"
	pb "github.com/micro/micro/v3/service/events/proto"
	"github.com/micro/micro/v3/service/events/util"
)

type CustomerEvent struct {
	Type     string
	Customer CustomerModel
}

type SubscriptionEvent struct {
	Type         string
	Subscription SubscriptionModel
}

// TODO copy pasted from subscription service, need centralised events repo
type SubscriptionModel struct {
	// ID in this service
	ID string
	// ID in the payment service
	PaymentSubscriptionID string
	CustomerID            string
	// developer, additional
	Type    string
	Created int64
	// If this sub has been cancelled this represents the end date
	Expires int64
	// If this subscription is paid for by another subscription, this is populated with the ID of the paying sybscription
	ParentSubscriptionID string
}

func ConsumeEvents() {
	go func() {
		var events <-chan events.Event
		start := time.Now()
		for {
			var err error
			events, err = mevents.Subscribe("subscriptions")
			if err == nil {
				break
			}
			// TODO fix me
			if time.Since(start) > 2*time.Minute {
				logger.Fatalf("Failed to subscribe to subscriptions topic %s", err) // TODO should be fatal
			}
			time.Sleep(20 * time.Second)
		}
		go processSubscriptionEvents(events)

	}()

}

func (c Customers) eventSubscribe(topic string, opts ...events.SubscribeOption) (<-chan events.Event, error) {
	// parse options
	var options events.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	// start the stream
	stream, err := c.streamService.Subscribe(context.DefaultContext, &pb.SubscribeRequest{
		Topic:       topic,
		Queue:       options.Queue,
		StartAtTime: options.StartAtTime.Unix(),
	}, goclient.WithAuthToken())
	if err != nil {
		return nil, err
	}

	evChan := make(chan events.Event)
	go func() {
		for {
			ev, err := stream.Recv()
			if err != nil {
				close(evChan)
				return
			}

			evChan <- util.DeserializeEvent(ev)
		}
	}()

	return evChan, nil
}

// TODO remove this and replace with publish from micro/micro
func (c *Customers) eventPublish(topic string, msg interface{}, opts ...events.PublishOption) error {
	// parse the options
	options := events.PublishOptions{
		Timestamp: time.Now(),
	}
	for _, o := range opts {
		o(&options)
	}

	// encode the message if it's not already encoded
	var payload []byte
	if p, ok := msg.([]byte); ok {
		payload = p
	} else {
		p, err := json.Marshal(msg)
		if err != nil {
			return events.ErrEncodingMessage
		}
		payload = p
	}

	// execute the RPC
	_, err := c.streamService.Publish(context.DefaultContext, &pb.PublishRequest{
		Topic:     topic,
		Payload:   payload,
		Metadata:  options.Metadata,
		Timestamp: options.Timestamp.Unix(),
	}, goclient.WithAuthToken())

	return err
}

func processSubscriptionEvents(ch <-chan events.Event) {
	// TODO need a mechanism to return the message to the queue for retry
	for ev := range ch {
		sub := &SubscriptionEvent{}
		if err := json.Unmarshal(ev.Payload, sub); err != nil {
			logger.Errorf("Error unmarshalling subscription event: $s", err)
			continue
		}
		switch sub.Type {
		case "subscriptions.created":
			if _, err := updateCustomerStatus(sub.Subscription.CustomerID, statusActive); err != nil {
				logger.Errorf("Error updating customers status for customers %s. %s", sub.Subscription.CustomerID, err)
				continue
			}
			logger.Infof("Updated customer status to active from subscriptions.created event %+v", sub)
		}

	}
	// TODO what do you do if the channel closes
}
