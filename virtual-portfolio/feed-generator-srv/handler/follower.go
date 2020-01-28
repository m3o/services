package handler

import (
	"context"
	"encoding/json"
	"fmt"

	feeditems "github.com/kytra-app/feed-items-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/broker"
)

// Resource is a stock or user, in the follower domain
type Resource struct {
	UUID string `json:"uuid"`
	Type string `json:"type"`
}

// Follow is a relationship between a follower and followee
type Follow struct {
	Follower Resource `json:"follower"`
	Followee Resource `json:"followee"`
}

// HandleFollow handles the event when a user follows another resource
func (h Handler) HandleFollow(e broker.Event) error {
	fmt.Printf("[HandleFollow] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var follow Follow
	if err := json.Unmarshal(e.Message().Body, &follow); err != nil {
		fmt.Println(err)
		return err
	}

	// Get the posts to remove from the feed
	var postUUIDs []string
	var postErr error

	if follow.Followee.Type == "User" {
		// Users can make posts on feeds other than their own, so we query by UserUUID and
		// not Feed, this will include posts they made on stocks for example.
		postUUIDs, postErr = h.postUUIDsMadeByUser(follow.Followee.UUID, 10)
	} else {
		postUUIDs, postErr = h.postUUIDsInFeed(follow.Followee.Type, follow.Followee.UUID)
	}

	if postErr != nil {
		fmt.Println(postErr)
		return postErr
	}

	// Get the description to attach to the posts, e.g. "Because you follow John"
	desc, err := h.descriptionForFollowee(follow.Followee.Type, follow.Followee.UUID)
	if err != nil {
		return err
	}

	// Insert the posts intos the followers feed
	for _, uuid := range postUUIDs {
		if err := h.addPostToUserFeed(follow.Follower.UUID, uuid, desc); err != nil {
			return err
		}
	}

	return nil
}

// HandleUnfollow handles the event when a user unfollows another resource
func (h Handler) HandleUnfollow(e broker.Event) error {
	fmt.Printf("[HandleUnfollow] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var follow Follow
	if err := json.Unmarshal(e.Message().Body, &follow); err != nil {
		fmt.Println(err)
		return err
	}

	// Get the posts to remove from the feed
	var postUUIDs []string
	var postErr error

	if follow.Followee.Type == "User" {
		// Users can make posts on feeds other than their own, so we query by UserUUID and
		// not Feed, this will include posts they made on stocks for example.
		postUUIDs, postErr = h.postUUIDsMadeByUser(follow.Followee.UUID, 50)
	} else {
		postUUIDs, postErr = h.postUUIDsInFeed(follow.Followee.Type, follow.Followee.UUID)
	}

	if postErr != nil {
		fmt.Println(postErr)
		return postErr
	}

	// Remove the posts from the followers feed
	q := feeditems.BulkDeleteRequest{
		FeedType:  follow.Follower.Type,
		FeedUuid:  follow.Follower.UUID,
		PostUuids: postUUIDs,
	}

	if _, err := h.feedItemsSrv.BulkDelete(context.Background(), &q); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (h Handler) postUUIDsInFeed(FeedType, FeedUUID string) ([]string, error) {
	query := &posts.Feed{Type: FeedType, Uuid: FeedUUID}

	pRsp, err := h.postsSrv.ListFeed(context.Background(), query)
	if err != nil {
		return []string{}, err
	}

	uuids := make([]string, len(pRsp.Posts))
	for i, post := range pRsp.Posts {
		uuids[i] = post.Uuid
	}

	return uuids, nil
}

func (h Handler) postUUIDsMadeByUser(UUID string, limit int32) ([]string, error) {
	pRsp, err := h.postsSrv.ListUser(context.Background(), &posts.ListRequest{
		UserUuid: UUID, Limit: limit,
	})

	if err != nil {
		return []string{}, err
	}

	uuids := make([]string, len(pRsp.Posts))
	for i, post := range pRsp.Posts {
		uuids[i] = post.Uuid
	}

	return uuids, nil
}

// Returns the descriptions to attach to the resources posts, e.g.  "Because you follow John"
func (h Handler) descriptionForFollowee(followeeType, followeeUUID string) (string, error) {
	switch followeeType {
	case "Stock":
		stockRsp, err := h.stocksSrv.Get(context.Background(), &stocks.Stock{Uuid: followeeUUID})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Because you follow <&Stock:%v>%v<&/Stock>", stockRsp.Stock.Uuid, stockRsp.Stock.Name), nil
	case "User":
		user, err := h.usersSrv.Find(context.Background(), &users.User{Uuid: followeeUUID})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Because you follow <&User:%v>%v<&/User>", user.Uuid, user.FirstName), nil
	default:
		return "", nil
	}
}
