package consumer

import (
	"context"

	comments "github.com/kytra-app/comments-srv/proto"
	followers "github.com/kytra-app/followers-srv/proto"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/helpers/textenhancer"
	"github.com/kytra-app/notifications-srv/storage"
	posts "github.com/kytra-app/posts-srv/proto"
	push "github.com/kytra-app/push-notifications-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
)

// New returns an instance of Consumer
func New(client client.Client, db storage.Service) Consumer {
	return Consumer{
		db:           db,
		textenhancer: textenhancer.Service{},
		usersSrv:     users.NewUsersService("kytra-srv-v1-users:8080", client),
		postsSrv:     posts.NewPostsService("kytra-srv-v1-posts:8080", client),
		followersSrv: followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		commentsSrv:  comments.NewCommentsService("kytra-srv-v1-comments:8080", client),
		pushSrv:      push.NewPushNotificationsService("kytra-srv-v1-push-notifications:8080", client),
	}
}

// Consumer is an object which processes various messages
type Consumer struct {
	db           storage.Service
	textenhancer textenhancer.Service
	usersSrv     users.UsersService
	postsSrv     posts.PostsService
	commentsSrv  comments.CommentsService
	followersSrv followers.FollowersService
	pushSrv      push.PushNotificationsService
}

func (c Consumer) usersTaggedInText(text string) []string {
	usernames := c.textenhancer.ListTaggedUsers(text)

	result := make([]string, len(usernames))
	for i, username := range usernames {
		user, err := c.usersSrv.Find(context.Background(), &users.User{Username: username})

		switch err {
		case nil:
			result[i] = user.Uuid
		case microgorm.ErrNotFound:
			continue
		default:
			return result
		}
	}

	return result
}
