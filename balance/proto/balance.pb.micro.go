// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/balance.proto

package balance

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

// Api Endpoints for Balance service

func NewBalanceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Balance service

type BalanceService interface {
	Increment(ctx context.Context, in *IncrementRequest, opts ...client.CallOption) (*IncrementResponse, error)
	Decrement(ctx context.Context, in *DecrementRequest, opts ...client.CallOption) (*DecrementResponse, error)
	Current(ctx context.Context, in *CurrentRequest, opts ...client.CallOption) (*CurrentResponse, error)
}

type balanceService struct {
	c    client.Client
	name string
}

func NewBalanceService(name string, c client.Client) BalanceService {
	return &balanceService{
		c:    c,
		name: name,
	}
}

func (c *balanceService) Increment(ctx context.Context, in *IncrementRequest, opts ...client.CallOption) (*IncrementResponse, error) {
	req := c.c.NewRequest(c.name, "Balance.Increment", in)
	out := new(IncrementResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balanceService) Decrement(ctx context.Context, in *DecrementRequest, opts ...client.CallOption) (*DecrementResponse, error) {
	req := c.c.NewRequest(c.name, "Balance.Decrement", in)
	out := new(DecrementResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *balanceService) Current(ctx context.Context, in *CurrentRequest, opts ...client.CallOption) (*CurrentResponse, error) {
	req := c.c.NewRequest(c.name, "Balance.Current", in)
	out := new(CurrentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Balance service

type BalanceHandler interface {
	Increment(context.Context, *IncrementRequest, *IncrementResponse) error
	Decrement(context.Context, *DecrementRequest, *DecrementResponse) error
	Current(context.Context, *CurrentRequest, *CurrentResponse) error
}

func RegisterBalanceHandler(s server.Server, hdlr BalanceHandler, opts ...server.HandlerOption) error {
	type balance interface {
		Increment(ctx context.Context, in *IncrementRequest, out *IncrementResponse) error
		Decrement(ctx context.Context, in *DecrementRequest, out *DecrementResponse) error
		Current(ctx context.Context, in *CurrentRequest, out *CurrentResponse) error
	}
	type Balance struct {
		balance
	}
	h := &balanceHandler{hdlr}
	return s.Handle(s.NewHandler(&Balance{h}, opts...))
}

type balanceHandler struct {
	BalanceHandler
}

func (h *balanceHandler) Increment(ctx context.Context, in *IncrementRequest, out *IncrementResponse) error {
	return h.BalanceHandler.Increment(ctx, in, out)
}

func (h *balanceHandler) Decrement(ctx context.Context, in *DecrementRequest, out *DecrementResponse) error {
	return h.BalanceHandler.Decrement(ctx, in, out)
}

func (h *balanceHandler) Current(ctx context.Context, in *CurrentRequest, out *CurrentResponse) error {
	return h.BalanceHandler.Current(ctx, in, out)
}