// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: provider/proto/provider.proto

package go_micro_srv_provider

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
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
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Provider service

type ProviderService interface {
	Test(ctx context.Context, in *TestRequest, opts ...client.CallOption) (*TestResponse, error)
}

type providerService struct {
	c    client.Client
	name string
}

func NewProviderService(name string, c client.Client) ProviderService {
	return &providerService{
		c:    c,
		name: name,
	}
}

func (c *providerService) Test(ctx context.Context, in *TestRequest, opts ...client.CallOption) (*TestResponse, error) {
	req := c.c.NewRequest(c.name, "Provider.Test", in)
	out := new(TestResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Provider service

type ProviderHandler interface {
	Test(context.Context, *TestRequest, *TestResponse) error
}

func RegisterProviderHandler(s server.Server, hdlr ProviderHandler, opts ...server.HandlerOption) error {
	type provider interface {
		Test(ctx context.Context, in *TestRequest, out *TestResponse) error
	}
	type Provider struct {
		provider
	}
	h := &providerHandler{hdlr}
	return s.Handle(s.NewHandler(&Provider{h}, opts...))
}

type providerHandler struct {
	ProviderHandler
}

func (h *providerHandler) Test(ctx context.Context, in *TestRequest, out *TestResponse) error {
	return h.ProviderHandler.Test(ctx, in, out)
}
