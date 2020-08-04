package main

import (
	"github.com/m3o/services/payments/provider"
	"github.com/m3o/services/payments/provider/stripe/handler"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
)

func main() {
	// Setup the service
	srv := service.New(
		service.Name(provider.ServicePrefix + "stripe"),
	)

	// Register the provider
	srv.Handle(handler.New(srv))

	// Run the service
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running service: %v", err)
	}
}
