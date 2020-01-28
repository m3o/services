package handler

import (
	"context"
	"fmt"

	"github.com/dongri/phonenumber"

	sms "github.com/kytra-app/sms-verification-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	proto "github.com/kytra-app/welcome-api/proto"
)

// RequestSMSCode sends an auth code to the user via SMS
func (h Handler) RequestSMSCode(ctx context.Context, req *proto.AuthRequest, rsp *proto.AuthResponse) error {
	// TEST USER FOR APPLE
	fmt.Println(req.PhoneNumber)
	if req.PhoneNumber == "+441111111111" {
		return nil
	}

	number := phonenumber.Parse(req.PhoneNumber, "UK")
	_, err := h.smsVer.Request(ctx, &sms.Verification{PhoneNumber: number})
	return err
}

// ValidateSMSCode checks the auth code submitted by the user
func (h Handler) ValidateSMSCode(ctx context.Context, req *proto.AuthRequest, rsp *proto.AuthResponse) error {
	// TEST USER FOR APPLE
	if req.PhoneNumber == "+441111111111" {
		user, _ := h.user.Find(ctx, &users.User{PhoneNumber: "441111111111"})
		rsp.User, _ = h.serializeUser(user)
		return nil
	}

	number := phonenumber.Parse(req.PhoneNumber, "UK")
	ver, err := h.smsVer.Verify(ctx, &sms.Verification{
		PhoneNumber: number,
		Code:        req.Code,
	})

	if err != nil {
		return err
	}

	// Find a user with the phone number
	user, err := h.user.Find(ctx, &users.User{PhoneNumber: ver.PhoneNumber})
	if err != nil {
		rsp.VerificationUuid = ver.Uuid
		return nil
	}

	rsp.User, err = h.serializeUser(user)
	return err
}
