package handler

import (
	"context"

	followers "github.com/kytra-app/followers-srv/proto"
	proto "github.com/kytra-app/investors-api/proto"
)

// Connections retries a list of investors the user currently follows
func (h Handler) Connections(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Step 1. Get the current user
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return err
	}

	// Step 2. Get the users their following
	fRsp, err := h.followers.Get(ctx, &followers.Resource{Type: "User", Uuid: user.UUID})
	if err != nil {
		return err
	}
	uuids := []string{}
	for _, r := range fRsp.GetFollowing() {
		if r.Type == "User" {
			uuids = append(uuids, r.Uuid)
		}
	}

	// Step 3. Serialize the response
	return h.serializeUsers(ctx, rsp, uuids)
}
