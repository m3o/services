package handler

import (
	"context"
	"time"

	auth "github.com/kytra-app/helpers/authentication"
	proto "github.com/kytra-app/notifications-api/proto"
	notifications "github.com/kytra-app/notifications-srv/proto"
	push "github.com/kytra-app/push-notifications-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth          auth.Authenticator
	push          push.PushNotificationsService
	notifications notifications.NotificationsService
}

// New returns an instance of Handler
func New(client client.Client, auth auth.Authenticator) Handler {
	return Handler{
		auth:          auth,
		push:          push.NewPushNotificationsService("kytra-srv-v1-push-notifications:8080", client),
		notifications: notifications.NewNotificationsService("kytra-srv-v1-notifications:8080", client),
	}
}

// Get returns the users recent notifications
func (h Handler) Get(ctx context.Context, req *proto.Query, rsp *proto.Response) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHORIZED", "You must be logged in to use this API")
	}

	// Fetch the notifications
	nRsp, err := h.notifications.ListNotifications(ctx, &notifications.Query{
		UserUuid: u.UUID, Page: req.Page, Limit: req.Limit,
	})
	if err != nil {
		return err
	}

	// Serialize the data
	rsp.Notifications = make([]*proto.Notification, len(nRsp.Notifications))
	for i, n := range nRsp.Notifications {
		rsp.Notifications[i] = &proto.Notification{
			Uuid:         n.Uuid,
			CreatedAt:    time.Unix(n.CreatedAt, 0).String(),
			Seen:         n.Seen,
			Title:        n.Title,
			Description:  n.Description,
			ResourceType: n.ResourceType,
			ResourceUuid: n.ResourceUuid,
		}
	}

	return nil
}

// Seen marks all of the users unseen notificatons as seen
func (h Handler) Seen(ctx context.Context, req *proto.Query, rsp *proto.Response) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHORIZED", "You must be logged in to use this API")
	}

	// Perform the request
	_, err = h.notifications.SetNotificationsSeen(ctx, &notifications.NotificationsSeenRequest{
		UserUuid: u.UUID,
	})

	return err
}

// RegisterPushToken saves the users mobile push notificaton token
func (h Handler) RegisterPushToken(ctx context.Context, req *proto.PushToken, rsp *proto.Response) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHORIZED", "You must be logged in to use this API")
	}

	// Create the token
	_, err = h.push.RegisterToken(ctx, &push.Token{Token: req.Token, UserUuid: u.UUID})
	return err
}
