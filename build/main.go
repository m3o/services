package main

import (
	"github.com/m3o/services/build/builder"
	"github.com/m3o/services/build/handler"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/metrics/prometheus"
	"github.com/micro/go-micro/v3/metrics/wrapper"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

const (
	defaultBaseImageURL  = "golang:alpine"
	defaultBuildImageURL = "golang:1.14-buster"
)

func main() {

	// Prepare a Prometheus metrics reporter:
	reporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("Error preparing a Prometheus reporter: %v", err)
	}

	// New service:
	srv := service.New(
		service.Name("build"),
		service.Version("latest"),
		service.WrapHandler(wrapper.New(reporter).HandlerFunc),
	)

	// Get some config:
	baseImageURL := config.Get("micro", "build", "baseImageURL").String(defaultBaseImageURL)
	buildImageURL := config.Get("micro", "build", "buildImageURL").String(defaultBuildImageURL)

	// Prepare a docker builder for the handler to use:
	dockerBuilder, err := builder.NewShellBuilder(reporter, baseImageURL, buildImageURL)
	if err != nil {
		log.Fatalf("Error preparing a Docker builder: %v", err)
	}

	// Register handler:
	if err := srv.Handle(handler.New(dockerBuilder)); err != nil {
		log.Fatalf("Error registering handler: %v", err)
	}

	// Run the service:
	if err := srv.Run(); err != nil {
		log.Fatalf("Error running the service: %v", err)
	}
}
