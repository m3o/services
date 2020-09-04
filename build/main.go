package main

import (
	"github.com/m3o/services/build/handler"
	pb "github.com/m3o/services/build/proto"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

const (
	defaultBaseImageURL  = "alpine:latest"
	defaultBuildImageURL = "golang:1.14-alpine"
)

func main() {

	// New Service:
	srv := service.New(
		service.Name("build"),
		service.Version("latest"),
	)

	// Get some config:
	baseImageURL := config.Get("micro", "build", "baseImageURL").String(defaultBaseImageURL)
	buildImageURL := config.Get("micro", "build", "buildImageURL").String(defaultBuildImageURL)

	// Prepare a new handler:
	buildHandler, err := handler.New(baseImageURL, buildImageURL)
	if err != nil {
		log.Fatalf("Error preparing a Build handler: %v", err)
	}

	// Register Handler:
	if err := srv.Handle(buildHandler); err != nil {
		log.Fatalf("Error registering handler: %v", err)
	}

	// Register Handler with proto:
	if err := pb.RegisterBuildHandler(srv.Server(), buildHandler); err != nil {
		log.Fatalf("Error registering handler: %v", err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running the service: %v", err)
	}
}
