package handler

import (
	"context"
	"encoding/json"
	"fmt"

	followers "github.com/kytra-app/followers-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/broker"
)

// Post is the JSON object published by the post-srv
type Post struct {
	UUID     string `json:"uuid"`
	UserUUID string `json:"user_uuid"`
	FeedUUID string `json:"feed_uuid"`
	FeedType string `json:"feed_type"`
}

// HandleNewPost handless the event when a post is created
func (h Handler) HandleNewPost(e broker.Event) error {
	fmt.Printf("[HandleNewPost] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var post Post
	if err := json.Unmarshal(e.Message().Body, &post); err != nil {
		return err
	}

	// Fetch the user who made the post
	user, err := h.usersSrv.Find(context.Background(), &users.User{Uuid: post.UserUUID})
	if err != nil {
		return err
	}

	// Fetch the users followers
	userFollowers, err := h.followersForFeed("User", user.Uuid)
	if err != nil {
		return err
	}

	// Add a post into their feeds
	for _, uuid := range userFollowers {
		d := fmt.Sprintf("Because you follow <&User:%v>%v<&/User>", user.Uuid, user.FirstName)
		fmt.Println(d)

		if err := h.addPostToUserFeed(uuid, post.UUID, d); err != nil {
			return err
		}
	}

	// Create a post for the posting user
	if err := h.addPostToUserFeed(user.Uuid, post.UUID, "Posted by you"); err != nil {
		return err
	}

	// Stop here if the post wasn't made on a stock's feed
	if post.FeedType != "Stock" {
		fmt.Println("Ending because the post is not made against a stock")
		return nil
	}

	// Fetch the stock
	stockRsp, err := h.stocksSrv.Get(context.Background(), &stocks.Stock{Uuid: post.FeedUUID})
	if err != nil {
		fmt.Println("Ending because the stock could not be found")
		return err
	}
	stock := stockRsp.Stock

	// Fetch the stocks followers
	stockFollowers, err := h.followersForFeed("Stock", stock.Uuid)
	if err != nil {
		fmt.Println("Ending because we could not get the stocks followers")
		return err
	}
	fmt.Printf("%v has %v followers", stock.Name, len(stockFollowers))

	// Add a post into their feeds
	for _, uuid := range stockFollowers {
		d := fmt.Sprintf("Because you follow <&Stock:%v>%v<&/Stock>", stock.Uuid, stock.Name)
		fmt.Println(d)

		if err := h.addPostToUserFeed(uuid, post.UUID, d); err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) followersForFeed(FeedType, FeedUUID string) ([]string, error) {
	query := followers.Resource{Type: FeedType, Uuid: FeedUUID}

	rsp, err := h.followerSrv.Get(context.Background(), &query)
	if err != nil {
		return []string{}, err
	}

	uuids := make([]string, len(rsp.Followers))
	for i, f := range rsp.Followers {
		uuids[i] = f.Uuid
	}

	return uuids, nil
}
