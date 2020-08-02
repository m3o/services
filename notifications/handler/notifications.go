package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/m3o/services/notifications/dao"
	nproto "github.com/m3o/services/notifications/proto/notifications"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
)

type Notifications struct {
}

// Subscribe subscribes the user to notifications for the given resource
func (n *Notifications) Subscribe(ctx context.Context, req *nproto.SubscribeRequest, rsp *nproto.SubscribeResponse) error {
	// Create subscription
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("com.micro.service.notifications.subscribe", "Failed to determine user for subscribe request %+v", acc)
	}

	return n.subscribe(acc.ID, req.ResourceType, req.ResourceId)
}

func (n *Notifications) subscribe(userID, resourceType, resourceID string) error {
	sub := dao.NewSubscription(userID, resourceType, resourceID)
	// Store subscription
	if err := dao.CreateSubscription(sub); err != nil {
		return errors.InternalServerError("com.micro.service.notifications.subscribe", "Error creating subscription %s", err)
	}
	return nil
}

// Unsubscribe unsubscribes the user from notifications about the given resource
func (n *Notifications) Unsubscribe(ctx context.Context, req *nproto.UnsubscribeRequest, rsp *nproto.UnsubscribeResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("com.micro.service.notifications.subscribe", "Failed to determine user for subscribe request %+v", acc)
	}

	return n.unsubscribe(acc.ID, req.ResourceType, req.ResourceId)
}

func (n *Notifications) unsubscribe(userID, resourceType, resourceID string) error {
	sub, err := dao.ReadSubscription(userID, resourceType, resourceID)
	if err == store.ErrNotFound {
		return errors.BadRequest("com.micro.service.notifications.unsubscribe", "Subscription not found")
	}
	if err != nil {
		return err
	}
	fmt.Printf("Deleting %+v\n", sub)
	return dao.DeleteSubscription(sub.ID, userID, resourceType, resourceID)
}

// MarkAsRead marks the given notifications as read
func (n *Notifications) MarkAsRead(ctx context.Context, req *nproto.MarkAsReadRequest, rsp *nproto.MarkAsReadResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("com.micro.service.notifications.subscribe", "Failed to determine user for subscribe request")
	}
	return n.markAsRead(acc.ID, req.Ids)
}

func (n *Notifications) markAsRead(userID string, notifIDs []string) error {

	for _, notifID := range notifIDs {
		notif, err := dao.ReadNotification(userID, notifID)
		if err != nil {
			return err
		}
		if !notif.Read.IsZero() {
			continue
		}
		notif.Read = time.Now()
		if err := dao.UpdateNotification(notif); err != nil {
			return err
		}
	}
	return nil
}

// Notify creates a notification
func (n *Notifications) Notify(ctx context.Context, req *nproto.NotifyRequest, rsp *nproto.NotifyResponse) error {
	// TODO do we need to do any auth check here?
	return n.notify(req.ResourceType, req.ResourceId, req.Message)
}

func (n *Notifications) notify(resourceType, resourceID, message string) error {
	subs, err := dao.ListSubscriptionsForResource(resourceType, resourceID)
	if err != nil {
		return errors.InternalServerError("com.micro.service.notifications.notify", "Error retrieving subscriptions %s", err.Error())
	}
	for _, v := range subs {
		fmt.Printf("Notification for sub %+v\n", v)
		notif := dao.NewNotification(message, resourceID, resourceType, v.SubscriberID)
		if err := dao.CreateNotification(notif); err != nil {
			return errors.InternalServerError("com.micro.service.notifications.notify", "Error creating notification %s", err.Error())
		}
	}
	return nil
}

// List lists all the notifications for the user TODO think about pagination
func (n *Notifications) List(ctx context.Context, req *nproto.ListRequest, rsp *nproto.ListResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("com.micro.service.notifications.subscribe", "Failed to determine user for subscribe request")
	}

	ret, err := n.listNotifsForUser(acc.ID)
	if err != nil {
		return err
	}
	rsp.Notifications = ret
	return nil
}

func (n *Notifications) listNotifsForUser(userID string) ([]*nproto.Notification, error) {
	// TODO - pagination
	notifs, err := dao.ListNotificationsForSubscriber(userID)
	if err != nil {
		return nil, errors.InternalServerError("com.micro.service.notifications.list", "Error retrieving notifications %s", err)
	}
	ret := make([]*nproto.Notification, len(notifs))
	for i, v := range notifs {
		ret[i] = &nproto.Notification{
			Id:            v.ID,
			Message:       v.Message,
			ReadTimestamp: v.Read.Unix(),
			ResourceId:    v.ResourceID,
			ResourceType:  v.ResourceType,
			Timestamp:     v.Created.Unix(),
		}
	}

	return ret, nil
}

// ListSubscriptions returns a list of subscriptions for the user
func (n *Notifications) ListSubscriptions(ctx context.Context, req *nproto.ListSubscribptionsRequest, rsp *nproto.ListSubscriptionsResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("com.micro.service.notifications.subscribe", "Failed to determine user for subscribe request")
	}
	subs, err := n.listSubscriptionsForSubscriber(acc.ID)
	if err != nil {
		return err
	}
	rsp.Subscriptions = subs
	return nil
}

func (n *Notifications) listSubscriptionsForSubscriber(userID string) ([]*nproto.Subscription, error) {
	subs, err := dao.ListSubscriptionsForSubscriber(userID)
	if err != nil {
		return nil, err
	}
	ret := make([]*nproto.Subscription, len(subs))
	for i, v := range subs {
		ret[i] = &nproto.Subscription{
			Id:           v.ID,
			ResourceId:   v.ResourceID,
			ResourceType: v.ResourceType,
		}
	}
	return ret, nil
}
