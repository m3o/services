package main

import (
	log "github.com/micro/go-micro/v3/logger"
	"github.com/m3o/services/payments/provider"
	pb "github.com/m3o/services/payments/provider/proto"
	"github.com/m3o/services/payments/provider/stripe/handler"
	"github.com/micro/micro/v3/service"
)

func main() {
	// Setup the service
	srv := service.New(
		service.Name(provider.ServicePrefix+"stripe"),
	)

	// Register the provider
	h := handler.NewHandler(srv)
	pb.RegisterProviderHandler(h)

	// Run the service
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running service: %v", err)
	}
}
