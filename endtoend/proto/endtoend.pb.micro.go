// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/endtoend.proto

package endtoend

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
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

// Api Endpoints for Endtoend service

func NewEndtoendEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Endtoend service

type EndtoendService interface {
	Mailin(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Check(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type endtoendService struct {
	c    client.Client
	name string
}

func NewEndtoendService(name string, c client.Client) EndtoendService {
	return &endtoendService{
		c:    c,
		name: name,
	}
}

func (c *endtoendService) Mailin(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Endtoend.Mailin", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endtoendService) Check(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Endtoend.Check", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Endtoend service

type EndtoendHandler interface {
	Mailin(context.Context, *Request, *Response) error
	Check(context.Context, *Request, *Response) error
}

func RegisterEndtoendHandler(s server.Server, hdlr EndtoendHandler, opts ...server.HandlerOption) error {
	type endtoend interface {
		Mailin(ctx context.Context, in *Request, out *Response) error
		Check(ctx context.Context, in *Request, out *Response) error
	}
	type Endtoend struct {
		endtoend
	}
	h := &endtoendHandler{hdlr}
	return s.Handle(s.NewHandler(&Endtoend{h}, opts...))
}

type endtoendHandler struct {
	EndtoendHandler
}

func (h *endtoendHandler) Mailin(ctx context.Context, in *Request, out *Response) error {
	return h.EndtoendHandler.Mailin(ctx, in, out)
}

func (h *endtoendHandler) Check(ctx context.Context, in *Request, out *Response) error {
	return h.EndtoendHandler.Check(ctx, in, out)
}
