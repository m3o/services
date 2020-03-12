package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	serverless "serverless/proto/serverless"
)

type Serverless struct{}

func (e *Serverless) Handle(ctx context.Context, msg *serverless.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *serverless.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
