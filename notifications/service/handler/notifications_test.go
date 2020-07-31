package handler

import (
	"testing"
	"time"

	"github.com/m3o/services/notifications/service/dao"
	"github.com/micro/go-micro/v3/store/memory"

	"github.com/stretchr/testify/assert"
)

func TestNotify(t *testing.T) {

	n := &Notifications{}
	dao.Init(memory.NewStore())

	// basic tests
	assert.NoError(t, n.subscribe("user1", "fooType", "1234"), "Error creating subscription")
	assert.NoError(t, n.notify("fooType", "1234", "Message 1"), "Error creating notification")
	assertNumberNotifsForUser(n, t, "user1", 1)
	assert.NoError(t, n.notify("fooType", "1234", "Message 2"), "Error creating notification")
	assertNumberNotifsForUser(n, t, "user1", 2)
	assertNumberSubsForUser(n, t, "user1", 1)

	// test with another subscriber
	assert.NoError(t, n.subscribe("user2", "fooType", "1234"), "Error creating subscription")
	assert.NoError(t, n.notify("fooType", "1234", "Message 3"), "Error creating notification")
	assertNumberNotifsForUser(n, t, "user1", 3)
	assertNumberNotifsForUser(n, t, "user2", 1)
	assertNumberSubsForUser(n, t, "user2", 1)

	// test unsubscribe
	assert.NoError(t, n.unsubscribe("user1", "fooType", "1234"), "Error unsubscribing")
	assert.NoError(t, n.notify("fooType", "1234", "Message 4"), "Error creating notification")
	assertNumberNotifsForUser(n, t, "user1", 3)
	assertNumberNotifsForUser(n, t, "user2", 2)
	assertNumberSubsForUser(n, t, "user1", 0)

	// mark as read
	notifs, err := n.listNotifsForUser("user1")
	assert.NoError(t, err, "Error retrieving notification")
	notifIDs := make([]string, len(notifs))
	for i, v := range notifs {
		notifIDs[i] = v.Id
	}
	now := time.Now().Unix()
	assert.NoError(t, n.markAsRead("user1", notifIDs))
	readNotifs, err := n.listNotifsForUser("user1")
	assert.Equal(t, len(notifs), len(readNotifs), "Different number of notifications")
	for _, v := range readNotifs {
		assert.LessOrEqual(t, now, v.ReadTimestamp, "Notification not marked as read correctly")
	}

}

func assertNumberNotifsForUser(n *Notifications, t *testing.T, userID string, length int) {
	ret, err := n.listNotifsForUser(userID)
	assert.NoError(t, err, "Error listing notifications for user")
	assert.Len(t, ret, length)
}

func assertNumberSubsForUser(n *Notifications, t *testing.T, userID string, length int) {
	subs, err := n.listSubscriptionsForSubscriber(userID)
	assert.NoError(t, err, "Error listing subscriptions")
	assert.Len(t, subs, length)

}
