package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"

	"github.com/micro/services/payments/provider"
	pb "github.com/micro/services/payments/provider/proto"
	"github.com/micro/services/payments/provider/stripe/handler"
)

func main() {
	// Setup the service
	service := micro.NewService(
		micro.Name(provider.ServicePrefix+"stripe"),
		micro.Version("latest"),
	)

	// Initialise the servicwe
	service.Init()

	// Register the provider
	h := handler.NewHandler(service)
	pb.RegisterProviderHandler(service.Server(), h)

	// Consume events from the users service
	micro.RegisterSubscriber("go.micro.service.users", service.Server(), h.HandleUserEvent, server.SubscriberQueue("queue.stripe"))

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatalf("Error running service: %v", err)
	}
}
