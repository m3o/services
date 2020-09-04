package handler

import (
	"docker.io/go-docker"
	"github.com/micro/go-micro/v3/logger"
)

// BuildHandler implements the build service interface:
type BuildHandler struct {
	baseImageURL  string
	buildImageURL string
	dockerClient  docker.ImageAPIClient
}

// New returns an initialised BuildHandler:
func New(baseImageURL, buildImageURL string) (*BuildHandler, error) {

	// Prepare a new Docker client:
	dockerClient, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}

	logger.Info("Prepared a new Build handler")

	return &BuildHandler{
		baseImageURL:  baseImageURL,
		buildImageURL: buildImageURL,
		dockerClient:  dockerClient,
	}, nil
}
