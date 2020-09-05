package builder

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"text/template"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/micro/go-micro/v3/logger"
)

// dockerfileTemplateRaw gets rendered to a buildable Dockerfile:
var dockerfileTemplateRaw = `# Build a service binary in a GoLang container:
FROM {{.BuildImage}} AS build
RUN go get {{.SourceGitRepo}}
RUN go build -o /service {{.SourceGitRepo}}

# Copy the service binary into a lean Alpine container:
FROM {{.BaseImage}}
COPY --from=build /service /service
CMD ["/service"]
`

// Docker implements the Builder interface:
type Docker struct {
	baseImageURL  string
	buildImageURL string
	dockerClient  docker.ImageAPIClient
}

// New returns a configured Docker builder:
func New(baseImageURL, buildImageURL string, dockerClient docker.ImageAPIClient) *Docker {
	return &Docker{
		baseImageURL:  baseImageURL,
		buildImageURL: buildImageURL,
		dockerClient:  dockerClient,
	}
}

// Build actually builds a Docker image:
func (b *Docker) Build(sourceGitRepo, targetImageTag string) error {

	// Prepare a Docker build context:
	dockerBuildContext, err := b.prepareBuildContext(sourceGitRepo)
	if err != nil {
		return err
	}

	// Prepare some ImageBuildOptions:
	imageBuildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Remove:     false,
		Tags:       []string{targetImageTag},
	}

	// Try to build an image:
	imageBuildResponse, err := b.dockerClient.ImageBuild(context.TODO(), dockerBuildContext, imageBuildOptions)
	if err != nil {
		logger.Warnf("Error building image: %v", err)
		return err
	}
	defer imageBuildResponse.Body.Close()

	// Log any message that came from Docker:
	var bodyBytes []byte
	imageBuildResponse.Body.Read(bodyBytes)
	logger.Debugf("ImageBuild response: %s", bodyBytes)

	return nil
}

// prepareBuildContext stuffs a Dockerfile into a TAR archive (I'm serious) which becomes the Docker "build context":
func (b *Docker) prepareBuildContext(sourceGitRepo string) (io.Reader, error) {

	// Render out the Dockerfile template:
	dockerfileContents, err := b.renderDockerFile(sourceGitRepo)
	if err != nil {
		return nil, err
	}

	// Now pack it into a TAR file (because this is how Docker rolls in 2020):
	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)

	// Create a header for the Dockerfile:
	dockerfileHeader := &tar.Header{
		Name: "Dockerfile",
		Mode: 0644,
		Size: int64(len(dockerfileContents)),
	}

	// Add the header to our tar archive:
	if err := tarWriter.WriteHeader(dockerfileHeader); err != nil {
		return nil, err
	}

	// And now write the Dockerfile contents:
	if _, err := tarWriter.Write(dockerfileContents); err != nil {
		return nil, err
	}

	// Close the tar archive:
	if err := tarWriter.Close(); err != nil {
		return nil, err
	}

	return tarBuffer, nil
}

// renderDockerFile uses parameters from config and from the RPC request to render the Dockerfile template:
func (b *Docker) renderDockerFile(sourceGitRepo string) ([]byte, error) {

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

	logger.Debugf("Generated Dockerfile: %s", buf.String())

	return buf.Bytes(), nil
}
