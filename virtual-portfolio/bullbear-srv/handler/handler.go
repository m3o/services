package handler

import (
	"context"

	proto "github.com/kytra-app/bullbear-srv/proto"
	"github.com/kytra-app/bullbear-srv/storage"
)

// New returns an instance of Handler
func New(storage storage.Service) *Handler {
	return &Handler{db: storage}
}

// Handler is an object can process RPC requests
type Handler struct{ db storage.Service }

// Get looks up a resource and returns the uuids of bullish & bearish users
func (h *Handler) Get(ctx context.Context, query *proto.Resource, rsp *proto.Response) error {
	if query.Uuid == "" || query.Type == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	r, err := h.db.Get(storage.Resource{UUID: query.Uuid, Type: query.Type})
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching resource"}
		return err
	}

	rsp.Resource = h.serializeResource(r)
	return nil
}

// Create sets a users opinion (bullish/bearish/none) on a resource
func (h *Handler) Create(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	if req.Resource == nil || req.UserUuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	err := h.db.Create(h.permittedParams(req))
	if err != nil {
		rsp.Error = &proto.Error{Code: 400, Message: err.Error()}
		return err
	}

	return nil
}

// List bulk looks up a resources and returns the number of bullish & bearish users
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	if req.ResourceType == "" || len(req.ResourceUuids) == 0 {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid Request"}
		return nil
	}

	rs, err := h.db.List(req.ResourceType, req.ResourceUuids, req.UserUuid)
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

func (h *Handler) permittedParams(req *proto.Request) storage.Opinion {
	return storage.Opinion{
		UserUUID:     req.UserUuid,
		Opinion:      req.Opinion.String(),
		ResourceUUID: req.Resource.Uuid,
		ResourceType: req.Resource.Type,
	}
}

func (h *Handler) serializeResource(r storage.Resource) *proto.Resource {
	var opinion proto.Opinion
	switch r.Opinion {
	case "BULLISH":
		opinion = proto.Opinion_BULLISH
		break
	case "BEARISH":
		opinion = proto.Opinion_BEARISH
		break
	default:
		opinion = proto.Opinion_NONE
	}

	return &proto.Resource{
		Uuid:       r.UUID,
		Type:       r.Type,
		Bulls:      r.Bulls,
		Bears:      r.Bears,
		BullsCount: int32(r.BullsCount),
		BearsCount: int32(r.BearsCount),
		Opinion:    opinion,
	}
}
