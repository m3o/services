package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	alert "github.com/m3o/services/alert/proto/alert"
)

type Alert struct{}

func (e *Alert) Handle(ctx context.Context, msg *alert.Event) error {
	log.Info("Handler Received message")
	return nil
}

func Handler(ctx context.Context, msg *alert.Event) error {
	log.Info("Function Received message")
	return nil
}
