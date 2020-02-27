// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/notes/notes.proto

package go_micro_srv_notes

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

// Client API for Notes service

type NotesService interface {
	Create(ctx context.Context, in *CreateNoteRequest, opts ...client.CallOption) (*CreateNoteResponse, error)
	Update(ctx context.Context, opts ...client.CallOption) (Notes_UpdateService, error)
	Delete(ctx context.Context, in *DeleteNoteRequest, opts ...client.CallOption) (*DeleteNoteResponse, error)
	List(ctx context.Context, in *ListNotesRequest, opts ...client.CallOption) (*ListNotesResponse, error)
}

type notesService struct {
	c    client.Client
	name string
}

func NewNotesService(name string, c client.Client) NotesService {
	return &notesService{
		c:    c,
		name: name,
	}
}

func (c *notesService) Create(ctx context.Context, in *CreateNoteRequest, opts ...client.CallOption) (*CreateNoteResponse, error) {
	req := c.c.NewRequest(c.name, "Notes.Create", in)
	out := new(CreateNoteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesService) Update(ctx context.Context, opts ...client.CallOption) (Notes_UpdateService, error) {
	req := c.c.NewRequest(c.name, "Notes.Update", &UpdateNoteRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &notesServiceUpdate{stream}, nil
}

type Notes_UpdateService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*UpdateNoteRequest) error
}

type notesServiceUpdate struct {
	stream client.Stream
}

func (x *notesServiceUpdate) Close() error {
	return x.stream.Close()
}

func (x *notesServiceUpdate) Context() context.Context {
	return x.stream.Context()
}

func (x *notesServiceUpdate) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *notesServiceUpdate) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *notesServiceUpdate) Send(m *UpdateNoteRequest) error {
	return x.stream.Send(m)
}

func (c *notesService) Delete(ctx context.Context, in *DeleteNoteRequest, opts ...client.CallOption) (*DeleteNoteResponse, error) {
	req := c.c.NewRequest(c.name, "Notes.Delete", in)
	out := new(DeleteNoteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesService) List(ctx context.Context, in *ListNotesRequest, opts ...client.CallOption) (*ListNotesResponse, error) {
	req := c.c.NewRequest(c.name, "Notes.List", in)
	out := new(ListNotesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Notes service

type NotesHandler interface {
	Create(context.Context, *CreateNoteRequest, *CreateNoteResponse) error
	Update(context.Context, Notes_UpdateStream) error
	Delete(context.Context, *DeleteNoteRequest, *DeleteNoteResponse) error
	List(context.Context, *ListNotesRequest, *ListNotesResponse) error
}

func RegisterNotesHandler(s server.Server, hdlr NotesHandler, opts ...server.HandlerOption) error {
	type notes interface {
		Create(ctx context.Context, in *CreateNoteRequest, out *CreateNoteResponse) error
		Update(ctx context.Context, stream server.Stream) error
		Delete(ctx context.Context, in *DeleteNoteRequest, out *DeleteNoteResponse) error
		List(ctx context.Context, in *ListNotesRequest, out *ListNotesResponse) error
	}
	type Notes struct {
		notes
	}
	h := &notesHandler{hdlr}
	return s.Handle(s.NewHandler(&Notes{h}, opts...))
}

type notesHandler struct {
	NotesHandler
}

func (h *notesHandler) Create(ctx context.Context, in *CreateNoteRequest, out *CreateNoteResponse) error {
	return h.NotesHandler.Create(ctx, in, out)
}

func (h *notesHandler) Update(ctx context.Context, stream server.Stream) error {
	return h.NotesHandler.Update(ctx, &notesUpdateStream{stream})
}

type Notes_UpdateStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*UpdateNoteRequest, error)
}

type notesUpdateStream struct {
	stream server.Stream
}

func (x *notesUpdateStream) Close() error {
	return x.stream.Close()
}

func (x *notesUpdateStream) Context() context.Context {
	return x.stream.Context()
}

func (x *notesUpdateStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *notesUpdateStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *notesUpdateStream) Recv() (*UpdateNoteRequest, error) {
	m := new(UpdateNoteRequest)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *notesHandler) Delete(ctx context.Context, in *DeleteNoteRequest, out *DeleteNoteResponse) error {
	return h.NotesHandler.Delete(ctx, in, out)
}

func (h *notesHandler) List(ctx context.Context, in *ListNotesRequest, out *ListNotesResponse) error {
	return h.NotesHandler.List(ctx, in, out)
}
