package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/m3o/services/build/builder/mock"
	build "github.com/m3o/services/build/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuildHandler(t *testing.T) {

	// Prepare a test handler with a mock builder:
	mockBuilder := new(mock.MockBuilder)
	testHandler := New(mockBuilder)
	assert.NotNil(t, testHandler)

	// Test that the request validation picks up incomplete requests:
	rsp := new(build.CreateImageResponse)
	assert.Error(t, testHandler.CreateImage(context.TODO(), &build.CreateImageRequest{}, rsp))
	assert.Error(t, testHandler.CreateImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master"}, rsp))
	assert.Error(t, testHandler.CreateImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master", GitRepo: "something"}, rsp))

	// Test that the request validation doesn't fail a valid request:
	assert.NoError(t, testHandler.CreateImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master", GitRepo: "something", ImageTag: "test/test:latest"}, rsp))

	// Test that the handler successfully marks a failed build request as an error:
	mockBuilder.BuildReturns(fmt.Errorf("Something went wrong"))
	assert.Error(t, testHandler.CreateImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master", GitRepo: "something", ImageTag: "test/test:latest"}, rsp))
}
