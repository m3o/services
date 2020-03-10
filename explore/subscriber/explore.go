package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	explore "explore/proto/explore"
)

type Explore struct{}

func (e *Explore) Handle(ctx context.Context, msg *explore.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *explore.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
