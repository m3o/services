package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	subscribe "subscribe/proto/subscribe"
)

type Subscribe struct{}

func (e *Subscribe) Handle(ctx context.Context, msg *subscribe.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *subscribe.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
