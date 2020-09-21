package main

import (
	"github.com/m3o/services/builder/handler"
	"github.com/micro/go-micro/v3/runtime/builder/golang"
	pb "github.com/micro/micro/v3/proto/runtime/builder"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("builder"),
		service.Version("latest"),
	)

	// Setup the builder
	builder, err := golang.NewBuilder()
	if err != nil {
		logger.Fatalf("Error setting up golang builder: %v", err)
	}

	// Register the handler
	pb.RegisterBuilderHandler(srv.Server(), &handler.Handler{Builder: builder})

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
