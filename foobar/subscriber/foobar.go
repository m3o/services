package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	foobar "foobar/proto/foobar"
)

type Foobar struct{}

func (e *Foobar) Handle(ctx context.Context, msg *foobar.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *foobar.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
