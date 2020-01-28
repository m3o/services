package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/errors"
	_ "github.com/micro/go-plugins/broker/rabbitmq"

	proto "github.com/kytra-app/posts-srv/proto"
	"github.com/kytra-app/posts-srv/storage"
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

// Create inserts a new post into the database
func (h *Handler) Create(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	p, err := h.db.Create(h.permittedParams(req))
	if err != nil {
		rsp.Error = &proto.Error{Code: 400, Message: err.Error()}
		return nil
	}

	rsp.Post = h.serializePost(&p)

	bytes, err := json.Marshal(&rsp.Post)
	if err != nil {
		return err
	}
	brokerErr := h.broker.Publish("kytra-srv-v1-posts-post-created", &broker.Message{Body: bytes})
	if brokerErr != nil {
		fmt.Printf("Error Sending Msg to broker: %v\n", err)
	} else {
		fmt.Printf("Message sent to broker\n")
	}

	return nil
}

// Count returns the number of posts which match the query
func (h *Handler) Count(ctx context.Context, req *proto.Post, rsp *proto.CountResponse) error {
	c, err := h.db.Count(storage.Post{UserUUID: req.UserUuid})
	if err != nil {
		return err
	}

	rsp.Count = c
	return nil
}

// CountByUser returns the number of posts made by each user
func (h *Handler) CountByUser(ctx context.Context, req *proto.CountByUserRequest, rsp *proto.CountByUserResponse) error {
	if req.StartTime == 0 {
		return errors.BadRequest("START_TIME_REQUIRED", "A start time is required")
	}
	if req.EndTime == 0 {
		return errors.BadRequest("END_TIME_REQUIRED", "An end time is required")
	}

	startTime := time.Unix(req.StartTime, 0)
	endTime := time.Unix(req.EndTime, 0)
	c, err := h.db.CountByUser(req.UserUuids, startTime, endTime)
	if err != nil {
		return err
	}

	rsp.Counts = make([]*proto.CountResponse, len(c))
	var i int32
	for uuid, count := range c {
		rsp.Counts[i] = &proto.CountResponse{
			UserUuid: uuid,
			Count:    count,
		}
		i++
	}

	return nil
}

// Recent returns X recent posts
func (h *Handler) Recent(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	limit := req.Limit
	if req.Limit == 0 {
		limit = 30
	}

	posts, err := h.db.Recent(limit, req.Page)
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching posts"}
		return err
	}

	rsp.Posts = make([]*proto.Post, len(posts))
	for i, p := range posts {
		rsp.Posts[i] = h.serializePost(p)
	}

	return nil
}

// Update amends the post, found using the UUID
func (h *Handler) Update(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	if req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing UUID"}
		return nil
	}

	p, err := h.db.Update(storage.Post{
		UUID:                req.Uuid,
		Text:                req.Text,
		Title:               req.Title,
		AttachmentPictureID: req.AttachmentPictureId,
		AttachmentLinkURL:   req.AttachmentLinkUrl,
	})

	if err != nil {
		rsp.Error = &proto.Error{Code: 404}
		return err
	}

	rsp.Post = h.serializePost(&p)
	return nil
}

// Get returns the post, found using the UUID
func (h *Handler) Get(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	if req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing UUID"}
		return nil
	}

	p, err := h.db.Get(storage.Post{UUID: req.Uuid})
	if err != nil {
		rsp.Error = &proto.Error{Code: 404}
		return err
	}

	rsp.Post = h.serializePost(&p)
	return nil
}

// Delete deletes the post found using the UUID
func (h *Handler) Delete(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	if req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing UUID"}
		return nil
	}

	if err := h.db.Delete(storage.Post{UUID: req.Uuid}); err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching post"}
		return err
	}

	return nil
}

// List returns all the posts matching the UUIDs provided
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	if len(req.Uuids) == 0 {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing UUIDs"}
		return nil
	}

	posts, err := h.db.List(req.Uuids)
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching posts"}
		return err
	}

	rsp.Posts = make([]*proto.Post, len(posts))
	for i, p := range posts {
		rsp.Posts[i] = h.serializePost(p)
	}

	return nil
}

// ListFeed returns all the posts matching the feed type and uuid provided
func (h *Handler) ListFeed(ctx context.Context, req *proto.Feed, rsp *proto.ListResponse) error {
	if req.Type == "" || req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing Type or UUID"}
		return nil
	}

	posts, err := h.db.ListFeed(req.Type, req.Uuid)
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching posts"}
		return err
	}

	rsp.Posts = make([]*proto.Post, len(posts))
	for i, p := range posts {
		rsp.Posts[i] = h.serializePost(p)
	}

	return nil
}

// ListUser returns all the posts made by the user
func (h *Handler) ListUser(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	if req.UserUuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Invalid request: missing UUID"}
		return nil
	}

	limit := req.Limit
	if req.Limit == 0 {
		limit = 30
	}

	posts, err := h.db.ListUser(req.UserUuid, limit, req.Page)
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching posts"}
		return err
	}

	rsp.Posts = make([]*proto.Post, len(posts))
	for i, p := range posts {
		rsp.Posts[i] = h.serializePost(p)
	}

	return nil
}

func (h *Handler) serializePost(p *storage.Post) *proto.Post {
	return &proto.Post{
		Uuid:                p.UUID,
		Text:                p.Text,
		Title:               p.Title,
		UserUuid:            p.UserUUID,
		FeedType:            p.FeedType,
		FeedUuid:            p.FeedUUID,
		CreatedAt:           p.CreatedAt.Unix(),
		AttachmentPictureId: p.AttachmentPictureID,
		AttachmentLinkUrl:   p.AttachmentLinkURL,
	}
}

func (h *Handler) permittedParams(p *proto.Post) storage.Post {
	return storage.Post{
		Text:                p.Text,
		Title:               p.Title,
		UserUUID:            p.UserUuid,
		FeedType:            p.FeedType,
		FeedUUID:            p.FeedUuid,
		AttachmentLinkURL:   p.AttachmentLinkUrl,
		AttachmentPictureID: p.AttachmentPictureId,
	}
}
