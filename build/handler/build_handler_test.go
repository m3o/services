package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/m3o/services/build/builder/mock"
	build "github.com/m3o/services/build/proto"
	pb "github.com/m3o/services/build/proto"
	"github.com/stretchr/testify/assert"
)

type mockStream struct {
	responses chan *pb.CreateImageResponse
}

func (ms *mockStream) Context() context.Context {
	return context.TODO()
}

func (ms *mockStream) SendMsg(interface{}) error {
	return nil
}

func (ms *mockStream) RecvMsg(interface{}) error {
	return nil
}

func (ms *mockStream) Close() error {
	return nil
}

func (ms *mockStream) Send(response *pb.CreateImageResponse) error {
	ms.responses <- response
	return nil
}

func TestBuildHandlerAsync(t *testing.T) {

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
}

func TestBuildHandlerSync(t *testing.T) {
	var response *pb.CreateImageResponse

	// Prepare a test handler with a mock builder:
	mockBuilder := new(mock.MockBuilder)
	testHandler := New(mockBuilder)
	assert.NotNil(t, testHandler)

	// Prepare a mock stream to receive streamed messages on:
	testStream := &mockStream{
		responses: make(chan *pb.CreateImageResponse, 100),
	}

	// Test that the request validation picks up incomplete requests:
	assert.Error(t, testHandler.StreamImage(context.TODO(), &build.CreateImageRequest{}, testStream))
	response = <-testStream.responses
	assert.NotEmpty(t, response.Error)
	assert.Equal(t, "GitCommit is required", response.Output)

	assert.Error(t, testHandler.StreamImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master"}, testStream))
	response = <-testStream.responses
	assert.NotEmpty(t, response.Error)
	assert.Equal(t, "GitRepo is required", response.Output)

	assert.Error(t, testHandler.StreamImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master", GitRepo: "something"}, testStream))
	response = <-testStream.responses
	assert.NotEmpty(t, response.Error)
	assert.Equal(t, "ImageTag is required", response.Output)

	// Test that the request validation doesn't fail a valid request:
	assert.NoError(t, testHandler.StreamImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master", GitRepo: "something", ImageTag: "test/test:latest"}, testStream))
	response = <-testStream.responses
	assert.Equal(t, "BUILDING", response.Status.String())
	response = <-testStream.responses
	assert.Equal(t, "PUSHING", response.Status.String())
	response = <-testStream.responses
	assert.Equal(t, "COMPLETE", response.Status.String())

	// Test that the handler successfully marks a failed build request as an error:
	mockBuilder.BuildReturns("problem", fmt.Errorf("Something went wrong"))
	assert.Error(t, testHandler.StreamImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master", GitRepo: "something", ImageTag: "test/test:latest"}, testStream))
	response = <-testStream.responses
	assert.Equal(t, "BUILDING", response.Status.String())
	response = <-testStream.responses
	assert.Equal(t, "BUILD_FAILED", response.Status.String())

	// Test that the handler successfully marks a failed push request as an error:
	mockBuilder.BuildReturns("", nil)
	mockBuilder.PushReturns("problem", fmt.Errorf("Something went wrong"))
	assert.Error(t, testHandler.StreamImage(context.TODO(), &build.CreateImageRequest{GitCommit: "master", GitRepo: "something", ImageTag: "test/test:latest"}, testStream))
	response = <-testStream.responses
	assert.Equal(t, "BUILDING", response.Status.String())
	response = <-testStream.responses
	assert.Equal(t, "PUSHING", response.Status.String())
	response = <-testStream.responses
	assert.Equal(t, "PUSHING_FAILED", response.Status.String())
}
