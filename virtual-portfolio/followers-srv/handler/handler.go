package handler

import (
	"context"
	"encoding/json"

	proto "github.com/kytra-app/followers-srv/proto"
	"github.com/kytra-app/followers-srv/storage"
	"github.com/kytra-app/helpers/microtime"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/errors"
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

// Count returns the number of followers and followees a resource (user/stock/group) has
func (h *Handler) Count(ctx context.Context, req *proto.Resource, rsp *proto.Response) (err error) {
	r := storage.Resource{UUID: req.Uuid, Type: req.Type}
	if r.UUID == "" || r.Type == "" {
		return errors.BadRequest("INVALID_RESOURCE", "Invalid followers resource requested")
	}

	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	if rsp.FollowerCount, err = h.db.CountFollowers(time, r); err != nil {
		return err
	}

	if rsp.FollowingCount, err = h.db.CountFollowing(time, r); err != nil {
		return err
	}

	return nil
}

// Get returns the followers and followees for a resource (e.g. user/stock/group)
func (h *Handler) Get(ctx context.Context, req *proto.Resource, rsp *proto.Response) (err error) {
	r := storage.Resource{UUID: req.Uuid, Type: req.Type}
	if r.UUID == "" || r.Type == "" {
		return errors.BadRequest("INVALID_RESOURCE", "Invalid followers resource requested")
	}

	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	followers, err := h.db.GetFollowers(time, r)
	if err != nil {
		return err
	}

	following, err := h.db.GetFollowing(time, r)
	if err != nil {
		return err
	}

	for _, r := range followers {
		rsp.Followers = append(rsp.Followers, &proto.Resource{Uuid: r.UUID, Type: r.Type})
	}

	for _, r := range following {
		rsp.Following = append(rsp.Following, &proto.Resource{Uuid: r.UUID, Type: r.Type})
	}

	return nil
}

// List is depricated due to an insufficiently descriptive name. Use ListRelationships instead.
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	return h.ListRelationships(ctx, req, rsp)
}

// ListRelationships checks if a resource follows other resources, e.g. which of these stocks does a user follow.
func (h *Handler) ListRelationships(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	follower := storage.Resource{UUID: req.Follower.Uuid, Type: req.Follower.Type}
	if follower.UUID == "" || follower.Type == "" {
		return errors.BadRequest("INVALID_RESOURCE", "Invalid follower resource requested")
	}

	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	data, err := h.db.ListRelationships(time, follower, req.FolloweeType, req.FolloweeUuids)
	if err != nil {
		return err
	}

	rsp.Resources = make([]*proto.Resource, len(data))
	for i, r := range data {
		rsp.Resources[i] = &proto.Resource{
			Uuid:      r.UUID,
			Type:      r.Type,
			Following: r.Following,
		}
	}

	return nil
}

// Follow creates a new relationship between a follower and a followee
func (h *Handler) Follow(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	follower := storage.Resource{UUID: req.Follower.Uuid, Type: req.Follower.Type}
	followee := storage.Resource{UUID: req.Followee.Uuid, Type: req.Followee.Type}

	if err := h.db.Follow(follower, followee); err != nil {
		return err
	}

	bytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	return h.broker.Publish("kytra-srv-v1-followers-new-follow", &broker.Message{Body: bytes})
}

// Unfollow deletes the relationship between a follower and a followee
func (h *Handler) Unfollow(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	follower := storage.Resource{UUID: req.Follower.Uuid, Type: req.Follower.Type}
	followee := storage.Resource{UUID: req.Followee.Uuid, Type: req.Followee.Type}

	if err := h.db.Unfollow(follower, followee); err != nil {
		return err
	}

	bytes, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	return h.broker.Publish("kytra-srv-v1-followers-new-unfollow", &broker.Message{Body: bytes})
}
