package main

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/services/update/api/update"
)

type Update struct{}

func (u *Update) Info(ctx context.Context, req json.RawMessage, rsp *json.RawMessage) error {
	// extract the data
	v := update.Get()
	b, _ := json.Marshal(v)
	*rsp = json.RawMessage(b)
	return nil
}

func (u *Update) Event(ctx context.Context, req json.RawMessage, rsp *json.RawMessage) error {
	return update.Event(ctx, req)
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.update"),
	)
	service.Init()

	// setup update
	update.Init()

	// register the handler
	service.Server().Handle(
		service.Server().NewHandler(new(Update)),
	)

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
