package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/client"

	proto "github.com/kytra-app/feed-api/proto"
	feeditems "github.com/kytra-app/feed-items-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/photos"
	"github.com/kytra-app/helpers/textenhancer"
	enhancer "github.com/kytra-app/post-enhancer-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
)

const resultsPerPage = 15

// Handler is an object can process RPC requests
type Handler struct {
	auth         auth.Authenticator
	photos       photos.Service
	feeditems    feeditems.FeedItemsService
	posts        posts.PostsService
	users        users.UsersService
	enhancer     enhancer.PostEnhancerService
	textenhancer textenhancer.Service
}

// New returns an instance of Handler
func New(auth auth.Authenticator, pics photos.Service, client client.Client) Handler {
	return Handler{
		auth:         auth,
		photos:       pics,
		textenhancer: textenhancer.Service{},
		users:        users.NewUsersService("kytra-srv-v1-users:8080", client),
		posts:        posts.NewPostsService("kytra-srv-v1-posts:8080", client),
		feeditems:    feeditems.NewFeedItemsService("kytra-srv-v1-feed-items:8080", client),
		enhancer:     enhancer.NewPostEnhancerService("kytra-srv-v1-post-enhancer:8080", client),
	}
}

// Get generates the users feed
func (h Handler) Get(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	// Get recent posts
	postUUIDs, postErr := h.getRecentPosts(req.Page)
	if postErr != nil {
		return postErr
	}

	// Enhance the posts
	postsRsp, err := h.enhancer.List(ctx, &enhancer.Request{PostUuids: postUUIDs})
	if err != nil {
		return err
	}

	rsp.Posts = make([]*proto.Post, len(postsRsp.Posts))
	for i, p := range postsRsp.Posts {
		rsp.Posts[i] = h.serializePost(p)
	}

	return nil
}

// GetAssetFeed finds the posts for an asset
func (h Handler) GetAssetFeed(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	if req.Page != 0 {
		return nil
	}

	// Fetch the recent posts for the stock's feed
	rawPostsRsp, err := h.posts.ListFeed(ctx, &posts.Feed{Type: "Stock", Uuid: req.Uuid})
	if err != nil {
		return nil
	}

	// Get the UUIDs from the posts
	postUUIDs := make([]string, len(rawPostsRsp.Posts))
	for i, p := range rawPostsRsp.Posts {
		postUUIDs[i] = p.Uuid
	}

	// Enhance the posts
	postsRsp, err := h.enhancer.List(ctx, &enhancer.Request{PostUuids: postUUIDs})
	if err != nil {
		return err
	}

	rsp.Posts = make([]*proto.Post, len(postsRsp.Posts))
	for i, p := range postsRsp.Posts {
		rsp.Posts[i] = h.serializePost(p)
	}

	return nil
}

// GetInvestorFeed finds the posts for an Investor
func (h Handler) GetInvestorFeed(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	if req.Page != 0 {
		return nil
	}

	user, err := h.users.Find(ctx, &users.User{Uuid: req.Uuid})
	if err != nil {
		return err
	}

	// Get the raw posts
	pRsp, err := h.posts.ListUser(ctx, &posts.ListRequest{UserUuid: user.Uuid})
	if err != nil {
		return err
	}

	// Get the post UUIDs
	uuids := make([]string, len(pRsp.Posts))
	for i, p := range pRsp.Posts {
		uuids[i] = p.Uuid
	}

	// Enhance the posts
	postsRsp, err := h.enhancer.List(ctx, &enhancer.Request{PostUuids: uuids})
	if err != nil {
		return err
	}

	rsp.Posts = make([]*proto.Post, len(postsRsp.Posts))
	for i, p := range postsRsp.Posts {
		rsp.Posts[i] = h.serializePost(p)
	}

	return nil
}

func (h Handler) serializePost(in *enhancer.EnhancedPost) *proto.Post {
	user := &proto.User{
		Uuid:              in.User.Uuid,
		FirstName:         in.User.FirstName,
		LastName:          in.User.LastName,
		Username:          in.User.Username,
		Following:         in.User.Following,
		ProfilePictureUrl: h.photos.GetURL(in.User.ProfilePictureId, 64, 64),
	}

	asset := &proto.Asset{}
	if in.Asset != nil {
		asset = &proto.Asset{
			Type:              in.Asset.Type,
			Uuid:              in.Asset.Uuid,
			Name:              in.Asset.Name,
			Color:             in.Asset.Color,
			Following:         in.Asset.Following,
			Description:       in.Asset.Description,
			ProfilePictureUrl: h.photos.GetURL(in.Asset.ProfilePictureId, 64, 64),
		}
	}

	comments := make([]*proto.Comment, len(in.Comments))
	for i, c := range in.Comments {
		comments[i] = &proto.Comment{
			Uuid:         c.Uuid,
			Text:         c.Text,
			EnhancedText: h.textenhancer.Enhance(c.Text),
			BullsCount:   c.BullsCount,
			BearsCount:   c.BearsCount,
			Opinion:      c.Opinion,
			User: &proto.User{
				Uuid:              c.User.Uuid,
				FirstName:         c.User.FirstName,
				LastName:          c.User.LastName,
				Username:          c.User.Username,
				ProfilePictureUrl: h.photos.GetURL(c.User.ProfilePictureId, 64, 64),
			},
		}
	}

	return &proto.Post{
		Uuid:                 in.Uuid,
		Text:                 in.Text,
		Title:                in.Title,
		BullsCount:           in.BullsCount,
		BearsCount:           in.BearsCount,
		Opinion:              in.Opinion,
		User:                 user,
		Comments:             comments,
		Asset:                asset,
		EnhancedText:         h.textenhancer.Enhance(in.Text),
		CreatedAt:            time.Unix(in.CreatedAt, 0).String(),
		AttachmentPictureUrl: h.photos.GetURL(in.AttachmentPictureId, 600, 600),
		AttachmentLinkUrl:    in.AttachmentLinkUrl,
	}
}

// Returns the UUIDs for the posts to be displayed, and their descriptions in a map
func (h Handler) getPostsForUser(UUID string, page int32) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	req := feeditems.GetFeedRequest{
		Type:  "User",
		Uuid:  UUID,
		Page:  page,
		Limit: resultsPerPage,
	}

	rsp, err := h.feeditems.GetFeed(ctx, &req)
	if err != nil {
		return []string{}, err
	}

	uuids := make([]string, len(rsp.Items))
	for i, p := range rsp.Items {
		uuids[i] = p.PostUuid
	}

	return uuids, nil
}

func (h Handler) getRecentPosts(page int32) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	rsp, err := h.posts.Recent(ctx, &posts.ListRequest{Page: page, Limit: resultsPerPage})
	if err != nil {
		return []string{}, err
	}

	uuids := make([]string, len(rsp.Posts))
	for i, p := range rsp.Posts {
		uuids[i] = p.Uuid
	}

	return uuids, nil
}
