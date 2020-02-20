package subscriber

import (
	"context"
	"github.com/micro/go-micro/v2/util/log"

	peartree "peartree/proto/peartree"
)

type Peartree struct{}

func (e *Peartree) Handle(ctx context.Context, msg *peartree.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *peartree.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
