package handler

import (
	"context"

	users "github.com/kytra-app/users-srv/proto"
	proto "github.com/kytra-app/welcome-api/proto"
)

// FindFriends looks up users by phone number
func (h Handler) FindFriends(ctx context.Context, req *proto.FindFriendsRequest, rsp *proto.FindFriendsResponse) error {
	if len(req.PhoneNumbers) == 0 {
		return nil
	}

	usersRsp, err := h.user.List(ctx, &users.ListRequest{PhoneNumbers: req.PhoneNumbers})
	if err != nil {
		return err
	}

	rsp.Users = make([]*proto.User, len(usersRsp.Users))
	for i, u := range usersRsp.Users {
		rsp.Users[i] = &proto.User{
			Uuid:              u.Uuid,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			Username:          u.Username,
			ProfilePictureUrl: h.photos.GetURL(u.ProfilePictureId),
		}
	}

	return nil
}
