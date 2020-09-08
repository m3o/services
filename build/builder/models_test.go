package builder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	expectedDockerFile = `# Build a service binary in a GoLang container:
FROM golang:alpine AS build
RUN mkdir -p /tmp/build
RUN git clone https://github.com/micro/micro.git /tmp/build
RUN cd /tmp/build && git checkout development && go build -o /service

# Copy the service binary into a lean Alpine container:
FROM some.docker.registry/repo/test:latest AS service
COPY --from=build /service /service
CMD ["/service"]
`
)

func TestBuildModel(t *testing.T) {

	// Prepare a test build:
	testBuild := &build{
		BaseImage:       "some.docker.registry/repo/test:latest",
		BuildImage:      "golang:alpine",
		SourceGitCommit: "development",
		SourceGitRepo:   "https://github.com/micro/micro.git",
	}

	// Render a Dockerfile for the test build:
	testDockerFileContents, err := testBuild.renderDockerFile()
	assert.NoError(t, err)
	assert.NotNil(t, testDockerFileContents)

	// Check what we got:
	testDockerFileBytes, ok := testDockerFileContents.(*bytes.Buffer)
	assert.True(t, ok)
	assert.Greater(t, len(testDockerFileBytes.Bytes()), 0)
	assert.Contains(t, testDockerFileBytes.String(), expectedDockerFile)
}
