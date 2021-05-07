package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	balance "github.com/m3o/services/balance/proto"
	onboarding "github.com/m3o/services/onboarding/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	mevents "github.com/micro/micro/v3/service/events"
	logger "github.com/micro/micro/v3/service/logger"
)

type Onboarding struct {
	balSvc       balance.BalanceService
	promoCredit  int64
	promoMessage string
}

func NewOnboarding(svc *service.Service) *Onboarding {
	cfg := struct {
		PromoCredit  int64
		PromoMessage string
	}{}
	v, err := config.Get("micro.onboarding")
	if err != nil {
		log.Fatalf("Failed to load config %s", err)
	}
	if err := v.Scan(cfg); err != nil {
		log.Fatalf("Failed to load config %s", err)
	}
	if cfg.PromoCredit == 0 || len(cfg.PromoMessage) == 0 {
		log.Fatalf("Missing config")
	}
	o := &Onboarding{
		balSvc:       balance.NewBalanceService("balance", svc.Client()),
		promoCredit:  cfg.PromoCredit,
		promoMessage: cfg.PromoMessage,
	}
	o.consumeEvents()
	return o
}

func (o *Onboarding) consumeEvents() {
	processTopic := func(topic string, handler func(ch <-chan mevents.Event)) {
		var evs <-chan mevents.Event
		start := time.Now()
		for {
			var err error
			evs, err = mevents.Consume(topic,
				mevents.WithAutoAck(false, 30*time.Second),
				mevents.WithRetryLimit(10), // 10 retries * 30 secs ackWait gives us 5 mins of tolerance for issues
				mevents.WithGroup("onboarding"))
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
	go processTopic(topic, o.processOnboardingEvents)
}

func (o *Onboarding) processOnboardingEvents(ch <-chan mevents.Event) {
	logger.Infof("Starting to process onboarding events")
	for ev := range ch {
		var ve *onboarding.Event
		if err := json.Unmarshal(ev.Payload, ve); err != nil {
			ev.Nack()
			logger.Errorf("Error unmarshalling onboarding event: $s", err)
			continue
		}
		switch ve.Type {
		case "newSignup":
			if err := o.processSignup(ve.NewSignup); err != nil {
				ev.Nack()
				logger.Errorf("Error processing signup event %s", err)
				continue
			}
		default:
			logger.Infof("Unrecognised event %+v", ve)
		}
		ev.Ack()
	}
}

func (o *Onboarding) processSignup(ev *onboarding.NewSignupEvent) error {
	// add a promo credit to their balance
	_, err := o.balSvc.Increment(context.Background(), &balance.IncrementRequest{
		CustomerId: ev.Id,
		Delta:      o.promoCredit,
		Visible:    true,
		Reference:  o.promoMessage,
	})
	return err
}
