// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/home/home.proto

package go_micro_api_home

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

// Client API for Home service

type HomeService interface {
	ReadUser(ctx context.Context, in *ReadUserRequest, opts ...client.CallOption) (*ReadUserResponse, error)
}

type homeService struct {
	c    client.Client
	name string
}

func NewHomeService(name string, c client.Client) HomeService {
	return &homeService{
		c:    c,
		name: name,
	}
}

func (c *homeService) ReadUser(ctx context.Context, in *ReadUserRequest, opts ...client.CallOption) (*ReadUserResponse, error) {
	req := c.c.NewRequest(c.name, "Home.ReadUser", in)
	out := new(ReadUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Home service

type HomeHandler interface {
	ReadUser(context.Context, *ReadUserRequest, *ReadUserResponse) error
}

func RegisterHomeHandler(s server.Server, hdlr HomeHandler, opts ...server.HandlerOption) error {
	type home interface {
		ReadUser(ctx context.Context, in *ReadUserRequest, out *ReadUserResponse) error
	}
	type Home struct {
		home
	}
	h := &homeHandler{hdlr}
	return s.Handle(s.NewHandler(&Home{h}, opts...))
}

type homeHandler struct {
	HomeHandler
}

func (h *homeHandler) ReadUser(ctx context.Context, in *ReadUserRequest, out *ReadUserResponse) error {
	return h.HomeHandler.ReadUser(ctx, in, out)
}
