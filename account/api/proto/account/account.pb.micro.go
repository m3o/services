// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: github.com/micro/services/account/api/proto/account/account.proto

package go_micro_api_account

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

// Api Endpoints for Account service

func NewAccountEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Account service

type AccountService interface {
	ReadUser(ctx context.Context, in *ReadUserRequest, opts ...client.CallOption) (*ReadUserResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...client.CallOption) (*UpdateUserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*DeleteUserResponse, error)
	ListPlans(ctx context.Context, in *ListPlansRequest, opts ...client.CallOption) (*ListPlansResponse, error)
	CreateSubscription(ctx context.Context, in *CreateSubscriptionRequest, opts ...client.CallOption) (*CreateSubscriptionResponse, error)
	CreatePaymentMethod(ctx context.Context, in *CreatePaymentMethodRequest, opts ...client.CallOption) (*CreatePaymentMethodResponse, error)
	DeletePaymentMethod(ctx context.Context, in *DeletePaymentMethodRequest, opts ...client.CallOption) (*DeletePaymentMethodResponse, error)
	DefaultPaymentMethod(ctx context.Context, in *DefaultPaymentMethodRequest, opts ...client.CallOption) (*DefaultPaymentMethodResponse, error)
	Signup(ctx context.Context, in *SignupRequest, opts ...client.CallOption) (*SignupResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	Token(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error)
}

type accountService struct {
	c    client.Client
	name string
}

func NewAccountService(name string, c client.Client) AccountService {
	return &accountService{
		c:    c,
		name: name,
	}
}

func (c *accountService) ReadUser(ctx context.Context, in *ReadUserRequest, opts ...client.CallOption) (*ReadUserResponse, error) {
	req := c.c.NewRequest(c.name, "Account.ReadUser", in)
	out := new(ReadUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...client.CallOption) (*UpdateUserResponse, error) {
	req := c.c.NewRequest(c.name, "Account.UpdateUser", in)
	out := new(UpdateUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*DeleteUserResponse, error) {
	req := c.c.NewRequest(c.name, "Account.DeleteUser", in)
	out := new(DeleteUserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) ListPlans(ctx context.Context, in *ListPlansRequest, opts ...client.CallOption) (*ListPlansResponse, error) {
	req := c.c.NewRequest(c.name, "Account.ListPlans", in)
	out := new(ListPlansResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) CreateSubscription(ctx context.Context, in *CreateSubscriptionRequest, opts ...client.CallOption) (*CreateSubscriptionResponse, error) {
	req := c.c.NewRequest(c.name, "Account.CreateSubscription", in)
	out := new(CreateSubscriptionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) CreatePaymentMethod(ctx context.Context, in *CreatePaymentMethodRequest, opts ...client.CallOption) (*CreatePaymentMethodResponse, error) {
	req := c.c.NewRequest(c.name, "Account.CreatePaymentMethod", in)
	out := new(CreatePaymentMethodResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) DeletePaymentMethod(ctx context.Context, in *DeletePaymentMethodRequest, opts ...client.CallOption) (*DeletePaymentMethodResponse, error) {
	req := c.c.NewRequest(c.name, "Account.DeletePaymentMethod", in)
	out := new(DeletePaymentMethodResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) DefaultPaymentMethod(ctx context.Context, in *DefaultPaymentMethodRequest, opts ...client.CallOption) (*DefaultPaymentMethodResponse, error) {
	req := c.c.NewRequest(c.name, "Account.DefaultPaymentMethod", in)
	out := new(DefaultPaymentMethodResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) Signup(ctx context.Context, in *SignupRequest, opts ...client.CallOption) (*SignupResponse, error) {
	req := c.c.NewRequest(c.name, "Account.Signup", in)
	out := new(SignupResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "Account.Login", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountService) Token(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error) {
	req := c.c.NewRequest(c.name, "Account.Token", in)
	out := new(TokenResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Account service

type AccountHandler interface {
	ReadUser(context.Context, *ReadUserRequest, *ReadUserResponse) error
	UpdateUser(context.Context, *UpdateUserRequest, *UpdateUserResponse) error
	DeleteUser(context.Context, *DeleteUserRequest, *DeleteUserResponse) error
	ListPlans(context.Context, *ListPlansRequest, *ListPlansResponse) error
	CreateSubscription(context.Context, *CreateSubscriptionRequest, *CreateSubscriptionResponse) error
	CreatePaymentMethod(context.Context, *CreatePaymentMethodRequest, *CreatePaymentMethodResponse) error
	DeletePaymentMethod(context.Context, *DeletePaymentMethodRequest, *DeletePaymentMethodResponse) error
	DefaultPaymentMethod(context.Context, *DefaultPaymentMethodRequest, *DefaultPaymentMethodResponse) error
	Signup(context.Context, *SignupRequest, *SignupResponse) error
	Login(context.Context, *LoginRequest, *LoginResponse) error
	Token(context.Context, *TokenRequest, *TokenResponse) error
}

func RegisterAccountHandler(s server.Server, hdlr AccountHandler, opts ...server.HandlerOption) error {
	type account interface {
		ReadUser(ctx context.Context, in *ReadUserRequest, out *ReadUserResponse) error
		UpdateUser(ctx context.Context, in *UpdateUserRequest, out *UpdateUserResponse) error
		DeleteUser(ctx context.Context, in *DeleteUserRequest, out *DeleteUserResponse) error
		ListPlans(ctx context.Context, in *ListPlansRequest, out *ListPlansResponse) error
		CreateSubscription(ctx context.Context, in *CreateSubscriptionRequest, out *CreateSubscriptionResponse) error
		CreatePaymentMethod(ctx context.Context, in *CreatePaymentMethodRequest, out *CreatePaymentMethodResponse) error
		DeletePaymentMethod(ctx context.Context, in *DeletePaymentMethodRequest, out *DeletePaymentMethodResponse) error
		DefaultPaymentMethod(ctx context.Context, in *DefaultPaymentMethodRequest, out *DefaultPaymentMethodResponse) error
		Signup(ctx context.Context, in *SignupRequest, out *SignupResponse) error
		Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error
		Token(ctx context.Context, in *TokenRequest, out *TokenResponse) error
	}
	type Account struct {
		account
	}
	h := &accountHandler{hdlr}
	return s.Handle(s.NewHandler(&Account{h}, opts...))
}

type accountHandler struct {
	AccountHandler
}

func (h *accountHandler) ReadUser(ctx context.Context, in *ReadUserRequest, out *ReadUserResponse) error {
	return h.AccountHandler.ReadUser(ctx, in, out)
}

func (h *accountHandler) UpdateUser(ctx context.Context, in *UpdateUserRequest, out *UpdateUserResponse) error {
	return h.AccountHandler.UpdateUser(ctx, in, out)
}

func (h *accountHandler) DeleteUser(ctx context.Context, in *DeleteUserRequest, out *DeleteUserResponse) error {
	return h.AccountHandler.DeleteUser(ctx, in, out)
}

func (h *accountHandler) ListPlans(ctx context.Context, in *ListPlansRequest, out *ListPlansResponse) error {
	return h.AccountHandler.ListPlans(ctx, in, out)
}

func (h *accountHandler) CreateSubscription(ctx context.Context, in *CreateSubscriptionRequest, out *CreateSubscriptionResponse) error {
	return h.AccountHandler.CreateSubscription(ctx, in, out)
}

func (h *accountHandler) CreatePaymentMethod(ctx context.Context, in *CreatePaymentMethodRequest, out *CreatePaymentMethodResponse) error {
	return h.AccountHandler.CreatePaymentMethod(ctx, in, out)
}

func (h *accountHandler) DeletePaymentMethod(ctx context.Context, in *DeletePaymentMethodRequest, out *DeletePaymentMethodResponse) error {
	return h.AccountHandler.DeletePaymentMethod(ctx, in, out)
}

func (h *accountHandler) DefaultPaymentMethod(ctx context.Context, in *DefaultPaymentMethodRequest, out *DefaultPaymentMethodResponse) error {
	return h.AccountHandler.DefaultPaymentMethod(ctx, in, out)
}

func (h *accountHandler) Signup(ctx context.Context, in *SignupRequest, out *SignupResponse) error {
	return h.AccountHandler.Signup(ctx, in, out)
}

func (h *accountHandler) Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.AccountHandler.Login(ctx, in, out)
}

func (h *accountHandler) Token(ctx context.Context, in *TokenRequest, out *TokenResponse) error {
	return h.AccountHandler.Token(ctx, in, out)
}
