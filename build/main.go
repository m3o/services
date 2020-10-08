package main

import (
	"github.com/m3o/services/build/handler"
	pb "github.com/micro/micro/v3/proto/runtime/build"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime/builder/golang"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("build"),
		service.Version("latest"),
	)

	// Setup the builder
	builder, err := golang.NewBuilder()
	if err != nil {
		logger.Fatalf("Error setting up golang builder: %v", err)
	}

	// Register the handler
	pb.RegisterBuildHandler(srv.Server(), &handler.Handler{Builder: builder})

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
