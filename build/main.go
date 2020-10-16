package main

import (
	"github.com/m3o/services/build/handler"
	pb "github.com/micro/micro/v3/proto/build"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/build/golang"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("build"),
		service.Version("latest"),
	)

	// Setup the build
	build, err := golang.NewBuilder()
	if err != nil {
		logger.Fatalf("Error setting up golang build: %v", err)
	}

	// Register the handler
	pb.RegisterBuildHandler(srv.Server(), &handler.Handler{Builder: build})

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
