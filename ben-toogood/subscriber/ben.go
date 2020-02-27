package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	ben "ben-toogood/proto/ben"
)

type Ben struct{}

func (e *Ben) Handle(ctx context.Context, msg *ben.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *ben.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
