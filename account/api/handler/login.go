package handler

import (
	"context"

	pb "github.com/micro/services/account/api/proto/account"
	login "github.com/micro/services/login/service/proto/login"
	users "github.com/micro/services/users/service/proto"
)

// EmailLogin looks up an account using an email and password
func (h *Handler) EmailLogin(ctx context.Context, req *pb.EmailLoginRequest, rsp *pb.EmailLoginResponse) error {
	// Verify the login credentials
	lRsp, err := h.login.VerifyLogin(ctx, &login.VerifyLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	// Lookup the user
	uRsp, err := h.users.Read(ctx, &users.ReadRequest{Id: lRsp.Id})
	if err != nil {
		return err
	}

	// Generate a token
	acc, err := h.auth.Generate(lRsp.Id)
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = acc.Token
	return nil
}

// EmailSignup creates an account using an email and password
func (h *Handler) EmailSignup(ctx context.Context, req *pb.EmailSignupRequest, rsp *pb.EmailSignupResponse) error {
	// Validate the user can be created
	_, err := h.users.Create(ctx, &users.CreateRequest{
		User:         &users.User{Email: req.Email},
		ValidateOnly: true,
	})
	if err != nil {
		return err
	}

	// Verify the login credentials
	_, err = h.login.CreateLogin(ctx, &login.CreateLoginRequest{
		Email:        req.Email,
		Password:     req.Password,
		ValidateOnly: true,
	})
	if err != nil {
		return err
	}

	// Create the user
	uRsp, err := h.users.Create(ctx, &users.CreateRequest{
		User: &users.User{Email: req.Email},
	})
	if err != nil {
		return err
	}

	// Create the login credentials
	_, err = h.login.CreateLogin(ctx, &login.CreateLoginRequest{
		Email:    req.Email,
		Password: req.Password,
		Id:       uRsp.User.Id,
	})
	if err != nil {
		h.users.Delete(ctx, &users.DeleteRequest{Id: uRsp.User.Id})
		return err
	}

	// Generate a token
	acc, err := h.auth.Generate(uRsp.User.Id)
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = acc.Token

	return nil
}
