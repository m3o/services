package main

import (
	"github.com/m3o/services/tests/handler"
	pb "github.com/m3o/services/tests/proto"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
)

func main() {
	service := service.New(
		service.Name("go.micro.service.tests"),
	)

	pb.RegisterTestsHandler(new(handler.Tests))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
