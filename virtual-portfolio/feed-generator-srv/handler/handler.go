package handler

import (
	"context"
	"fmt"

	feeditems "github.com/kytra-app/feed-items-srv/proto"
	followers "github.com/kytra-app/followers-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
)

// New returns an instance of Handler
func New(client client.Client) Handler {
	return Handler{
		usersSrv:     users.NewUsersService("kytra-srv-v1-users:8080", client),
		postsSrv:     posts.NewPostsService("kytra-srv-v1-posts:8080", client),
		stocksSrv:    stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		followerSrv:  followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		feedItemsSrv: feeditems.NewFeedItemsService("kytra-srv-v1-feed-items:8080", client),
	}
}

// Handler is an object which processes various messages
type Handler struct {
	usersSrv     users.UsersService
	postsSrv     posts.PostsService
	stocksSrv    stocks.StocksService
	followerSrv  followers.FollowersService
	feedItemsSrv feeditems.FeedItemsService
}

func (h Handler) addPostToUserFeed(userUUID, postUUID, description string) error {
	item := feeditems.FeedItem{
		FeedType:    "User",
		FeedUuid:    userUUID,
		Tag:         "POST",
		PostUuid:    postUUID,
		Description: description,
	}

	rsp, err := h.feedItemsSrv.Create(context.Background(), &item)
	if err != nil {
		fmt.Printf("Error creating post: %v \n", err)
		return err
	}

	fmt.Printf("Created post! UUID: %v\n", rsp.Item.Uuid)
	return nil
}
