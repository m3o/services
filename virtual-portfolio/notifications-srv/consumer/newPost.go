package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	followers "github.com/kytra-app/followers-srv/proto"
	notifiableEvent "github.com/kytra-app/notifications-srv/helpers/notifiableevent"
	"github.com/kytra-app/notifications-srv/storage"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/broker"
)

// Post is the JSON object published by the posts-srv
type Post struct {
	UUID     string `json:"uuid"`
	UserUUID string `json:"user_uuid"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

// ConsumeNewPost creates notifications for any users who are tagged in a post
func (c Consumer) ConsumeNewPost(e broker.Event) error {
	fmt.Printf("[ConsumeNewPost] Processing Message: %v\n", string(e.Message().Body))

	// Create a new notifiable event to keep track of the notifications
	event := notifiableEvent.New(c.db, c.pushSrv)

	// Decode the message
	var post Post
	if err := json.Unmarshal(e.Message().Body, &post); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	// Find the user who made the post
	user, err := c.usersSrv.Find(context.Background(), &users.User{Uuid: post.UserUUID})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	// Create the notificiations
	for _, uuid := range c.usersTaggedInText(post.Text) {
		event.SendNotification(storage.Notification{
			UserUUID:     uuid,
			Title:        fmt.Sprintf("%v tagged you in a post", user.FirstName),
			Description:  post.Title,
			Emoji:        "✍️",
			ResourceType: "Post",
			ResourceUUID: post.UUID,
		})
	}

	// Notify the users who follow the posting user
	fRsp, err := c.followersSrv.Get(context.Background(), &followers.Resource{Type: "User", Uuid: user.Uuid})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, f := range fRsp.Followers {
		event.SendNotification(storage.Notification{
			UserUUID:     f.Uuid,
			Title:        fmt.Sprintf("%v just shared a post", user.FirstName),
			Description:  post.Title,
			Emoji:        "✍️",
			ResourceType: "Post",
			ResourceUUID: post.UUID,
		})
	}

	fmt.Println("Finished successfully")
	return nil
}
