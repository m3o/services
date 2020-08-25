// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/namespace.proto

package namespace

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

// Api Endpoints for Namespace service

func NewNamespaceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Namespace service

type NamespaceService interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
	// Adds a new user to an existing namespace
	AddUser(ctx context.Context, in *AddUserRequest, opts ...client.CallOption) (*AddUserResponse, error)
	// Remove a user from a namespace
	RemoveUser(ctx context.Context, in *RemoveUserRequest, opts ...client.CallOption) (*RemoveUserResponse, error)
}

type namespaceService struct {
	c    client.Client
	name string
}

func NewNamespaceService(name string, c client.Client) NamespaceService {
	return &namespaceService{
		c:    c,
		name: name,
	}
}

func (c *namespaceService) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.name, "Namespace.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceService) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.name, "Namespace.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceService) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.name, "Namespace.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Namespace.List", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceService) AddUser(ctx context.Context, in *AddUserRequest, opts ...client.CallOption) (*AddUserResponse, error) {
	req := c.c.NewRequest(c.name, "Namespace.AddUser", in)
	out := new(AddUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceService) RemoveUser(ctx context.Context, in *RemoveUserRequest, opts ...client.CallOption) (*RemoveUserResponse, error) {
	req := c.c.NewRequest(c.name, "Namespace.RemoveUser", in)
	out := new(RemoveUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Namespace service

type NamespaceHandler interface {
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Read(context.Context, *ReadRequest, *ReadResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	List(context.Context, *ListRequest, *ListResponse) error
	// Adds a new user to an existing namespace
	AddUser(context.Context, *AddUserRequest, *AddUserResponse) error
	// Remove a user from a namespace
	RemoveUser(context.Context, *RemoveUserRequest, *RemoveUserResponse) error
}

func RegisterNamespaceHandler(s server.Server, hdlr NamespaceHandler, opts ...server.HandlerOption) error {
	type namespace interface {
		Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error
		Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error
		Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error
		List(ctx context.Context, in *ListRequest, out *ListResponse) error
		AddUser(ctx context.Context, in *AddUserRequest, out *AddUserResponse) error
		RemoveUser(ctx context.Context, in *RemoveUserRequest, out *RemoveUserResponse) error
	}
	type Namespace struct {
		namespace
	}
	h := &namespaceHandler{hdlr}
	return s.Handle(s.NewHandler(&Namespace{h}, opts...))
}

type namespaceHandler struct {
	NamespaceHandler
}

func (h *namespaceHandler) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.NamespaceHandler.Create(ctx, in, out)
}

func (h *namespaceHandler) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.NamespaceHandler.Read(ctx, in, out)
}

func (h *namespaceHandler) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.NamespaceHandler.Delete(ctx, in, out)
}

func (h *namespaceHandler) List(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.NamespaceHandler.List(ctx, in, out)
}

func (h *namespaceHandler) AddUser(ctx context.Context, in *AddUserRequest, out *AddUserResponse) error {
	return h.NamespaceHandler.AddUser(ctx, in, out)
}

func (h *namespaceHandler) RemoveUser(ctx context.Context, in *RemoveUserRequest, out *RemoveUserResponse) error {
	return h.NamespaceHandler.RemoveUser(ctx, in, out)
}
