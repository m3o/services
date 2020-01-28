package api

import (
	"context"
	"time"

	proto "github.com/kytra-app/notifications-srv/proto"
	"github.com/kytra-app/notifications-srv/storage"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(client client.Client, db storage.Service) Handler {
	return Handler{
		db: db,
	}
}

// Handler is an object which processes various messages
type Handler struct {
	db storage.Service
}

// ListNotifications returns the notifications for a user, which matches the given query
func (h Handler) ListNotifications(ctx context.Context, req *proto.Query, rsp *proto.Response) error {
	query := storage.Query{UserUUID: req.UserUuid, OnlyUnseen: req.OnlyUnseen}
	if req.StartTime != 0 {
		startTime := time.Unix(req.StartTime, 0)
		query.StartTime = &startTime
	}
	if req.EndTime != 0 {
		endTime := time.Unix(req.EndTime, 0)
		query.EndTime = &endTime
	}

	notifications, err := h.db.List(query)
	if err != nil {
		return err
	}

	rsp.Notifications = make([]*proto.Notification, len(notifications))
	for i, n := range notifications {
		rsp.Notifications[i] = &proto.Notification{
			Uuid:         n.UUID,
			CreatedAt:    n.CreatedAt.Unix(),
			UserUuid:     n.UserUUID,
			Seen:         n.Seen,
			Title:        n.Title,
			Description:  n.Description,
			ResourceType: n.ResourceType,
			ResourceUuid: n.ResourceUUID,
		}
	}

	return nil
}

// SetNotificationsSeen marks a users unseen notifications as seen
func (h Handler) SetNotificationsSeen(ctx context.Context, req *proto.NotificationsSeenRequest, rsp *proto.Response) error {
	if req.UserUuid == "" {
		return errors.BadRequest("MISSING_UUID", "A UserUUID is required to list notifications")
	}

	return h.db.SetNotificationsSeen(req.UserUuid)
}
