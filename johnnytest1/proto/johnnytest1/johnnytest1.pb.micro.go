// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/johnnytest1/johnnytest1.proto

package go_micro_srv_johnnytest1

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

// Client API for Johnnytest1 service

type Johnnytest1Service interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Johnnytest1_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Johnnytest1_PingPongService, error)
}

type johnnytest1Service struct {
	c    client.Client
	name string
}

func NewJohnnytest1Service(name string, c client.Client) Johnnytest1Service {
	return &johnnytest1Service{
		c:    c,
		name: name,
	}
}

func (c *johnnytest1Service) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Johnnytest1.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *johnnytest1Service) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Johnnytest1_StreamService, error) {
	req := c.c.NewRequest(c.name, "Johnnytest1.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &johnnytest1ServiceStream{stream}, nil
}

type Johnnytest1_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type johnnytest1ServiceStream struct {
	stream client.Stream
}

func (x *johnnytest1ServiceStream) Close() error {
	return x.stream.Close()
}

func (x *johnnytest1ServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *johnnytest1ServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *johnnytest1ServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *johnnytest1ServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *johnnytest1Service) PingPong(ctx context.Context, opts ...client.CallOption) (Johnnytest1_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Johnnytest1.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &johnnytest1ServicePingPong{stream}, nil
}

type Johnnytest1_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type johnnytest1ServicePingPong struct {
	stream client.Stream
}

func (x *johnnytest1ServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *johnnytest1ServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *johnnytest1ServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *johnnytest1ServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *johnnytest1ServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *johnnytest1ServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Johnnytest1 service

type Johnnytest1Handler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Johnnytest1_StreamStream) error
	PingPong(context.Context, Johnnytest1_PingPongStream) error
}

func RegisterJohnnytest1Handler(s server.Server, hdlr Johnnytest1Handler, opts ...server.HandlerOption) error {
	type johnnytest1 interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Johnnytest1 struct {
		johnnytest1
	}
	h := &johnnytest1Handler{hdlr}
	return s.Handle(s.NewHandler(&Johnnytest1{h}, opts...))
}

type johnnytest1Handler struct {
	Johnnytest1Handler
}

func (h *johnnytest1Handler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.Johnnytest1Handler.Call(ctx, in, out)
}

func (h *johnnytest1Handler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.Johnnytest1Handler.Stream(ctx, m, &johnnytest1StreamStream{stream})
}

type Johnnytest1_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type johnnytest1StreamStream struct {
	stream server.Stream
}

func (x *johnnytest1StreamStream) Close() error {
	return x.stream.Close()
}

func (x *johnnytest1StreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *johnnytest1StreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *johnnytest1StreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *johnnytest1StreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *johnnytest1Handler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.Johnnytest1Handler.PingPong(ctx, &johnnytest1PingPongStream{stream})
}

type Johnnytest1_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type johnnytest1PingPongStream struct {
	stream server.Stream
}

func (x *johnnytest1PingPongStream) Close() error {
	return x.stream.Close()
}

func (x *johnnytest1PingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *johnnytest1PingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *johnnytest1PingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *johnnytest1PingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *johnnytest1PingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
