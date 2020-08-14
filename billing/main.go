package main

import (
	"github.com/m3o/services/billing/handler"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("billing"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(new(handler.Billing))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
