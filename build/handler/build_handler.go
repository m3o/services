package handler

import (
	"context"
	"time"

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

	// Validate GitCommit:
	if request.GetGitCommit() == "" {
		return errors.BadRequest("request.validation", "GitCommit is required")
	}

	// Validate GitRepo:
	if request.GetGitRepo() == "" {
		return errors.BadRequest("request.validation", "GitRepo is required")
	}

	// Validate GitImageTag:
	if request.GetImageTag() == "" {
		return errors.BadRequest("request.validation", "ImageTag is required")
	}

	// Run the build and push in the background:
	go func() {

		// Build:
		logger.Debugf("Requesting image build (%s)", request.ImageTag)
		buildBeginTime := time.Now()
		buildOutput, err := b.builder.Build(request.GitRepo, request.GitCommit, request.ImageTag)
		logger.Debugf("Build output (%s): %s", request.ImageTag, buildOutput)
		if err != nil {
			logger.Errorf("Error requesting image build (%s): %v", request.ImageTag, err)
			return
		}
		logger.Infof("Image (%s) has been built in %s", request.ImageTag, time.Since(buildBeginTime).String())

		// Push:
		logger.Debugf("Requesting image push (%s)", request.ImageTag)
		pushBeginTime := time.Now()
		pushOutput, err := b.builder.Push(request.ImageTag)
		logger.Debugf("Push output (%s): %s", request.ImageTag, pushOutput)
		if err != nil {
			logger.Errorf("Error requesting image push (%s): %v", request.ImageTag, err)
			return
		}
		logger.Infof("Image (%s) has been pushed in %s", request.ImageTag, time.Since(pushBeginTime).String())
	}()

	logger.Infof("Build requested (%s)", request.ImageTag)

	return nil
}

// StreamImage builds a service from source (a git repo), pushes to a Docker registry, and returns the image URL:
func (b *Build) StreamImage(ctx context.Context, request *pb.CreateImageRequest, stream pb.Build_StreamImageStream) error {

	// Validate GitCommit:
	if request.GetGitCommit() == "" {
		stream.Send(&pb.CreateImageResponse{
			Status: pb.Status_REQUEST_FAILED,
			Error:  "request.validation",
			Output: "GitCommit is required",
		})
		return errors.BadRequest("request.validation", "GitCommit is required")
	}

	// Validate GitRepo:
	if request.GetGitRepo() == "" {
		stream.Send(&pb.CreateImageResponse{
			Status: pb.Status_REQUEST_FAILED,
			Error:  "request.validation",
			Output: "GitRepo is required",
		})
		return errors.BadRequest("request.validation", "GitRepo is required")
	}

	// Validate GitImageTag:
	if request.GetImageTag() == "" {
		stream.Send(&pb.CreateImageResponse{
			Status: pb.Status_REQUEST_FAILED,
			Error:  "request.validation",
			Output: "ImageTag is required",
		})
		return errors.BadRequest("request.validation", "ImageTag is required")
	}

	// Build the image:
	logger.Debugf("Requesting image build (%s)", request.ImageTag)
	stream.Send(&pb.CreateImageResponse{Status: pb.Status_BUILDING})
	buildBeginTime := time.Now()
	buildOutput, err := b.builder.Build(request.GitRepo, request.GitCommit, request.ImageTag)
	logger.Debugf("Build output (%s): %s", request.ImageTag, buildOutput)
	if err != nil {
		stream.Send(&pb.CreateImageResponse{
			Status: pb.Status_BUILD_FAILED,
			Error:  err.Error(),
			Output: buildOutput,
		})
		return errors.InternalServerError("image.build", "Error requesting image build: %v", err)
	}
	logger.Infof("Image (%s) has been built in %s", request.ImageTag, time.Since(buildBeginTime).String())

	// Push the image:
	logger.Debugf("Requesting image push (%s)", request.ImageTag)
	stream.Send(&pb.CreateImageResponse{Status: pb.Status_PUSHING})
	pushBeginTime := time.Now()
	pushOutput, err := b.builder.Push(request.ImageTag)
	logger.Debugf("Push output (%s): %s", request.ImageTag, pushOutput)
	if err != nil {
		stream.Send(&pb.CreateImageResponse{
			Status: pb.Status_PUSHING_FAILED,
			Error:  err.Error(),
			Output: pushOutput,
		})
		return errors.InternalServerError("image.push", "Error requesting image push: %v", err)
	}
	logger.Infof("Image (%s) has been pushed in %s", request.ImageTag, time.Since(pushBeginTime).String())

	// Complete:
	stream.Send(&pb.CreateImageResponse{
		Status: pb.Status_COMPLETE,
		Output: pushOutput,
	})

	return nil
}
