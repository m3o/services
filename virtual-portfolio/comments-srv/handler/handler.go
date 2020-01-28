package handler

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/rabbitmq"

	proto "github.com/kytra-app/comments-srv/proto"
	"github.com/kytra-app/comments-srv/storage"
)

// New returns an instance of Handler
func New(storage storage.Service, broker broker.Broker) *Handler {
	return &Handler{db: storage, broker: broker}
}

// Handler is an object can process RPC requests
type Handler struct {
	db     storage.Service
	broker broker.Broker
}

// GetResource looks up a resource and returns the comments
func (h *Handler) GetResource(ctx context.Context, query *proto.Resource, rsp *proto.Response) error {
	if query.Uuid == "" || query.Type == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	r, err := h.db.GetResource(storage.Resource{UUID: query.Uuid, Type: query.Type})
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching resource"}
		return err
	}

	rsp.Resource = h.serializeResource(r)
	return nil
}

// Get finds a comment by UUID
func (h *Handler) Get(ctx context.Context, req *proto.Comment, rsp *proto.Response) error {
	if req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	c, err := h.db.Get(storage.Comment{UUID: req.Uuid})
	if err != nil {
		return err
	}

	rsp.Comment = &proto.Comment{
		Uuid:     c.UUID,
		Text:     c.Text,
		UserUuid: c.UserUUID,
	}

	return nil
}

// Create creates a comment on the resource
func (h *Handler) Create(ctx context.Context, req *proto.Comment, rsp *proto.Response) error {
	if req.Resource == nil || req.UserUuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	c, err := h.db.Create(h.permittedParams(req))
	if err != nil {
		rsp.Error = &proto.Error{Code: 400, Message: err.Error()}
		return err
	}

	rsp.Comment = &proto.Comment{
		Uuid:     c.UUID,
		Text:     c.Text,
		UserUuid: c.UserUUID,
		Resource: &proto.Resource{
			Uuid: c.ResourceUUID,
			Type: c.ResourceType,
		},
	}

	bytes, err := json.Marshal(&rsp.Comment)
	if err != nil {
		return err
	}
	return h.broker.Publish("kytra-srv-v1-comments-comment-created", &broker.Message{Body: bytes})
}

// Delete destroys a comment
func (h *Handler) Delete(ctx context.Context, req *proto.Comment, rsp *proto.Response) error {
	if req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	err := h.db.Delete(req.Uuid)
	if err != nil {
		rsp.Error = &proto.Error{Code: 400, Message: err.Error()}
		return err
	}

	return nil
}

// ListResources bulk looks up a resources and returns the comments
func (h *Handler) ListResources(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	if req.ResourceType == "" || len(req.ResourceUuids) == 0 {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	rs, err := h.db.ListResources(req.ResourceType, req.ResourceUuids)
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching resource"}
		return err
	}

	rsp.Resources = make([]*proto.Resource, len(rs))
	for i, r := range rs {
		rsp.Resources[i] = h.serializeResource(*r)
	}

	return nil
}

func (h *Handler) permittedParams(req *proto.Comment) storage.Comment {
	return storage.Comment{
		Text:         req.Text,
		UserUUID:     req.UserUuid,
		ResourceUUID: req.Resource.Uuid,
		ResourceType: req.Resource.Type,
	}
}

func (h *Handler) serializeResource(r storage.Resource) *proto.Resource {
	comments := make([]*proto.Comment, len(r.Comments))

	for i, c := range r.Comments {
		comments[i] = &proto.Comment{Uuid: c.UUID, Text: c.Text, UserUuid: c.UserUUID}
	}

	return &proto.Resource{Uuid: r.UUID, Type: r.Type, Comments: comments}
}
