package subscriber

import (
	"context"
	"github.com/micro/go-micro/v2/util/log"

	vtolstov "vtolstov/proto/vtolstov"
)

type Vtolstov struct{}

func (e *Vtolstov) Handle(ctx context.Context, msg *vtolstov.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *vtolstov.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
