package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/services/account/api/proto/account"
	login "github.com/micro/services/login/service/proto/login"
	users "github.com/micro/services/users/service/proto"
)

// Token generates a new JWT using a secret token
func (h *Handler) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	tok, err := h.auth.Token(req.Id, req.Secret, auth.WithTokenExpiry(time.Hour*24))
	if err != nil {
		return err
	}
	rsp.Token = serializeToken(tok)
	return nil
}

// Login looks up an account using an email and password
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	// Generate a context with elevated privelages
	privCtx := auth.ContextWithToken(ctx, h.authToken)

	// Verify the login credentials
	lRsp, err := h.login.VerifyLogin(privCtx, &login.VerifyLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	// Lookup the user
	uRsp, err := h.users.Read(privCtx, &users.ReadRequest{Id: lRsp.Id})
	if err != nil {
		return err
	}

	// Generate an auth account and token
	acc, err := h.auth.Generate(lRsp.Id, auth.WithRoles("user"))
	if err != nil {
		return err
	}
	tok, err := h.auth.Token(acc.ID, acc.Secret, auth.WithTokenExpiry(time.Hour*24))
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = serializeToken(tok)
	rsp.Secret = acc.Secret
	return nil
}

// Signup creates an account using an email and password
func (h *Handler) Signup(ctx context.Context, req *pb.SignupRequest, rsp *pb.SignupResponse) error {
	// Generate a context with elevated privelages
	privCtx := auth.ContextWithToken(ctx, h.authToken)

	// Validate the user can be created
	_, err := h.users.Create(privCtx, &users.CreateRequest{
		User:         &users.User{Email: req.Email},
		ValidateOnly: true,
	})
	if err != nil {
		return err
	}

	// Verify the login credentials
	_, err = h.login.CreateLogin(privCtx, &login.CreateLoginRequest{
		Email:        req.Email,
		Password:     req.Password,
		ValidateOnly: true,
	})
	if err != nil {
		return err
	}

	// Create the user
	uRsp, err := h.users.Create(privCtx, &users.CreateRequest{
		User: &users.User{Email: req.Email},
	})
	if err != nil {
		return err
	}

	// Create the login credentials
	_, err = h.login.CreateLogin(privCtx, &login.CreateLoginRequest{
		Email:    req.Email,
		Password: req.Password,
		Id:       uRsp.User.Id,
	})
	if err != nil {
		h.users.Delete(privCtx, &users.DeleteRequest{Id: uRsp.User.Id})
		return err
	}

	// Generate an account and token
	acc, err := h.auth.Generate(uRsp.User.Id, auth.WithRoles("user"))
	if err != nil {
		return err
	}
	tok, err := h.auth.Token(acc.ID, acc.Secret, auth.WithTokenExpiry(time.Hour*24))
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = serializeToken(tok)
	rsp.Secret = acc.Secret
	return nil
}
