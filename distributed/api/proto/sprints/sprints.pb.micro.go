// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/sprints/sprints.proto

package go_micro_api_distributed

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

// Client API for DistributedSprints service

type DistributedSprintsService interface {
	CreateSprint(ctx context.Context, in *CreateSprintRequest, opts ...client.CallOption) (*CreateSprintResponse, error)
	ListSprints(ctx context.Context, in *ListSprintsRequest, opts ...client.CallOption) (*ListSprintsResponse, error)
	ReadSprint(ctx context.Context, in *ReadSprintRequest, opts ...client.CallOption) (*ReadSprintResponse, error)
	CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...client.CallOption) (*CreateTaskResponse, error)
	UpdateTask(ctx context.Context, in *UpdateTaskRequest, opts ...client.CallOption) (*UpdateTaskResponse, error)
	DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...client.CallOption) (*DeleteTaskResponse, error)
	CreateObjective(ctx context.Context, in *CreateObjectiveRequest, opts ...client.CallOption) (*CreateObjectiveResponse, error)
	UpdateObjective(ctx context.Context, in *UpdateObjectiveRequest, opts ...client.CallOption) (*UpdateObjectiveResponse, error)
	DeleteObjective(ctx context.Context, in *DeleteObjectiveRequest, opts ...client.CallOption) (*DeleteObjectiveResponse, error)
}

type distributedSprintsService struct {
	c    client.Client
	name string
}

func NewDistributedSprintsService(name string, c client.Client) DistributedSprintsService {
	return &distributedSprintsService{
		c:    c,
		name: name,
	}
}

func (c *distributedSprintsService) CreateSprint(ctx context.Context, in *CreateSprintRequest, opts ...client.CallOption) (*CreateSprintResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.CreateSprint", in)
	out := new(CreateSprintResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) ListSprints(ctx context.Context, in *ListSprintsRequest, opts ...client.CallOption) (*ListSprintsResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.ListSprints", in)
	out := new(ListSprintsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) ReadSprint(ctx context.Context, in *ReadSprintRequest, opts ...client.CallOption) (*ReadSprintResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.ReadSprint", in)
	out := new(ReadSprintResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...client.CallOption) (*CreateTaskResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.CreateTask", in)
	out := new(CreateTaskResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) UpdateTask(ctx context.Context, in *UpdateTaskRequest, opts ...client.CallOption) (*UpdateTaskResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.UpdateTask", in)
	out := new(UpdateTaskResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...client.CallOption) (*DeleteTaskResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.DeleteTask", in)
	out := new(DeleteTaskResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) CreateObjective(ctx context.Context, in *CreateObjectiveRequest, opts ...client.CallOption) (*CreateObjectiveResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.CreateObjective", in)
	out := new(CreateObjectiveResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) UpdateObjective(ctx context.Context, in *UpdateObjectiveRequest, opts ...client.CallOption) (*UpdateObjectiveResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.UpdateObjective", in)
	out := new(UpdateObjectiveResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *distributedSprintsService) DeleteObjective(ctx context.Context, in *DeleteObjectiveRequest, opts ...client.CallOption) (*DeleteObjectiveResponse, error) {
	req := c.c.NewRequest(c.name, "DistributedSprints.DeleteObjective", in)
	out := new(DeleteObjectiveResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DistributedSprints service

type DistributedSprintsHandler interface {
	CreateSprint(context.Context, *CreateSprintRequest, *CreateSprintResponse) error
	ListSprints(context.Context, *ListSprintsRequest, *ListSprintsResponse) error
	ReadSprint(context.Context, *ReadSprintRequest, *ReadSprintResponse) error
	CreateTask(context.Context, *CreateTaskRequest, *CreateTaskResponse) error
	UpdateTask(context.Context, *UpdateTaskRequest, *UpdateTaskResponse) error
	DeleteTask(context.Context, *DeleteTaskRequest, *DeleteTaskResponse) error
	CreateObjective(context.Context, *CreateObjectiveRequest, *CreateObjectiveResponse) error
	UpdateObjective(context.Context, *UpdateObjectiveRequest, *UpdateObjectiveResponse) error
	DeleteObjective(context.Context, *DeleteObjectiveRequest, *DeleteObjectiveResponse) error
}

func RegisterDistributedSprintsHandler(s server.Server, hdlr DistributedSprintsHandler, opts ...server.HandlerOption) error {
	type distributedSprints interface {
		CreateSprint(ctx context.Context, in *CreateSprintRequest, out *CreateSprintResponse) error
		ListSprints(ctx context.Context, in *ListSprintsRequest, out *ListSprintsResponse) error
		ReadSprint(ctx context.Context, in *ReadSprintRequest, out *ReadSprintResponse) error
		CreateTask(ctx context.Context, in *CreateTaskRequest, out *CreateTaskResponse) error
		UpdateTask(ctx context.Context, in *UpdateTaskRequest, out *UpdateTaskResponse) error
		DeleteTask(ctx context.Context, in *DeleteTaskRequest, out *DeleteTaskResponse) error
		CreateObjective(ctx context.Context, in *CreateObjectiveRequest, out *CreateObjectiveResponse) error
		UpdateObjective(ctx context.Context, in *UpdateObjectiveRequest, out *UpdateObjectiveResponse) error
		DeleteObjective(ctx context.Context, in *DeleteObjectiveRequest, out *DeleteObjectiveResponse) error
	}
	type DistributedSprints struct {
		distributedSprints
	}
	h := &distributedSprintsHandler{hdlr}
	return s.Handle(s.NewHandler(&DistributedSprints{h}, opts...))
}

type distributedSprintsHandler struct {
	DistributedSprintsHandler
}

func (h *distributedSprintsHandler) CreateSprint(ctx context.Context, in *CreateSprintRequest, out *CreateSprintResponse) error {
	return h.DistributedSprintsHandler.CreateSprint(ctx, in, out)
}

func (h *distributedSprintsHandler) ListSprints(ctx context.Context, in *ListSprintsRequest, out *ListSprintsResponse) error {
	return h.DistributedSprintsHandler.ListSprints(ctx, in, out)
}

func (h *distributedSprintsHandler) ReadSprint(ctx context.Context, in *ReadSprintRequest, out *ReadSprintResponse) error {
	return h.DistributedSprintsHandler.ReadSprint(ctx, in, out)
}

func (h *distributedSprintsHandler) CreateTask(ctx context.Context, in *CreateTaskRequest, out *CreateTaskResponse) error {
	return h.DistributedSprintsHandler.CreateTask(ctx, in, out)
}

func (h *distributedSprintsHandler) UpdateTask(ctx context.Context, in *UpdateTaskRequest, out *UpdateTaskResponse) error {
	return h.DistributedSprintsHandler.UpdateTask(ctx, in, out)
}

func (h *distributedSprintsHandler) DeleteTask(ctx context.Context, in *DeleteTaskRequest, out *DeleteTaskResponse) error {
	return h.DistributedSprintsHandler.DeleteTask(ctx, in, out)
}

func (h *distributedSprintsHandler) CreateObjective(ctx context.Context, in *CreateObjectiveRequest, out *CreateObjectiveResponse) error {
	return h.DistributedSprintsHandler.CreateObjective(ctx, in, out)
}

func (h *distributedSprintsHandler) UpdateObjective(ctx context.Context, in *UpdateObjectiveRequest, out *UpdateObjectiveResponse) error {
	return h.DistributedSprintsHandler.UpdateObjective(ctx, in, out)
}

func (h *distributedSprintsHandler) DeleteObjective(ctx context.Context, in *DeleteObjectiveRequest, out *DeleteObjectiveResponse) error {
	return h.DistributedSprintsHandler.DeleteObjective(ctx, in, out)
}
