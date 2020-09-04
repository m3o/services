// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/build/build.proto

package build

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v3/api"
	client "github.com/micro/go-micro/v3/client"
	server "github.com/micro/go-micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Build service

func NewBuildEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Build service

type BuildService interface {
	BuildImageFromGitRepo(ctx context.Context, in *BuildImageFromGitRepoRequest, opts ...client.CallOption) (*BuildImageFromGitRepoResponse, error)
}

type buildService struct {
	c    client.Client
	name string
}

func NewBuildService(name string, c client.Client) BuildService {
	return &buildService{
		c:    c,
		name: name,
	}
}

func (c *buildService) BuildImageFromGitRepo(ctx context.Context, in *BuildImageFromGitRepoRequest, opts ...client.CallOption) (*BuildImageFromGitRepoResponse, error) {
	req := c.c.NewRequest(c.name, "Build.BuildImageFromGitRepo", in)
	out := new(BuildImageFromGitRepoResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Build service

type BuildHandler interface {
	BuildImageFromGitRepo(context.Context, *BuildImageFromGitRepoRequest, *BuildImageFromGitRepoResponse) error
}

func RegisterBuildHandler(s server.Server, hdlr BuildHandler, opts ...server.HandlerOption) error {
	type build interface {
		BuildImageFromGitRepo(ctx context.Context, in *BuildImageFromGitRepoRequest, out *BuildImageFromGitRepoResponse) error
	}
	type Build struct {
		build
	}
	h := &buildHandler{hdlr}
	return s.Handle(s.NewHandler(&Build{h}, opts...))
}

type buildHandler struct {
	BuildHandler
}

func (h *buildHandler) BuildImageFromGitRepo(ctx context.Context, in *BuildImageFromGitRepoRequest, out *BuildImageFromGitRepoResponse) error {
	return h.BuildHandler.BuildImageFromGitRepo(ctx, in, out)
}
