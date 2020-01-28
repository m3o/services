package notifableevent

import (
	"context"
	"fmt"

	push "github.com/kytra-app/push-notifications-srv/proto"

	"github.com/kytra-app/notifications-srv/storage"
)

// New returns an instance of NotifiableEvent
func New(db storage.Service, push push.PushNotificationsService) NotifiableEvent {
	return NotifiableEvent{
		db:                db,
		push:              push,
		notifiedUserUUIDs: []string{},
	}
}

// NotifiableEvent is an event a group of users should be notified about
type NotifiableEvent struct {
	db                storage.Service
	push              push.PushNotificationsService
	notifiedUserUUIDs []string
}

// SendNotification sends a notoification to a user
func (ne NotifiableEvent) SendNotification(n storage.Notification) error {
	if !ne.shouldSendToUser(n.UserUUID) {
		return nil
	}

	if _, err := ne.db.Create(n); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	go ne.sendPush(n)
	ne.notifiedUserUUIDs = append(ne.notifiedUserUUIDs, n.UserUUID)

	return nil
}

// shouldSendToUser checks if a user has already been notified for this event
func (ne NotifiableEvent) shouldSendToUser(uuid string) bool {
	for _, u := range ne.notifiedUserUUIDs {
		if u == uuid {
			return false
		}
	}

	return true
}

// sendPush notifies the user via APNS or FCM
func (ne NotifiableEvent) sendPush(n storage.Notification) {
	title := n.Title
	if n.Emoji != "" {
		title = fmt.Sprintf("%v %v", n.Emoji, n.Title)
	}

	subtitle := n.Description
	if len(n.Description) > 60 {
		subtitle = fmt.Sprintf("%v...", subtitle[0:60])
	}

	ne.push.SendNotification(context.Background(), &push.Notification{
		UserUuid: n.UserUUID, Title: title, Subtitle: subtitle,
	})
}
