package main

import (
	"github.com/m3o/services/customers/handler"
	"github.com/micro/micro/v3/service"
	eventsproto "github.com/micro/micro/v3/service/events/proto"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("customers"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(handler.New(eventsproto.NewStreamService("events", srv.Client())))
	handler.ConsumeEvents()
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
