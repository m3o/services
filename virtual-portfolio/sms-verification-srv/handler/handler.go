package handler

import (
	"context"
	"fmt"

	"github.com/kytra-app/helpers/sms"
	proto "github.com/kytra-app/sms-verification-srv/proto"
	"github.com/kytra-app/sms-verification-srv/storage"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(storage storage.Service, sms sms.Service) *Handler {
	return &Handler{
		db:  storage,
		sms: sms,
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	db  storage.Service
	sms sms.Service
}

// Request triggers an SMS verification code to be sent
func (h Handler) Request(ctx context.Context, req *proto.Verification, rsp *proto.Verification) error {
	ver, err := h.db.Request(req.PhoneNumber)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Your Kytra verification code is %v", ver.Code)
	if err = h.sms.Send(ver.PhoneNumber, msg); err != nil {
		fmt.Println(err)
		return storage.ErrInvalidNumber
	}

	*rsp = proto.Verification{
		Uuid:        ver.UUID,
		PhoneNumber: ver.PhoneNumber,
		Expired:     ver.Expired(),
		Verified:    ver.Verified,
	}

	return nil
}

// Verify attempts to verify the users verification code
func (h Handler) Verify(ctx context.Context, req *proto.Verification, rsp *proto.Verification) error {
	if req.PhoneNumber == "" {
		return errors.BadRequest("INVALID_NUMBER", "The phone number was not provided")
	}

	if req.Code == "" {
		return errors.BadRequest("INVALID_CODE", "The code was not provided")
	}

	ver, err := h.db.Verify(req.PhoneNumber, req.Code)
	if err != nil {
		return err
	}

	*rsp = proto.Verification{
		Uuid:        ver.UUID,
		PhoneNumber: ver.PhoneNumber,
		Expired:     ver.Expired(),
		Verified:    ver.Verified,
	}

	return nil
}

// Get attempts to find the verification code by UUIDD
func (h Handler) Get(ctx context.Context, req *proto.Verification, rsp *proto.Verification) error {
	if req.Uuid == "" {
		return errors.BadRequest("INVALID_UUID", "The UUID was not provided")
	}

	ver, err := h.db.Get(req.Uuid)
	if err != nil {
		return err
	}

	*rsp = proto.Verification{
		Uuid:        ver.UUID,
		PhoneNumber: ver.PhoneNumber,
		Verified:    ver.Verified,
		Expired:     ver.Expired(),
	}

	return nil
}
