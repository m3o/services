package builder

import (
	"testing"

	"github.com/micro/go-micro/v3/metrics/logging"
	"github.com/stretchr/testify/assert"
)

func TestCmdBuilder(t *testing.T) {

	// Get some config:
	builderConfig := &Config{
		BaseImageURL:     "alpine:base",
		BuildImageURL:    "golang:alpine",
		BuildRegistryURL: "some.host/registry",
		DockerCommand:    "echo",
		RegistryUsername: "username",
		RegistryPassword: "password",
	}

	// Prepare a docker builder for the handler to use:
	testBuilder, err := New(logging.New(), builderConfig)
	assert.NoError(t, err)
	assert.NotNil(t, testBuilder)

	// Check that the build command is being put together properly:
	buildOutput, err := testBuilder.Build("https://github.com/micro/micro.git", "master", "some.host/registry/micro:latest")
	assert.NoError(t, err)
	assert.Contains(t, buildOutput, "build --force-rm --rm -t some.host/registry/micro:latest -")

	// Check that the push command is being put together properly:
	pushOutput, err := testBuilder.Push("some.host/registry/micro:latest")
	assert.NoError(t, err)
	assert.Contains(t, pushOutput, "push some.host/registry/micro:latest")

	// Now set ourselves up for failure:
	testBuilder.config.DockerCommand = "something-nonexistent"

	// Try to build:
	buildOutput, err = testBuilder.Build("https://github.com/micro/micro.git", "master", "some.host/registry/micro:latest")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unable to build image (some.host/registry/micro:latest):")

	// Try to push:
	pushOutput, err = testBuilder.Push("some.host/registry/micro:latest")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unable to push image (some.host/registry/micro:latest):")

	// Try to login:
	err = testBuilder.dockerLogin(builderConfig.BuildRegistryURL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unable to login to Docker registry (some.host/registry):")
}
