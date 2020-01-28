package handler

import (
	"context"
	"time"

	proto "github.com/kytra-app/posts-api/proto"

	bullbear "github.com/kytra-app/bullbear-srv/proto"
	comments "github.com/kytra-app/comments-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	posts "github.com/kytra-app/posts-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/errors"
)

// Create inserts a new post into the posts-srv
func (h Handler) Create(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("AUTH_REQUIRED", err.Error())
	}

	params, err := h.permittedParams(req, &u)
	if err != nil {
		return err
	}

	p, err := h.posts.Create(ctx, params)
	if err != nil {
		return err
	}
	if p.Error != nil {
		return err
	}

	rsp.Post, err = h.serializePost(p.Post, true)
	return err
}

// SetOpinion sets the bull/bear on a post
func (h Handler) SetOpinion(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("AUTH_REQUIRED", err.Error())
	}

	var opinion bullbear.Opinion
	switch req.Opinion {
	case "BEARISH":
		opinion = bullbear.Opinion_BEARISH
		break
	case "BULLISH":
		opinion = bullbear.Opinion_BULLISH
		break
	case "NONE":
		opinion = bullbear.Opinion_NONE
		break
	default:
		return errors.BadRequest("INVALID_OPINION", "An invalid opinion was provided")
	}

	bbReq := &bullbear.Request{
		Resource: &bullbear.Resource{Type: "Post", Uuid: req.Uuid},
		UserUuid: u.UUID,
		Opinion:  opinion,
	}

	_, err = h.bullBear.Create(ctx, bbReq)
	return err
}

// Get retrieves a post from the posts-srv
func (h Handler) Get(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	if len(req.Uuid) == 0 {
		rsp.Error = &proto.Error{Code: 400, Message: "Bad Request: Missing UUID"}
		return nil
	}

	p, err := h.posts.Get(ctx, &posts.Post{Uuid: req.Uuid})
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching post"}
		return err
	} else if p.Error != nil {
		rsp.Error = &proto.Error{Code: p.Error.Code, Message: p.Error.Message}
		return nil
	}

	rsp.Post, err = h.serializePost(p.Post, true)
	return err
}

// Delete destroys a post in the posts-srv
func (h Handler) Delete(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("AUTH_REQUIRED", err.Error())
	}

	p, err := h.posts.Get(ctx, &posts.Post{Uuid: req.Uuid})
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching post"}
		return err
	} else if p.Error != nil {
		rsp.Error = &proto.Error{Code: p.Error.Code, Message: p.Error.Message}
		return nil
	} else if p.Post.UserUuid != u.UUID {
		rsp.Error = &proto.Error{Code: 403}
		return nil
	}

	r, err := h.posts.Delete(ctx, &posts.Post{Uuid: req.Uuid})
	if err != nil || r.Error != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching post"}
		return err
	}

	return nil
}

// Update amends a post in the posts-srv
func (h Handler) Update(ctx context.Context, req *proto.Post, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("AUTH_REQUIRED", err.Error())
	}

	p, err := h.posts.Get(ctx, &posts.Post{Uuid: req.Uuid})
	if err != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching post"}
		return err
	} else if p.Error != nil {
		rsp.Error = &proto.Error{Code: p.Error.Code, Message: p.Error.Message}
		return nil
	} else if p.Post.UserUuid != u.UUID {
		rsp.Error = &proto.Error{Code: 403}
		return nil
	}

	r, err := h.posts.Update(ctx, &posts.Post{Uuid: req.Uuid, Text: req.Text})
	if err != nil || r.Error != nil {
		rsp.Error = &proto.Error{Code: 500, Message: "Error fetching post"}
		return err
	}

	rsp.Post, err = h.serializePost(r.Post, false)
	return err
}

func (h Handler) permittedParams(p *proto.Post, u *auth.User) (*posts.Post, error) {
	FeedType := "User"
	FeedUUID := u.UUID

	if p.FeedType == "Stock" && p.FeedUuid != "" {
		FeedType = p.FeedType
		FeedUUID = p.FeedUuid
	}

	var pictureID string
	var pictureErr error
	if p.AttachmentPictureBase64 != "" {
		pictureID, pictureErr = h.photos.Upload(p.AttachmentPictureBase64)
	}
	if pictureErr != nil {
		return &posts.Post{}, pictureErr
	}

	return &posts.Post{
		Text:                p.Text,
		Title:               p.Title,
		UserUuid:            u.UUID,
		FeedType:            FeedType,
		FeedUuid:            FeedUUID,
		AttachmentPictureId: pictureID,
		AttachmentLinkUrl:   p.AttachmentLinkUrl,
	}, nil
}

func (h Handler) serializePost(p *posts.Post, loadExtra bool) (*proto.Post, error) {
	rsp := &proto.Post{
		Uuid:                 p.Uuid,
		Text:                 p.Text,
		Title:                p.Title,
		EnhancedText:         h.textenhancer.Enhance(p.Text),
		FeedType:             p.FeedType,
		FeedUuid:             p.FeedUuid,
		User:                 &proto.User{Uuid: p.UserUuid},
		CreatedAt:            time.Unix(p.CreatedAt, 0).String(),
		AttachmentPictureUrl: h.photos.GetURL(p.AttachmentPictureId),
		AttachmentLinkUrl:    p.AttachmentLinkUrl,
	}

	// Add the user who created the post
	usrs, err := h.fetchUsers([]string{p.UserUuid})
	if err != nil {
		return rsp, err
	}
	rsp.User = usrs[p.UserUuid]

	if !loadExtra {
		return rsp, nil
	}

	// Add the bull & bear counts
	if bulls, bears, err := h.bullsAndBearsForPost(p); err == nil {
		rsp.BullsCount = bulls
		rsp.BearsCount = bears
	}

	// Add the comments to the post
	if comments, err := h.getCommentsForPost(p); err == nil {
		rsp.Comments = comments
	}

	// Add the asset to the post
	if asset, err := h.getAssetForPost(p); err == nil {
		rsp.Asset = asset
	}

	return rsp, nil
}

func (h Handler) bullsAndBearsForPost(p *posts.Post) (int32, int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	rsp, err := h.bullBear.Get(ctx, &bullbear.Resource{Type: "Post", Uuid: p.Uuid})
	if err != nil {
		return 0, 0, err
	}

	return rsp.Resource.BullsCount, rsp.Resource.BearsCount, nil
}

// fetchUsers returns a list of users, grouped by their UUIDs
func (h Handler) fetchUsers(uuids []string) (map[string]*proto.User, error) {
	rsp := make(map[string]*proto.User, len(uuids))

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	usersRsp, err := h.users.List(ctx, &users.ListRequest{Uuids: uuids})
	if err != nil {
		return rsp, err
	}

	for _, u := range usersRsp.Users {
		rsp[u.Uuid] = &proto.User{
			Uuid:              u.Uuid,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			Username:          u.Username,
			ProfilePictureUrl: h.photos.GetURL(u.ProfilePictureId),
		}
	}

	return rsp, nil
}

func (h Handler) getCommentsForPost(p *posts.Post) ([]*proto.Comment, error) {
	// Initialize the response
	var result []*proto.Comment

	// Create a context for the request
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// Get the comments
	rsp, err := h.comments.GetResource(ctx, &comments.Resource{Type: "Post", Uuid: p.Uuid})
	if err != nil {
		return result, err
	}

	// Get the users who made those comments
	userUuids := make([]string, len(rsp.Resource.Comments))
	for i, c := range rsp.Resource.Comments {
		userUuids[i] = c.UserUuid
	}
	usrs, err := h.fetchUsers(userUuids)
	if err != nil {
		return result, err
	}

	// Serialize the data
	result = make([]*proto.Comment, len(rsp.Resource.Comments))
	for i, c := range rsp.Resource.Comments {
		result[i] = &proto.Comment{
			Uuid:         c.Uuid,
			Text:         c.Text,
			User:         usrs[c.UserUuid],
			EnhancedText: h.textenhancer.Enhance(c.Text),
		}
	}

	return result, nil
}

func (h Handler) getAssetForPost(p *posts.Post) (*proto.Asset, error) {
	if p.FeedType != "Stock" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	sRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: p.FeedUuid})
	if err != nil {
		return nil, err
	}

	return &proto.Asset{
		Type:              "Stock",
		Uuid:              sRsp.Stock.Uuid,
		Name:              sRsp.Stock.Name,
		ProfilePictureUrl: h.photos.GetURL(sRsp.Stock.ProfilePictureId),
	}, nil
}
