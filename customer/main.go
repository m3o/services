package main

import (
	"github.com/m3o/services/customer/handler"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("customer"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(new(handler.Customer))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
