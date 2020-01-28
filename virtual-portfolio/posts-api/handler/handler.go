package handler

import (
	bullbear "github.com/kytra-app/bullbear-srv/proto"
	comments "github.com/kytra-app/comments-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	photos "github.com/kytra-app/helpers/photos"
	"github.com/kytra-app/helpers/textenhancer"
	post "github.com/kytra-app/posts-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/client"
)

// Handler is an object can process RPC requests
type Handler struct {
	textenhancer textenhancer.Service
	photos       photos.Service
	auth         auth.Authenticator
	users        users.UsersService
	stocks       stocks.StocksService
	posts        post.PostsService
	bullBear     bullbear.BullBearService
	comments     comments.CommentsService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, pics photos.Service, client client.Client) Handler {
	return Handler{
		auth:         auth,
		photos:       pics,
		textenhancer: textenhancer.Service{},
		users:        users.NewUsersService("kytra-srv-v1-users:8080", client),
		stocks:       stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		posts:        posts.NewPostsService("kytra-srv-v1-posts:8080", client),
		bullBear:     bullbear.NewBullBearService("kytra-srv-v1-bullbear:8080", client),
		comments:     comments.NewCommentsService("kytra-srv-v1-comments:8080", client),
	}
}
