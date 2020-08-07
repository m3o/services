// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: invite/proto/invite.proto

package go_micro_service_invite

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

// Api Endpoints for Invite service

func NewInviteEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Invite service

type InviteService interface {
	// Create an invite
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	// Delete an invite
	Delete(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	// Validate an email
	Validate(ctx context.Context, in *ValidateRequest, opts ...client.CallOption) (*ValidateResponse, error)
}

type inviteService struct {
	c    client.Client
	name string
}

func NewInviteService(name string, c client.Client) InviteService {
	return &inviteService{
		c:    c,
		name: name,
	}
}

func (c *inviteService) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.name, "Invite.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inviteService) Delete(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.name, "Invite.Delete", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inviteService) Validate(ctx context.Context, in *ValidateRequest, opts ...client.CallOption) (*ValidateResponse, error) {
	req := c.c.NewRequest(c.name, "Invite.Validate", in)
	out := new(ValidateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Invite service

type InviteHandler interface {
	// Create an invite
	Create(context.Context, *CreateRequest, *CreateResponse) error
	// Delete an invite
	Delete(context.Context, *CreateRequest, *CreateResponse) error
	// Validate an email
	Validate(context.Context, *ValidateRequest, *ValidateResponse) error
}

func RegisterInviteHandler(s server.Server, hdlr InviteHandler, opts ...server.HandlerOption) error {
	type invite interface {
		Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error
		Delete(ctx context.Context, in *CreateRequest, out *CreateResponse) error
		Validate(ctx context.Context, in *ValidateRequest, out *ValidateResponse) error
	}
	type Invite struct {
		invite
	}
	h := &inviteHandler{hdlr}
	return s.Handle(s.NewHandler(&Invite{h}, opts...))
}

type inviteHandler struct {
	InviteHandler
}

func (h *inviteHandler) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.InviteHandler.Create(ctx, in, out)
}

func (h *inviteHandler) Delete(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.InviteHandler.Delete(ctx, in, out)
}

func (h *inviteHandler) Validate(ctx context.Context, in *ValidateRequest, out *ValidateResponse) error {
	return h.InviteHandler.Validate(ctx, in, out)
}
