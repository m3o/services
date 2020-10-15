package main

import (
	"github.com/m3o/services/tests/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	service := service.New(
		service.Name("tests"),
	)

	service.Handle(new(handler.Tests))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
