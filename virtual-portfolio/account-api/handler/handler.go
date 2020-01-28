package handler

import (
	"context"
	"time"

	proto "github.com/kytra-app/account-api/proto"
	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/photos"
	posts "github.com/kytra-app/posts-srv/proto"
	user "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth      auth.Authenticator
	photos    photos.Service
	user      user.UsersService
	posts     posts.PostsService
	followers followers.FollowersService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, pics photos.Service, client client.Client) Handler {
	return Handler{
		auth:      auth,
		photos:    pics,
		user:      user.NewUsersService("kytra-srv-v1-users:8080", client),
		posts:     posts.NewPostsService("kytra-srv-v1-posts:8080", client),
		followers: followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
	}
}

// Health is used my k8s to healthcheck micro-api
func (h Handler) Health(ctx context.Context, req *proto.Healthcheck, rsp *proto.Healthcheck) error {
	return nil
}

// Get returns the account of the current user
func (h Handler) Get(ctx context.Context, req *proto.User, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	user, err := h.user.Find(ctx, &user.User{Uuid: u.UUID})
	if err != nil {
		return err
	}

	a := auth.User{UUID: user.Uuid}
	if rsp.Jwt, err = h.auth.EncodeUser(a); err != nil {
		return errors.InternalServerError("JWT_ERROR", "Error encoding JWT")
	}

	rsp.User = h.serializeUser(ctx, *user)
	return nil
}

// Update amends the account of the current user
func (h Handler) Update(ctx context.Context, req *proto.User, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	params := h.reverseSerializeUser(req)
	params.Uuid = u.UUID // Don't allow user to path another user by passing a UUID

	user, err := h.user.Update(ctx, params)
	if err != nil {
		return err
	}

	rsp.User = h.serializeUser(ctx, *user)
	return nil
}

// UpdatePassword amends the password of the current user. Requires current password.
func (h Handler) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordRequest, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	// Find the user
	usr, err := h.user.Find(ctx, &user.User{Uuid: u.UUID})
	if err != nil {
		return err
	}

	// Validate the current password
	credentials := &user.User{Email: usr.Email, Password: req.CurrentPassword}
	if _, err := h.user.ValidatePassword(ctx, credentials); err != nil {
		return err
	}

	// Update the password
	params := &user.User{Uuid: u.UUID, Password: req.NewPassword}
	if _, err = h.user.Update(ctx, params); err != nil {
		return err
	}

	return nil
}

// UpdateProfilePicture uploads the photo and then updates the users ProfilePictureId
func (h Handler) UpdateProfilePicture(ctx context.Context, req *proto.Picture, rsp *proto.Response) error {
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	// Upload the photo and handle any errors
	pictureID, err := h.photos.Upload(req.Base64)
	switch err {
	case photos.ErrInvalidData:
		return errors.BadRequest("INVALID_PHOTO", err.Error())
	case photos.ErrUnknown:
		return errors.InternalServerError("PHOTO_UPLOAD_ERROR", err.Error())
	}

	// Update the users ProfilePictureId
	params := &user.User{Uuid: u.UUID, ProfilePictureId: pictureID}
	user, err := h.user.Update(ctx, params)
	if err != nil {
		return err
	}

	rsp.User = h.serializeUser(ctx, *user)
	return nil
}

// Login allows a user to get their account details and auth token once their password is validated
func (h Handler) Login(ctx context.Context, req *proto.User, rsp *proto.Response) error {
	credentials := &user.User{Email: req.Email, Password: req.Password}
	user, err := h.user.ValidatePassword(ctx, credentials)

	if err != nil {
		return err
	}

	rsp.Jwt, err = h.auth.EncodeUser(auth.User{UUID: user.Uuid})
	if err != nil {
		return errors.InternalServerError("JWT_ERROR", "Error encoding JWT")
	}

	rsp.User = h.serializeUser(ctx, *user)
	return nil
}

func (h Handler) serializeUser(ctx context.Context, u user.User) *proto.User {
	out := &proto.User{
		Uuid:              u.Uuid,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Email:             u.Email,
		Username:          u.Username,
		ProfilePictureUrl: h.photos.GetURL(u.ProfilePictureId),
	}

	// Lookup Followers / Following Count. This is not critical information
	fContext, fCancel := context.WithTimeout(ctx, 150*time.Millisecond)
	fQuery := &followers.Resource{Type: "User", Uuid: u.Uuid}
	if res, err := h.followers.Count(fContext, fQuery); err == nil {
		out.FollowersCount = res.FollowerCount
		out.FollowingCount = res.FollowingCount
	}
	defer fCancel()

	// Lookup number of posts. This is not critical information
	pContext, pCancel := context.WithTimeout(ctx, 150*time.Millisecond)
	pQuery := &posts.Post{UserUuid: u.Uuid}
	if res, err := h.posts.Count(pContext, pQuery); err == nil {
		out.PostsCount = res.Count
	}
	defer pCancel()

	return out
}

func (h Handler) reverseSerializeUser(u *proto.User) *user.User {
	return &user.User{
		Uuid:      u.Uuid,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Username:  u.Username,
	}
}
