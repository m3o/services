package handler

import (
	"context"

	bullbear "github.com/kytra-app/bullbear-srv/proto"
	comments "github.com/kytra-app/comments-srv/proto"
	proto "github.com/kytra-app/posts-api/proto"
	"github.com/micro/go-micro/errors"
)

// CreateComment inserts a new comment on the resource
func (h Handler) CreateComment(ctx context.Context, req *proto.Comment, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		rsp.Error = &proto.Error{Code: 401, Message: err.Error()}
		return nil
	}

	comment := &comments.Comment{
		UserUuid: u.UUID,
		Text:     req.Text,
		Resource: &comments.Resource{
			Type: "Post",
			Uuid: req.Post.Uuid,
		},
	}

	cRsp, err := h.comments.Create(ctx, comment)
	if err != nil {
		return err
	}

	rsp.Comment = &proto.Comment{
		Uuid:         cRsp.Comment.Uuid,
		Text:         cRsp.Comment.Text,
		EnhancedText: h.textenhancer.Enhance(cRsp.Comment.Text),
	}
	return nil
}

// DeleteComment deletes a comment from a resource
func (h Handler) DeleteComment(ctx context.Context, req *proto.Comment, rsp *proto.Response) error {
	if req.Uuid == "" {
		rsp.Error = &proto.Error{Code: 400, Message: "Missing UUID"}
		return nil
	}

	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		rsp.Error = &proto.Error{Code: 401, Message: err.Error()}
		return nil
	}

	c, err := h.comments.Get(ctx, &comments.Comment{Uuid: req.Uuid})
	if err != nil {
		return err
	}
	// Ensure user has permission to delete comment
	if c.Comment.UserUuid != u.UUID {
		rsp.Error = &proto.Error{Code: 403}
		return nil
	}

	_, err = h.comments.Delete(ctx, &comments.Comment{Uuid: req.Uuid})
	if err != nil {
		return err
	}

	return nil
}

// SetCommentOpinion sets the bull/bear on a comment
func (h Handler) SetCommentOpinion(ctx context.Context, req *proto.Comment, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		rsp.Error = &proto.Error{Code: 401, Message: err.Error()}
		return nil
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
		Resource: &bullbear.Resource{Type: "Comment", Uuid: req.Uuid},
		UserUuid: u.UUID,
		Opinion:  opinion,
	}

	_, err = h.bullBear.Create(ctx, bbReq)
	return err
}
