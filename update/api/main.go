package main

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/codec/proto"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/services/update/api/update"
)

type Update struct{}

func (u *Update) Info(ctx context.Context, req *proto.Message, rsp *proto.Message) error {
	// extract the data
	v := update.Get()
	b, _ := json.Marshal(v)
	rsp.Data = b
	return nil
}

func (u *Update) Event(ctx context.Context, req *proto.Message, rsp *proto.Message) error {
	return update.Event(ctx, req.Data)
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.update"),
	)
	service.Init()

	// register the handler
	service.Server().Handle(
		service.Server().NewHandler(new(Update)),
	)

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
