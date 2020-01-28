package handler

import (
	"context"

	proto "github.com/kytra-app/push-notifications-srv/proto"
	"github.com/kytra-app/push-notifications-srv/storage"
	"github.com/micro/go-micro/errors"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

// New returns an instance of Handler
func New(storage storage.Service) *Handler {
	return &Handler{
		db:   storage,
		expo: expo.NewPushClient(nil),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	db   storage.Service
	expo *expo.PushClient
}

// RegisterToken stores a push notification token in the DB
func (h Handler) RegisterToken(ctx context.Context, req *proto.Token, rsp *proto.Token) error {
	_, err := h.db.CreateToken(req.Token, req.UserUuid)
	return err
}

// SendNotification sends a notification to a user
func (h Handler) SendNotification(ctx context.Context, req *proto.Notification, rsp *proto.Notification) error {
	// Find the users notification token
	token, err := h.db.GetToken(req.UserUuid)
	if err != nil {
		return err
	}

	// Generate the expo token
	pushToken, err := expo.NewExponentPushToken(token.Token)
	if err != nil {
		return errors.BadRequest("INVALID_TOKEN", "An invalid token was sent to Expo")
	}

	// Publish message
	_, err = h.expo.Publish(&expo.PushMessage{
		To:    pushToken,
		Body:  req.Subtitle,
		Title: req.Title,
	})

	// Check errors
	if err != nil {
		return errors.InternalServerError("EXPO_ERROR", err.Error())
	}

	return nil
}
