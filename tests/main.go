package main

import (
	"github.com/m3o/services/tests/handler"
	pb "github.com/m3o/services/tests/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.tests"),
	)

	service.Init()

	pb.RegisterTestsHandler(service.Server(), new(handler.Tests))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
