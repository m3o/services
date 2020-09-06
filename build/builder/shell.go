package builder

import (
	"bytes"
	"io"
	"os/exec"
	"text/template"

	"github.com/micro/go-micro/v3/logger"
	"github.com/pkg/errors"
)

// ShellBuilder implements the Builder interface:
// It shells out to the "docker" command to do builds, image pulls, and image pushes.
type ShellBuilder struct {
	baseImageURL  string
	buildImageURL string
}

// NewShellBuilder returns a configured shell Docker builder:
func NewShellBuilder(baseImageURL, buildImageURL string) (*ShellBuilder, error) {

	// Pull the base image:
	logger.Debugf("Pulling base image (%s)", baseImageURL)
	if err := exec.Command("docker", "pull", baseImageURL).Run(); err != nil {
		return nil, errors.Wrapf(err, "Unable to pull base image (%s)", baseImageURL)
	}

	// Pull the build image:
	logger.Debugf("Pulling build image (%s)", buildImageURL)
	if err := exec.Command("docker", "pull", buildImageURL).Run(); err != nil {
		return nil, errors.Wrapf(err, "Unable to pull build image (%s)", buildImageURL)
	}

	return &ShellBuilder{
		baseImageURL:  baseImageURL,
		buildImageURL: buildImageURL,
	}, nil
}

// Build actually builds a Docker image:
func (b *ShellBuilder) Build(sourceGitRepo, targetImageTag string) error {

	// Render out the Dockerfile template:
	dockerfileContents, err := b.renderDockerFile(sourceGitRepo)
	if err != nil {
		return err
	}

	// Try to build an image (Dockerfile contents provided via StdIn):
	go func() {
		buildCommand := exec.Command("docker", "build", "--rm", "-t", targetImageTag, "-")
		buildCommand.Stdin = dockerfileContents
		if err := buildCommand.Run(); err != nil {
			logger.Errorf("Unable to build image (%s): %v", targetImageTag, err)
		}
		logger.Infof("Build finished (%s)", targetImageTag)
	}()

	return nil
}

// renderDockerFile uses parameters from config and from the RPC request to render the Dockerfile template:
func (b *ShellBuilder) renderDockerFile(sourceGitRepo string) (io.Reader, error) {

	// Prepare a build with the metadata we need to render a Dockerfile template:
	build := build{
		BaseImage:     b.baseImageURL,
		BuildImage:    b.buildImageURL,
		SourceGitRepo: sourceGitRepo,
	}

	// Create the template:
	dockerfileTemplate := template.New("Dockerfile")
	dockerfileTemplateParsed, err := dockerfileTemplate.Parse(dockerfileTemplateRaw)
	if err != nil {
		return nil, err
	}

	// Render the template with our build:
	buf := new(bytes.Buffer)
	if err := dockerfileTemplateParsed.Execute(buf, build); err != nil {
		return nil, err
	}

	logger.Debugf("Generated Dockerfile: %s", buf)

	return buf, nil
}
