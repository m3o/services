// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/build.proto

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
	CreateImage(ctx context.Context, in *CreateImageRequest, opts ...client.CallOption) (*CreateImageResponse, error)
	// Creates an image from Go source in a Git repo, streams back the status as it progresses:
	StreamImage(ctx context.Context, in *CreateImageRequest, opts ...client.CallOption) (Build_StreamImageService, error)
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

func (c *buildService) CreateImage(ctx context.Context, in *CreateImageRequest, opts ...client.CallOption) (*CreateImageResponse, error) {
	req := c.c.NewRequest(c.name, "Build.CreateImage", in)
	out := new(CreateImageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *buildService) StreamImage(ctx context.Context, in *CreateImageRequest, opts ...client.CallOption) (Build_StreamImageService, error) {
	req := c.c.NewRequest(c.name, "Build.StreamImage", &CreateImageRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &buildServiceStreamImage{stream}, nil
}

type Build_StreamImageService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*CreateImageResponse, error)
}

type buildServiceStreamImage struct {
	stream client.Stream
}

func (x *buildServiceStreamImage) Close() error {
	return x.stream.Close()
}

func (x *buildServiceStreamImage) Context() context.Context {
	return x.stream.Context()
}

func (x *buildServiceStreamImage) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *buildServiceStreamImage) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *buildServiceStreamImage) Recv() (*CreateImageResponse, error) {
	m := new(CreateImageResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Build service

type BuildHandler interface {
	CreateImage(context.Context, *CreateImageRequest, *CreateImageResponse) error
	// Creates an image from Go source in a Git repo, streams back the status as it progresses:
	StreamImage(context.Context, *CreateImageRequest, Build_StreamImageStream) error
}

func RegisterBuildHandler(s server.Server, hdlr BuildHandler, opts ...server.HandlerOption) error {
	type build interface {
		CreateImage(ctx context.Context, in *CreateImageRequest, out *CreateImageResponse) error
		StreamImage(ctx context.Context, stream server.Stream) error
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

func (h *buildHandler) CreateImage(ctx context.Context, in *CreateImageRequest, out *CreateImageResponse) error {
	return h.BuildHandler.CreateImage(ctx, in, out)
}

func (h *buildHandler) StreamImage(ctx context.Context, stream server.Stream) error {
	m := new(CreateImageRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.BuildHandler.StreamImage(ctx, m, &buildStreamImageStream{stream})
}

type Build_StreamImageStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*CreateImageResponse) error
}

type buildStreamImageStream struct {
	stream server.Stream
}

func (x *buildStreamImageStream) Close() error {
	return x.stream.Close()
}

func (x *buildStreamImageStream) Context() context.Context {
	return x.stream.Context()
}

func (x *buildStreamImageStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *buildStreamImageStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *buildStreamImageStream) Send(m *CreateImageResponse) error {
	return x.stream.Send(m)
}
