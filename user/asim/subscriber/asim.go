package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	asim "asim/proto/asim"
)

type Asim struct{}

func (e *Asim) Handle(ctx context.Context, msg *asim.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *asim.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
