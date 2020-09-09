package builder

import (
	"bytes"
	"html/template"
	"io"

	"github.com/micro/go-micro/v3/logger"
)

// Config has the attributes we need to know how to authenticate, build, and push:
type Config struct {
	BaseImageURL     string
	BuildImageURL    string
	BuildRegistryURL string
	DockerCommand    string
	RegistryUsername string
	RegistryPassword string
}

// build has what we need to render the Dockerfile template:
type build struct {
	BaseImage       string
	BuildImage      string
	SourceGitCommit string
	SourceGitRepo   string
}

// renderDockerFile uses parameters from config and from the RPC request to render the Dockerfile template:
func (b *build) renderDockerFile() (io.Reader, error) {

	// Create the template:
	dockerfileTemplate := template.New("Dockerfile")
	dockerfileTemplateParsed, err := dockerfileTemplate.Parse(dockerfileTemplateRaw)
	if err != nil {
		return nil, err
	}

	// Render the template with our build:
	buf := new(bytes.Buffer)
	if err := dockerfileTemplateParsed.Execute(buf, b); err != nil {
		return nil, err
	}

	logger.Debugf("Generated Dockerfile: %s", buf)

	return buf, nil
}
