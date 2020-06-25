package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/services/onboarding/handler"
	"github.com/micro/services/onboarding/subscriber"

	onboarding "github.com/micro/services/onboarding/proto/onboarding"
	paymentsproto "github.com/micro/services/payments/provider/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.onboarding"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	onboarding.RegisterOnboardingHandler(service.Server(), handler.NewOnboarding(
		paymentsproto.NewProviderService("go.micro.payment.service.stripe", service.Options().Client),
		service.Options().Store,
		service.Options().Config,
		service.Options().Auth,
	))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.onboarding", service.Server(), new(subscriber.Onboarding))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
