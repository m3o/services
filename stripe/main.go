package main

import (
	"github.com/m3o/services/stripe/handler"
	pb "github.com/m3o/services/stripe/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("stripe"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterStripeHandler(srv.Server(), new(handler.Stripe))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
