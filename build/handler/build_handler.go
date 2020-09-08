package handler

import (
	"context"

	"github.com/m3o/services/build/builder"
	pb "github.com/m3o/services/build/proto"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/logger"
)

// Build implements the build service interface:
type Build struct {
	builder builder.Builder
}

// New returns an initialised build handler:
func New(builder builder.Builder) *Build {

	logger.Info("Prepared a new Build handler")

	return &Build{
		builder: builder,
	}
}

// CreateImage builds a service from source (a git repo), pushes to a Docker registry, and returns the image URL:
func (b *Build) CreateImage(ctx context.Context, request *pb.CreateImageRequest, response *pb.CreateImageResponse) error {

	if request.GetGitCommit() == "" {
		return errors.BadRequest("request.validation", "GitCommit is required")
	}

	if request.GetGitRepo() == "" {
		return errors.BadRequest("request.validation", "GitRepo is required")
	}

	if request.GetImageTag() == "" {
		return errors.BadRequest("request.validation", "ImageTag is required")
	}

	if err := b.builder.Build(request.GitRepo, request.GitCommit, request.ImageTag); err != nil {
		return errors.InternalServerError("docker.build", "Error request Docker image build: %v", err)
	}

	logger.Infof("Build requested (%s)", request.ImageTag)

	return nil
}
