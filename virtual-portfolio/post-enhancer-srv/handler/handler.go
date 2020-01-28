package handler

import (
	"context"

	auth "github.com/kytra-app/helpers/authentication"
	unique "github.com/kytra-app/helpers/unique"
	stocks "github.com/kytra-app/stocks-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"

	bullbear "github.com/kytra-app/bullbear-srv/proto"
	comments "github.com/kytra-app/comments-srv/proto"
	followers "github.com/kytra-app/followers-srv/proto"
	proto "github.com/kytra-app/post-enhancer-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth      auth.Authenticator
	users     users.UsersService
	stocks    stocks.StocksService
	posts     posts.PostsService
	comments  comments.CommentsService
	bullbear  bullbear.BullBearService
	followers followers.FollowersService
}

// feedData is an object containing all the data needed to render a feed
type feedData struct {
	posts           []*posts.Post
	users           []*users.User
	stocks          []*stocks.Stock
	comments        []*comments.Resource
	bullbears       []*bullbear.Resource
	userFollowings  map[string]bool
	stockFollowings map[string]bool
}

// New returns an instance of Handler
func New(auth auth.Authenticator, client client.Client) Handler {
	return Handler{
		auth:      auth,
		users:     users.NewUsersService("kytra-srv-v1-users:8080", client),
		posts:     posts.NewPostsService("kytra-srv-v1-posts:8080", client),
		stocks:    stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		bullbear:  bullbear.NewBullBearService("kytra-srv-v1-bullbear:8080", client),
		comments:  comments.NewCommentsService("kytra-srv-v1-comments:8080", client),
		followers: followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
	}
}

// List generates enhanced posts using the UUIDs requested
func (h Handler) List(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	// Fetch the posts
	posts, err := h.fetchPosts(ctx, req.PostUuids)
	if err != nil {
		return err
	}

	// Fetch the comments
	comments, err := h.fetchCommentsForPosts(ctx, posts)
	if err != nil {
		return err
	}

	// Fetch the users
	users, err := h.fetchUsersForPostsAndComments(ctx, posts, comments)
	if err != nil {
		return err
	}

	// Fetch the followings for the users (which do I follow)
	userFollowings, err := h.fetchFollowingsForUsers(ctx, users)
	if err != nil {
		return err
	}

	// Fetch the stocks
	stocks, err := h.fetchStocksForPosts(ctx, posts)
	if err != nil {
		return err
	}

	// Fetch the followings for the stocks (which do I follow)
	stockFollowings, err := h.fetchFollowingsForStocks(ctx, stocks)
	if err != nil {
		return err
	}

	// Fetch the bulls and bears for the posts
	postBullBears, err := h.fetchBullsAndBearsForPosts(ctx, posts)
	if err != nil {
		return err
	}

	// Fetch the bulls and bears for the comments
	commentBullBears, err := h.fetchBullsAndBearsForComments(ctx, comments)
	if err != nil {
		return err
	}

	// Concat the bulls and bears
	bullbears := append(postBullBears, commentBullBears...)

	// Initialize the response array to the correct length
	rsp.Posts = h.serializeFeed(feedData{
		posts, users, stocks, comments, bullbears, userFollowings, stockFollowings,
	})

	return nil
}

func (h Handler) fetchCommentsForPosts(ctx context.Context, posts []*posts.Post) ([]*comments.Resource, error) {
	// Get the Post UUIDs
	uuids := make([]string, len(posts))
	for i, post := range posts {
		uuids[i] = post.Uuid
	}

	// Request the comments and return the result
	query := &comments.ListRequest{ResourceType: "Post", ResourceUuids: uuids}
	rsp, err := h.comments.ListResources(ctx, query)
	if err != nil {
		return []*comments.Resource{}, err
	}

	return rsp.Resources, nil
}

func (h Handler) fetchBullsAndBearsForPosts(ctx context.Context, posts []*posts.Post) ([]*bullbear.Resource, error) {
	// Get the Post UUIDs
	uuids := make([]string, len(posts))
	for i, post := range posts {
		uuids[i] = post.Uuid
	}

	// Get the user UUID
	var userUUID string
	if u, err := h.auth.UserFromContext(ctx); err == nil {
		userUUID = u.UUID
	}

	// Request the bulls&bears and return the result
	q := &bullbear.ListRequest{ResourceType: "Post", ResourceUuids: uuids, UserUuid: userUUID}
	rsp, err := h.bullbear.List(ctx, q)
	if err != nil {
		return []*bullbear.Resource{}, err
	}

	return rsp.Resources, nil
}

func (h Handler) fetchBullsAndBearsForComments(ctx context.Context, commentResources []*comments.Resource) ([]*bullbear.Resource, error) {
	// Get the Comment UUIDs
	var uuids []string
	for _, resource := range commentResources {
		for _, comment := range resource.Comments {
			uuids = append(uuids, comment.Uuid)
		}
	}

	// Get the user UUID
	var userUUID string
	if u, err := h.auth.UserFromContext(ctx); err == nil {
		userUUID = u.UUID
	}

	// Request the bulls&bears and return the result
	q := &bullbear.ListRequest{ResourceType: "Comment", ResourceUuids: uuids, UserUuid: userUUID}
	rsp, err := h.bullbear.List(ctx, q)
	if err != nil {
		return []*bullbear.Resource{}, err
	}

	return rsp.Resources, nil
}

func (h Handler) fetchPosts(ctx context.Context, uuids []string) ([]*posts.Post, error) {
	// Request the posts and return the result
	rsp, err := h.posts.List(ctx, &posts.ListRequest{Uuids: uuids})
	if err != nil {
		return []*posts.Post{}, nil
	}

	return rsp.Posts, err
}

func (h Handler) fetchStocksForPosts(ctx context.Context, posts []*posts.Post) ([]*stocks.Stock, error) {
	// Get the posting user UUIDs
	var uuids []string
	for _, post := range posts {
		if post.FeedType == "Stock" {
			uuids = append(uuids, post.FeedUuid)
		}
	}

	// Ensure we only request unique stocks
	uuids = unique.Strings(uuids)

	// Request the stocks and return the result
	rsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: uuids})
	if err != nil {
		return []*stocks.Stock{}, err
	}

	return rsp.Stocks, err
}

func (h Handler) fetchUsersForPostsAndComments(ctx context.Context, posts []*posts.Post, comments []*comments.Resource) ([]*users.User, error) {
	// Get the posting user UUIDs
	uuids := make([]string, len(posts))
	for i, post := range posts {
		uuids[i] = post.UserUuid
	}

	// Get the commenting user UUIDs
	for _, post := range comments {
		for _, c := range post.Comments {
			uuids = append(uuids, c.UserUuid)
		}
	}

	// Ensure we only request unique users
	uuids = unique.Strings(uuids)

	// Request the users and return the result
	rsp, err := h.users.List(ctx, &users.ListRequest{Uuids: uuids})
	if err != nil {
		return []*users.User{}, err
	}

	return rsp.Users, nil
}

func (h Handler) fetchFollowingsForStocks(ctx context.Context, stocks []*stocks.Stock) (map[string]bool, error) {
	// Setup response, default to not following users
	rsp := make(map[string]bool, len(stocks))
	for _, u := range stocks {
		rsp[u.Uuid] = false
	}

	// Try and get the user from the context
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return rsp, nil
	}

	// Get the user UUIDS
	uuids := make([]string, len(stocks))
	for i, user := range stocks {
		uuids[i] = user.Uuid
	}

	// Construct the query
	req := &followers.ListRequest{
		Follower:      &followers.Resource{Uuid: u.UUID, Type: "User"},
		FolloweeType:  "Stock",
		FolloweeUuids: uuids,
	}

	// Request the data
	data, err := h.followers.List(ctx, req)
	if err != nil {
		return rsp, err
	}

	// Update the response
	for _, r := range data.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}

func (h Handler) fetchFollowingsForUsers(ctx context.Context, users []*users.User) (map[string]bool, error) {
	// Setup response, default to not following users
	rsp := make(map[string]bool, len(users))
	for _, u := range users {
		rsp[u.Uuid] = false
	}

	// Try and get the user from the context
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return rsp, nil
	}

	// Get the user UUIDS
	uuids := make([]string, len(users))
	for i, user := range users {
		uuids[i] = user.Uuid
	}

	// Construct the query
	req := &followers.ListRequest{
		Follower:      &followers.Resource{Uuid: u.UUID, Type: "User"},
		FolloweeType:  "User",
		FolloweeUuids: uuids,
	}

	// Request the data
	data, err := h.followers.List(ctx, req)
	if err != nil {
		return rsp, err
	}

	// Update the response
	for _, r := range data.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}

func (h Handler) serializeFeed(data feedData) []*proto.EnhancedPost {
	// usersMap is a map of users, indexed by their UUID
	usersMap := make(map[string]*users.User, len(data.users))
	for _, user := range data.users {
		usersMap[user.Uuid] = user
	}

	// stocksMap is a map of stocks, indexed by their UUID
	stocksMap := make(map[string]*stocks.Stock, len(data.stocks))
	for _, stock := range data.stocks {
		stocksMap[stock.Uuid] = stock
	}

	// bullbearsPostMap is a map of bullbears resources, indexed by PostUUID
	bullbearsPostMap := make(map[string]*bullbear.Resource, len(data.posts))
	for _, r := range data.bullbears {
		if r.Type == "Post" {
			bullbearsPostMap[r.Uuid] = r
		}
	}

	// bullbearsCommentMap is a map of bullbears resources, indexed by CommentUUID
	bullbearsCommentMap := make(map[string]*bullbear.Resource, len(data.comments))
	for _, r := range data.bullbears {
		if r.Type == "Comment" {
			bullbearsCommentMap[r.Uuid] = r
		}
	}

	// commentResourceMap is a map of comment resources, indexed by PostUUID
	commentResourceMap := make(map[string]*comments.Resource, len(data.posts))
	for _, r := range data.comments {
		commentResourceMap[r.Uuid] = r
	}

	// Serialize the data
	rsp := make([]*proto.EnhancedPost, len(data.posts))
	for i, post := range data.posts {
		bb := bullbearsPostMap[post.Uuid]
		cr := commentResourceMap[post.Uuid]

		rsp[i] = &proto.EnhancedPost{
			Uuid:                post.Uuid,
			Text:                post.Text,
			Title:               post.Title,
			Comments:            h.serializeComments(cr.Comments, usersMap, bullbearsCommentMap),
			BullsCount:          bb.BullsCount,
			BearsCount:          bb.BearsCount,
			Opinion:             bb.Opinion.String(),
			CreatedAt:           post.CreatedAt,
			AttachmentPictureId: post.AttachmentPictureId,
			AttachmentLinkUrl:   post.AttachmentLinkUrl,
		}

		// Serialize the user who made the post
		if user, ok := usersMap[post.UserUuid]; ok {
			rsp[i].User = &proto.User{
				Uuid:             user.Uuid,
				FirstName:        user.FirstName,
				LastName:         user.LastName,
				Username:         user.Username,
				ProfilePictureId: user.ProfilePictureId,
				Following:        data.userFollowings[user.Uuid],
			}
		}

		// Serialize the feed the user posted on
		if post.FeedType == "Stock" {
			stock, ok := stocksMap[post.FeedUuid]
			if !ok {
				continue
			}

			rsp[i].Asset = &proto.Asset{
				Type:             "Stock",
				Uuid:             stock.Uuid,
				Name:             stock.Name,
				Color:            stock.Color,
				Description:      stock.Description,
				ProfilePictureId: stock.ProfilePictureId,
				Following:        data.stockFollowings[stock.Uuid],
			}
		}
	}

	return rsp
}

func (h Handler) serializeComments(comments []*comments.Comment, usersMap map[string]*users.User, bullbearsMap map[string]*bullbear.Resource) []*proto.Comment {
	res := make([]*proto.Comment, len(comments))

	for i, c := range comments {
		bb := bullbearsMap[c.Uuid]

		res[i] = &proto.Comment{
			Uuid:       c.Uuid,
			Text:       c.Text,
			BullsCount: bb.BullsCount,
			BearsCount: bb.BearsCount,
			Opinion:    bb.Opinion.String(),
		}

		if user, ok := usersMap[c.UserUuid]; ok {
			res[i].User = &proto.User{
				Uuid:             user.Uuid,
				FirstName:        user.FirstName,
				LastName:         user.LastName,
				Username:         user.Username,
				ProfilePictureId: user.ProfilePictureId,
			}
		}
	}

	return res
}
