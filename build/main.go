package main

import (
	"github.com/m3o/services/build/handler"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

const (
	defaultBaseImageURL = "alpine:latest"
)

func main() {

	// New Service:
	srv := service.New(
		service.Name("build"),
	)

	// Get some config:
	baseImageURL := config.Get("micro", "build", "baseimageurl").String(defaultBaseImageURL)

	// Prepare a new handler:
	buildHandler, err := handler.New(baseImageURL)
	if err != nil {
		log.Fatalf("Error preparing a Build handler: %v", err)
	}

	// Register Handler
	srv.Handle(buildHandler)

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running the service: %v", err)
	}
}
