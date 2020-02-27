// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/ben/ben.proto

package go_micro_srv_ben

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

// Client API for Ben service

type BenService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Ben_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Ben_PingPongService, error)
}

type benService struct {
	c    client.Client
	name string
}

func NewBenService(name string, c client.Client) BenService {
	return &benService{
		c:    c,
		name: name,
	}
}

func (c *benService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Ben.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *benService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Ben_StreamService, error) {
	req := c.c.NewRequest(c.name, "Ben.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &benServiceStream{stream}, nil
}

type Ben_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type benServiceStream struct {
	stream client.Stream
}

func (x *benServiceStream) Close() error {
	return x.stream.Close()
}

func (x *benServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *benServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *benServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *benServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *benService) PingPong(ctx context.Context, opts ...client.CallOption) (Ben_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Ben.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &benServicePingPong{stream}, nil
}

type Ben_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type benServicePingPong struct {
	stream client.Stream
}

func (x *benServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *benServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *benServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *benServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *benServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *benServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Ben service

type BenHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Ben_StreamStream) error
	PingPong(context.Context, Ben_PingPongStream) error
}

func RegisterBenHandler(s server.Server, hdlr BenHandler, opts ...server.HandlerOption) error {
	type ben interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Ben struct {
		ben
	}
	h := &benHandler{hdlr}
	return s.Handle(s.NewHandler(&Ben{h}, opts...))
}

type benHandler struct {
	BenHandler
}

func (h *benHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.BenHandler.Call(ctx, in, out)
}

func (h *benHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.BenHandler.Stream(ctx, m, &benStreamStream{stream})
}

type Ben_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type benStreamStream struct {
	stream server.Stream
}

func (x *benStreamStream) Close() error {
	return x.stream.Close()
}

func (x *benStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *benStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *benStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *benStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *benHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.BenHandler.PingPong(ctx, &benPingPongStream{stream})
}

type Ben_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type benPingPongStream struct {
	stream server.Stream
}

func (x *benPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *benPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *benPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *benPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *benPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *benPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
