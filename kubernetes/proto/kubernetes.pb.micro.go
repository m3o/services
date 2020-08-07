// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: kubernetes/proto/kubernetes.proto

package go_micro_service_kubernetes

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

// Api Endpoints for Kubernetes service

func NewKubernetesEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Kubernetes service

type KubernetesService interface {
	CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, opts ...client.CallOption) (*CreateNamespaceResponse, error)
	DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, opts ...client.CallOption) (*DeleteNamespaceResponse, error)
	CreateImagePullSecret(ctx context.Context, in *CreateImagePullSecretRequest, opts ...client.CallOption) (*CreateImagePullSecretResponse, error)
	DeleteImagePullSecret(ctx context.Context, in *DeleteImagePullSecretRequest, opts ...client.CallOption) (*DeleteImagePullSecretResponse, error)
	CreateServiceAccount(ctx context.Context, in *CreateServiceAccountRequest, opts ...client.CallOption) (*CreateServiceAccountResponse, error)
	DeleteServiceAccount(ctx context.Context, in *DeleteServiceAccountRequest, opts ...client.CallOption) (*DeleteServiceAccountResponse, error)
}

type kubernetesService struct {
	c    client.Client
	name string
}

func NewKubernetesService(name string, c client.Client) KubernetesService {
	return &kubernetesService{
		c:    c,
		name: name,
	}
}

func (c *kubernetesService) CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, opts ...client.CallOption) (*CreateNamespaceResponse, error) {
	req := c.c.NewRequest(c.name, "Kubernetes.CreateNamespace", in)
	out := new(CreateNamespaceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubernetesService) DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, opts ...client.CallOption) (*DeleteNamespaceResponse, error) {
	req := c.c.NewRequest(c.name, "Kubernetes.DeleteNamespace", in)
	out := new(DeleteNamespaceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubernetesService) CreateImagePullSecret(ctx context.Context, in *CreateImagePullSecretRequest, opts ...client.CallOption) (*CreateImagePullSecretResponse, error) {
	req := c.c.NewRequest(c.name, "Kubernetes.CreateImagePullSecret", in)
	out := new(CreateImagePullSecretResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubernetesService) DeleteImagePullSecret(ctx context.Context, in *DeleteImagePullSecretRequest, opts ...client.CallOption) (*DeleteImagePullSecretResponse, error) {
	req := c.c.NewRequest(c.name, "Kubernetes.DeleteImagePullSecret", in)
	out := new(DeleteImagePullSecretResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubernetesService) CreateServiceAccount(ctx context.Context, in *CreateServiceAccountRequest, opts ...client.CallOption) (*CreateServiceAccountResponse, error) {
	req := c.c.NewRequest(c.name, "Kubernetes.CreateServiceAccount", in)
	out := new(CreateServiceAccountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kubernetesService) DeleteServiceAccount(ctx context.Context, in *DeleteServiceAccountRequest, opts ...client.CallOption) (*DeleteServiceAccountResponse, error) {
	req := c.c.NewRequest(c.name, "Kubernetes.DeleteServiceAccount", in)
	out := new(DeleteServiceAccountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Kubernetes service

type KubernetesHandler interface {
	CreateNamespace(context.Context, *CreateNamespaceRequest, *CreateNamespaceResponse) error
	DeleteNamespace(context.Context, *DeleteNamespaceRequest, *DeleteNamespaceResponse) error
	CreateImagePullSecret(context.Context, *CreateImagePullSecretRequest, *CreateImagePullSecretResponse) error
	DeleteImagePullSecret(context.Context, *DeleteImagePullSecretRequest, *DeleteImagePullSecretResponse) error
	CreateServiceAccount(context.Context, *CreateServiceAccountRequest, *CreateServiceAccountResponse) error
	DeleteServiceAccount(context.Context, *DeleteServiceAccountRequest, *DeleteServiceAccountResponse) error
}

func RegisterKubernetesHandler(s server.Server, hdlr KubernetesHandler, opts ...server.HandlerOption) error {
	type kubernetes interface {
		CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, out *CreateNamespaceResponse) error
		DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, out *DeleteNamespaceResponse) error
		CreateImagePullSecret(ctx context.Context, in *CreateImagePullSecretRequest, out *CreateImagePullSecretResponse) error
		DeleteImagePullSecret(ctx context.Context, in *DeleteImagePullSecretRequest, out *DeleteImagePullSecretResponse) error
		CreateServiceAccount(ctx context.Context, in *CreateServiceAccountRequest, out *CreateServiceAccountResponse) error
		DeleteServiceAccount(ctx context.Context, in *DeleteServiceAccountRequest, out *DeleteServiceAccountResponse) error
	}
	type Kubernetes struct {
		kubernetes
	}
	h := &kubernetesHandler{hdlr}
	return s.Handle(s.NewHandler(&Kubernetes{h}, opts...))
}

type kubernetesHandler struct {
	KubernetesHandler
}

func (h *kubernetesHandler) CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, out *CreateNamespaceResponse) error {
	return h.KubernetesHandler.CreateNamespace(ctx, in, out)
}

func (h *kubernetesHandler) DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, out *DeleteNamespaceResponse) error {
	return h.KubernetesHandler.DeleteNamespace(ctx, in, out)
}

func (h *kubernetesHandler) CreateImagePullSecret(ctx context.Context, in *CreateImagePullSecretRequest, out *CreateImagePullSecretResponse) error {
	return h.KubernetesHandler.CreateImagePullSecret(ctx, in, out)
}

func (h *kubernetesHandler) DeleteImagePullSecret(ctx context.Context, in *DeleteImagePullSecretRequest, out *DeleteImagePullSecretResponse) error {
	return h.KubernetesHandler.DeleteImagePullSecret(ctx, in, out)
}

func (h *kubernetesHandler) CreateServiceAccount(ctx context.Context, in *CreateServiceAccountRequest, out *CreateServiceAccountResponse) error {
	return h.KubernetesHandler.CreateServiceAccount(ctx, in, out)
}

func (h *kubernetesHandler) DeleteServiceAccount(ctx context.Context, in *DeleteServiceAccountRequest, out *DeleteServiceAccountResponse) error {
	return h.KubernetesHandler.DeleteServiceAccount(ctx, in, out)
}
