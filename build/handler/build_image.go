package handler

import (
	"context"
	"fmt"

	pb "github.com/m3o/services/build/proto/build"
	"github.com/micro/go-micro/v3/client"
)

// BuildImageFromGitRepo builds a service from source (a git repo), pushes to a Docker registry, and returns the image URL:
func (h *BuildHandler) BuildImageFromGitRepo(ctx context.Context, in *pb.BuildImageFromGitRepoRequest, opts ...client.CallOption) (*pb.BuildImageFromGitRepoResponse, error) {
	return nil, fmt.Errorf("BuildImageFromGitRepo is unimplemented")
}
