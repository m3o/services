// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/usage.proto

package usage

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

// Api Endpoints for Usage service

func NewUsageEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Usage service

type UsageService interface {
	ListSamples(ctx context.Context, in *ListSamplesRequest, opts ...client.CallOption) (*ListSamplesResponse, error)
}

type usageService struct {
	c    client.Client
	name string
}

func NewUsageService(name string, c client.Client) UsageService {
	return &usageService{
		c:    c,
		name: name,
	}
}

func (c *usageService) ListSamples(ctx context.Context, in *ListSamplesRequest, opts ...client.CallOption) (*ListSamplesResponse, error) {
	req := c.c.NewRequest(c.name, "Usage.ListSamples", in)
	out := new(ListSamplesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Usage service

type UsageHandler interface {
	ListSamples(context.Context, *ListSamplesRequest, *ListSamplesResponse) error
}

func RegisterUsageHandler(s server.Server, hdlr UsageHandler, opts ...server.HandlerOption) error {
	type usage interface {
		ListSamples(ctx context.Context, in *ListSamplesRequest, out *ListSamplesResponse) error
	}
	type Usage struct {
		usage
	}
	h := &usageHandler{hdlr}
	return s.Handle(s.NewHandler(&Usage{h}, opts...))
}

type usageHandler struct {
	UsageHandler
}

func (h *usageHandler) ListSamples(ctx context.Context, in *ListSamplesRequest, out *ListSamplesResponse) error {
	return h.UsageHandler.ListSamples(ctx, in, out)
}
