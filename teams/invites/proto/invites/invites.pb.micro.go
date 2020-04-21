// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: github.com/micro/services/teams/invites/proto/invites/invites.proto

package go_micro_service_teams_invites

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
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
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Invites service

func NewInvitesEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Invites service

type InvitesService interface {
	// Generate an invite to a user. An email will be sent to this
	// user containing a code which is valid for 24 hours.
	Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error)
	// Verify is called to ensure a code is valid, e.g has not expired.
	// This rpc should be called when the user opens the link in their
	// email before they create a profile.
	Verify(ctx context.Context, in *VerifyRequest, opts ...client.CallOption) (*VerifyResponse, error)
	// Redeem is used called after user completes signup and has an account.
	// Now they have an account we can redeem the invite and add the user
	// to the team. Once this rpc is called, the invite code can no longer
	// be used. The email address used when generating the invite must match
	// the email of the user redeeming the token.
	Redeem(ctx context.Context, in *RedeemRequest, opts ...client.CallOption) (*RedeemResponse, error)
}

type invitesService struct {
	c    client.Client
	name string
}

func NewInvitesService(name string, c client.Client) InvitesService {
	return &invitesService{
		c:    c,
		name: name,
	}
}

func (c *invitesService) Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error) {
	req := c.c.NewRequest(c.name, "Invites.Generate", in)
	out := new(GenerateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *invitesService) Verify(ctx context.Context, in *VerifyRequest, opts ...client.CallOption) (*VerifyResponse, error) {
	req := c.c.NewRequest(c.name, "Invites.Verify", in)
	out := new(VerifyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *invitesService) Redeem(ctx context.Context, in *RedeemRequest, opts ...client.CallOption) (*RedeemResponse, error) {
	req := c.c.NewRequest(c.name, "Invites.Redeem", in)
	out := new(RedeemResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Invites service

type InvitesHandler interface {
	// Generate an invite to a user. An email will be sent to this
	// user containing a code which is valid for 24 hours.
	Generate(context.Context, *GenerateRequest, *GenerateResponse) error
	// Verify is called to ensure a code is valid, e.g has not expired.
	// This rpc should be called when the user opens the link in their
	// email before they create a profile.
	Verify(context.Context, *VerifyRequest, *VerifyResponse) error
	// Redeem is used called after user completes signup and has an account.
	// Now they have an account we can redeem the invite and add the user
	// to the team. Once this rpc is called, the invite code can no longer
	// be used. The email address used when generating the invite must match
	// the email of the user redeeming the token.
	Redeem(context.Context, *RedeemRequest, *RedeemResponse) error
}

func RegisterInvitesHandler(s server.Server, hdlr InvitesHandler, opts ...server.HandlerOption) error {
	type invites interface {
		Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error
		Verify(ctx context.Context, in *VerifyRequest, out *VerifyResponse) error
		Redeem(ctx context.Context, in *RedeemRequest, out *RedeemResponse) error
	}
	type Invites struct {
		invites
	}
	h := &invitesHandler{hdlr}
	return s.Handle(s.NewHandler(&Invites{h}, opts...))
}

type invitesHandler struct {
	InvitesHandler
}

func (h *invitesHandler) Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error {
	return h.InvitesHandler.Generate(ctx, in, out)
}

func (h *invitesHandler) Verify(ctx context.Context, in *VerifyRequest, out *VerifyResponse) error {
	return h.InvitesHandler.Verify(ctx, in, out)
}

func (h *invitesHandler) Redeem(ctx context.Context, in *RedeemRequest, out *RedeemResponse) error {
	return h.InvitesHandler.Redeem(ctx, in, out)
}
