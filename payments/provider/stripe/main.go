package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

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
	h, err := handler.NewHandler()
	if err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}
	pb.RegisterProviderHandler(service.Server(), &h)

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatalf("Error running service: %v", err)
	}
}
