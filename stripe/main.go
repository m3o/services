package main

import (
	"github.com/m3o/services/stripe/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/stripe/stripe-go/v71"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("stripe"),
		service.Version("latest"),
	)
	srv.Server().Handle(
		srv.Server().NewHandler(
			handler.NewHandler(srv),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "Stripe.Webhook",
					Handler: "http",
					Method:  []string{"POST"},
				}),
		))

	v, err := config.Get("micro.stripe.stripe.api_key")
	if err != nil || v.String("") == "" {
		logger.Fatal("Failed to retrieve stripe key")
	}
	stripe.Key = v.String("")
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
