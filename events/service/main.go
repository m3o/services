package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/events/service/handler"
	pb "github.com/micro/services/events/service/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.events"),
		micro.Version("latest"),
	)

	service.Init()

	pb.RegisterEventsHandler(service.Server(), handler.New(service))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
