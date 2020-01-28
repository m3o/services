package handler

import (
	"context"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	proto "github.com/kytra-app/investors-api/proto"
	"github.com/micro/go-micro/errors"
)

// Follow creates a follower relationship between the requested user and the authenticated uses
func (h Handler) Follow(ctx context.Context, req *proto.User, rsp *proto.User) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	if req.Uuid == "" {
		return errors.BadRequest("UUID_REQUIRED", "A UUID is required")
	}

	if u.UUID == req.Uuid {
		return errors.BadRequest("UUID_INVALID", "You cannot follow yourself")
	}

	followReq := followers.Request{
		Followee: &followers.Resource{Uuid: req.Uuid, Type: "User"},
		Follower: &followers.Resource{Uuid: u.UUID, Type: "User"},
	}

	_, err = h.followers.Follow(ctx, &followReq)
	return err
}

// Unfollow deletes the follower relationship between the requested user and the authenticated uses
func (h Handler) Unfollow(ctx context.Context, req *proto.User, rsp *proto.User) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	unFollowReq := followers.Request{
		Followee: &followers.Resource{Uuid: req.Uuid, Type: "User"},
		Follower: &followers.Resource{Uuid: u.UUID, Type: "User"},
	}

	_, err = h.followers.Unfollow(ctx, &unFollowReq)
	return err
}

func (h Handler) getFollowingStatus(ctx context.Context, userUUID string) (bool, error) {
	rsp, err := h.getFollowingStatuses(ctx, []string{userUUID})

	if following, ok := rsp[userUUID]; ok {
		return following, nil
	}

	return false, err
}

func (h Handler) getFollowingStatuses(ctx context.Context, uuids []string) (map[string]bool, error) {
	follower, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return map[string]bool{}, nil
	}

	fContext, fCancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer fCancel()
	fRsp, err := h.followers.List(fContext, &followers.ListRequest{
		Follower:      &followers.Resource{Type: "User", Uuid: follower.UUID},
		FolloweeType:  "User",
		FolloweeUuids: uuids,
	})

	if err != nil {
		return map[string]bool{}, err
	}

	rsp := make(map[string]bool, len(fRsp.Resources))
	for _, r := range fRsp.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}
