package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	comments "github.com/kytra-app/comments-srv/proto"
	notifiableEvent "github.com/kytra-app/notifications-srv/helpers/notifiableevent"
	"github.com/kytra-app/notifications-srv/storage"
	posts "github.com/kytra-app/posts-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/broker"
)

// Comment is the JSON object published by the comments-srv
type Comment struct {
	UUID     string `json:"uuid"`
	UserUUID string `json:"user_uuid"`
	Text     string `json:"text"`
	Resource struct {
		UUID string `json:"uuid"`
		Type string `json:"type"`
	} `json:"resource"`
}

// ConsumeNewComment creates notifications for any users who are tagged in a comment
func (c Consumer) ConsumeNewComment(e broker.Event) error {
	fmt.Printf("[ConsumeNewComment] Processing Message: %v\n", string(e.Message().Body))

	// Create a new notifiable event to keep track of the notifications
	event := notifiableEvent.New(c.db, c.pushSrv)

	// Decode the message
	var comment Comment
	if err := json.Unmarshal(e.Message().Body, &comment); err != nil {
		return err
	}
	if comment.Resource.Type != "Post" {
		fmt.Println("This comment was not made about a post")
		return nil
	}

	// Find the user who made the comment
	commentingUser, err := c.usersSrv.Find(context.Background(), &users.User{Uuid: comment.UserUUID})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	// Notify the tagged users
	for _, uuid := range c.usersTaggedInText(comment.Text) {
		event.SendNotification(storage.Notification{
			UserUUID:     uuid,
			Title:        fmt.Sprintf("%v tagged you in a comment", commentingUser.FirstName),
			Description:  comment.Text,
			Emoji:        "✍️",
			ResourceType: "Post",
			ResourceUUID: comment.Resource.UUID,
		})
	}

	// Get the post the comment was made about
	pRsp, err := c.postsSrv.Get(context.Background(), &posts.Post{Uuid: comment.Resource.UUID})
	if err != nil {
		return err
	}
	post := pRsp.Post

	// Notify the user who made the post
	if post.UserUuid != commentingUser.Uuid {
		event.SendNotification(storage.Notification{
			UserUUID:     post.UserUuid,
			Title:        fmt.Sprintf("%v commented on your post", commentingUser.FirstName),
			Description:  comment.Text,
			Emoji:        "✍️",
			ResourceType: "Post",
			ResourceUUID: comment.Resource.UUID,
		})
	}

	// Notify other users who commented on the post
	resource := &comments.Resource{Type: comment.Resource.Type, Uuid: comment.Resource.UUID}
	commentsRsp, err := c.commentsSrv.GetResource(context.Background(), resource)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	for _, com := range commentsRsp.Resource.Comments {
		if com.UserUuid == commentingUser.Uuid {
			continue
		}

		event.SendNotification(storage.Notification{
			UserUUID:     com.UserUuid,
			Title:        fmt.Sprintf("%v commented on a post", commentingUser.FirstName),
			Description:  fmt.Sprintf("%v left a comment after you on '%v'", commentingUser.FirstName, post.Title),
			Emoji:        "✍️",
			ResourceType: "Post",
			ResourceUUID: comment.Resource.UUID,
		})
	}

	fmt.Println("Finished successfully")
	return nil
}
