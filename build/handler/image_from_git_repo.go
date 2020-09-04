package handler

import (
	"bytes"
	"context"
	"text/template"

	"docker.io/go-docker/api/types"
	"github.com/m3o/services/build/docker"
	pb "github.com/m3o/services/build/proto"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/logger"
)

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

	// Render a Dockerfile:
	dockerfile, err := h.renderDockerFile(request.SourceGitRepo)
	if err != nil {
		return err
	}

	// Use the Dockerfile to prepare build options:
	imageBuildOptions := types.ImageBuildOptions{
		Tags:       []string{request.TargetImageTag},
		Dockerfile: dockerfile,
	}

	// Try to build an image:
	imageBuildResponse, err := h.dockerClient.ImageBuild(context.TODO(), nil, imageBuildOptions)
	if err != nil {
		return err
	}

	response.BuiltImageTag = "cruft/cruft"
	response.BuiltImageURL = "cruft"

	logger.Infof("Built an image: %v", imageBuildResponse)

	return nil
}

func (h *BuildHandler) renderDockerFile(source string) (string, error) {

	// Prepare a build with the metadata we need to render a Dockerfile template:
	build := docker.Build{
		BaseImage:  h.baseImageURL,
		BuildImage: h.buildImageURL,
		Source:     source,
	}

	// Create the template:
	dockerfileTemplate := template.New("Dockerfile")
	dockerfileTemplateParsed, err := dockerfileTemplate.Parse(docker.DockerfileTemplate)
	if err != nil {
		return "", err
	}

	// Render the template with our build:
	buf := new(bytes.Buffer)
	if err := dockerfileTemplateParsed.Execute(buf, build); err != nil {
		return "", err
	}

	return buf.String(), nil
}
