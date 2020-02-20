package subscriber

import (
	"context"
	"github.com/micro/go-micro/v2/util/log"

	cruft "cruft/proto/cruft"
)

type Cruft struct{}

func (e *Cruft) Handle(ctx context.Context, msg *cruft.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *cruft.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
