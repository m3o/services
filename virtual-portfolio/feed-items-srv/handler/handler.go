package handler

import (
	"context"

	proto "github.com/kytra-app/feeditems-srv/proto"
	"github.com/kytra-app/feeditems-srv/storage"
)

// New returns an instance of Handler
func New(storage storage.Service) *Handler {
	return &Handler{db: storage}
}

// Handler is an object can process RPC requests
type Handler struct{ db storage.Service }

// Create inserts a new feed item into the database
func (h *Handler) Create(ctx context.Context, req *proto.FeedItem, rsp *proto.Response) error {
	p, err := h.db.Create(storage.FeedItem{
		FeedType:    req.FeedType,
		FeedUUID:    req.FeedUuid,
		Tag:         req.Tag,
		PostUUID:    req.PostUuid,
		Description: req.Description,
	})

	if err != nil {
		rsp.Error = &proto.Error{Code: 400, Message: err.Error()}
		return nil
	}

	rsp.Item = h.serializeFeedItem(&p)
	return nil
}

// Get looks up a feed item using the UUID
func (h *Handler) Get(ctx context.Context, req *proto.FeedItem, rsp *proto.Response) error {
	i, err := h.db.Get(storage.FeedItem{UUID: req.Uuid})

	if err != nil {
		rsp.Error = &proto.Error{Code: 404, Message: err.Error()}
		return nil
	}

	rsp.Item = h.serializeFeedItem(&i)
	return nil
}

// Delete removes a feed item using the UUID
func (h *Handler) Delete(ctx context.Context, req *proto.FeedItem, rsp *proto.Response) error {
	err := h.db.Delete(storage.FeedItem{UUID: req.Uuid})
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: err.Error()}
		return err
	}

	return nil
}

// BulkDelete removes a feed item using the UUID
func (h *Handler) BulkDelete(ctx context.Context, req *proto.BulkDeleteRequest, rsp *proto.Response) error {
	err := h.db.BulkDelete(req.FeedType, req.FeedUuid, req.PostUuids)

	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: err.Error()}
		return err
	}

	return nil
}

// GetFeed returns the recent posts for a feed
func (h *Handler) GetFeed(ctx context.Context, req *proto.GetFeedRequest, rsp *proto.GetFeedResponse) error {
	limit := req.Limit
	if req.Limit == 0 {
		limit = 30
	}

	items, err := h.db.GetFeed(req.Type, req.Uuid, req.Page, limit)
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: err.Error()}
		return err
	}

	rsp.Items = make([]*proto.FeedItem, len(items))
	for i, item := range items {
		rsp.Items[i] = h.serializeFeedItem(item)
	}

	return nil
}

func (h *Handler) serializeFeedItem(p *storage.FeedItem) *proto.FeedItem {
	return &proto.FeedItem{
		Uuid:        p.UUID,
		FeedType:    p.FeedType,
		FeedUuid:    p.FeedUUID,
		Tag:         p.Tag,
		PostUuid:    p.PostUUID,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Unix(),
	}
}
