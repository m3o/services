package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kytra-app/insights-srv/storage"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/broker"
)

// Post is the JSON object published by the post-srv
type Post struct {
	UUID     string `json:"uuid"`
	Title    string `json:"title"`
	UserUUID string `json:"user_uuid"`
	FeedUUID string `json:"feed_uuid"`
	FeedType string `json:"feed_type"`
}

// HandleNewPost handles the event when a post is created
func (h *Handler) HandleNewPost(e broker.Event) error {
	fmt.Printf("[HandleNewPost] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var post Post
	if err := json.Unmarshal(e.Message().Body, &post); err != nil {
		return err
	}
	fmt.Println(post)
	if post.FeedType != "Stock" {
		fmt.Println("Skipping non-stock post")
		return nil
	}

	// Fetch the user who made the post
	user, err := h.users.Find(context.Background(), &users.User{Uuid: post.UserUUID})
	if err != nil {
		return err
	}

	// Create a post in the feed
	i, err := h.db.CreateInsight(storage.Insight{
		Title:     fmt.Sprintf("%v shared a post: %v", user.FirstName, post.Title),
		Type:      "POST",
		AssetUUID: post.FeedUUID,
		AssetType: post.FeedType,
		PostUUID:  post.UUID,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	h.publishNewInsight(i)

	return nil
}
