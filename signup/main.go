package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/services/signup/handler"
	"github.com/micro/services/signup/subscriber"

	paymentsproto "github.com/micro/services/payments/provider/proto"
	signup "github.com/micro/services/signup/proto/signup"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.signup"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	signup.RegisterSignupHandler(service.Server(), handler.NewSignup(
		paymentsproto.NewProviderService("go.micro.service.payment.stripe", service.Options().Client),
		service.Options().Store,
		service.Options().Config,
		service.Options().Auth,
	))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.signup", service.Server(), new(subscriber.Signup))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
