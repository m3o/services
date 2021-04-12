package main

import (
	"github.com/m3o/services/stripe/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("stripe"),
		service.Version("latest"),
	)
	srv.Server().Handle(
		srv.Server().NewHandler(
			new(handler.Stripe),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "Stripe.Webhook",
					Handler: "api",
					Method:  []string{"POST"},
				}),
		))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
