package handler

import (
	"encoding/json"
	"time"

	mstore "github.com/micro/micro/v3/service/store"

	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/logger"
	mevents "github.com/micro/micro/v3/service/events"
)

type SignupEvent struct {
	Type   string
	Signup SignupModel
}

type SignupModel struct {
	Email      string
	Namespace  string
	CustomerID string
}

func (s *Signup) consumeEvents() {
	var evs <-chan events.Event
	start := time.Now()
	for {
		var err error
		evs, err = mevents.Subscribe("signup",
			events.WithAutoAck(false, 30*time.Second),
			events.WithRetryLimit(10)) // 10 retries * 30 secs ackWait gives us 5 mins of tolerance for issues
		if err == nil {
			s.processSignupEvents(evs)
			start = time.Now()
			continue // if for some reason evs closes we loop and try subscribing again
		}
		// TODO fix me
		if time.Since(start) > 2*time.Minute {
			logger.Fatalf("Failed to subscribe to subscriptions topic %s", err)
		}
		logger.Warnf("Unable to subscribe to evs %s. Will retry in 20 secs", err)
		time.Sleep(20 * time.Second)
	}

}

func (s *Signup) processSignupEvents(ch <-chan events.Event) {
	for ev := range ch {
		sue := &SignupEvent{}
		if err := json.Unmarshal(ev.Payload, sub); err != nil {
			ev.Nack()
			logger.Errorf("Error unmarshalling subscription event: $s", err)
			continue
		}
		switch sue.Type {
		case "signup.completed":
			// do cleanup of any data so that customer could signup again if they cancelled
			if err := s.processSignupCompleted(sue); err != nil {
				ev.Nack()
				logger.Errorf("Error processing %s event for customer %s. %s", sue.Type, sue.Signup.CustomerID, err)
				continue
			}
		}
		ev.Ack()
	}
}

func (s *Signup) processSignupCompleted(signup *SignupModel) error {
	// delete all the things
	// payment method
	if err := mstore.Delete(prefixPaymentMethod + signup.Email); err != nil {
		return err
	}
	recs, err := mstore.Read(signup.Email)
	if err != nil {
		return err
	}
	tok := &tokenToEmail{}
	if err := json.Unmarshal(recs[0].Value, tok); err != nil {
		return err
	}
	if err := mstore.Delete(signup.Email); err != nil {
		return err
	}
	if err := mstore.Delete(tok.Token); err != nil {
		return err
	}

	return nil
}
