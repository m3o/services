package handler

import (
	"context"

	"docker.io/go-docker"
	"github.com/m3o/services/build/builder"
	pb "github.com/m3o/services/build/proto"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/logger"
)

// BuildHandler implements the build service interface:
type BuildHandler struct {
	baseImageURL  string
	buildImageURL string
	builder       builder.Builder
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
		builder:       builder.New(baseImageURL, buildImageURL, dockerClient),
		dockerClient:  dockerClient,
	}, nil
}

// ImageFromGitRepo builds a service from source (a git repo), pushes to a Docker registry, and returns the image URL:
func (h *BuildHandler) ImageFromGitRepo(ctx context.Context, request *pb.ImageFromGitRepoRequest, response *pb.ImageFromGitRepoResponse) error {

	if request.GetSourceGitRepo() == "" {
		return errors.BadRequest("request.validation", "SourceGitRepo is required")
	}

	if request.GetTargetDockerRegistry() == "" {
		return errors.BadRequest("request.validation", "TargetDockerRegistry is required")
	}

	if request.TargetImageTag == "" {
		return errors.BadRequest("request.validation", "TargetImageTag is required")
	}

	if err := h.builder.Build(request.SourceGitRepo, request.TargetImageTag); err != nil {
		return errors.InternalServerError("docker.build", "Error building Docker image: %v", err)
	}

	response.BuiltImageTag = "cruft/cruft"
	response.BuiltImageURL = "cruft"

	logger.Info("Built an image")

	return nil
}
