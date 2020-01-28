package handler

import (
	"context"

	proto "github.com/kytra-app/investors-api/proto"
	users "github.com/kytra-app/users-srv/proto"
)


// Search retrieves all investors which match the criteria
func (h Handler) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.ListResponse) error {
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return err
	}

	uRsp, err := h.users.Search(ctx, &users.SearchRequest{
		Query: req.Query, Limit: 20,
	})
	if err != nil {
		return err
	}

	uuids := []string{}
	for _, u := range uRsp.GetUsers() {
		if u.Uuid == user.UUID {
			continue
		}

		uuids = append(uuids, u.Uuid)
	}

	return h.serializeUsers(ctx, rsp, uuids)
}
