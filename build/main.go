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
	defaultBaseImageURL     = "rg.fr-par.scw.cloud/build/base:latest"
	defaultBuildImageURL    = "rg.fr-par.scw.cloud/build/build:latest"
	defaultBuildRegistryURL = "rg.fr-par.scw.cloud/build"
	defaultDockerCommand    = "docker"
	defaultRegistryUsername = "nologin"
	defaultRegistryPassword = "changeme"
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
	builderConfig := &builder.Config{
		BaseImageURL:     config.Get("micro", "build", "baseimageurl").String(defaultBaseImageURL),
		BuildImageURL:    config.Get("micro", "build", "buildimageurl").String(defaultBuildImageURL),
		BuildRegistryURL: config.Get("micro", "build", "buildregistryurl").String(defaultBuildRegistryURL),
		DockerCommand:    config.Get("micro", "build", "dockerCommand").String(defaultDockerCommand),
		RegistryUsername: config.Get("micro", "build", "registryusername").String(defaultRegistryUsername),
		RegistryPassword: config.Get("micro", "build", "registrypassword").String(defaultRegistryPassword),
	}

	// Prepare a docker builder for the handler to use:
	dockerBuilder, err := builder.New(reporter, builderConfig)
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
