package main

import (
	paymentsproto "github.com/m3o/services/payments/provider/proto"
	"github.com/m3o/services/subscriptions/handler"
	"github.com/micro/micro/v3/service"
	eventsproto "github.com/micro/micro/v3/service/events/proto"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("subscriptions"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(handler.New(
		paymentsproto.NewProviderService("payment.stripe", srv.Client()),
		eventsproto.NewStreamService("events", srv.Client()),
	))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
