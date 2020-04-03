package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/client"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/services/account/api/proto/account"
	users "github.com/micro/services/users/service/proto"
)

// Token generates a new JWT using a RefreshToken token
func (h *Handler) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	tok, err := h.auth.Token(auth.WithExpiry(time.Hour*24), auth.WithToken(req.RefreshToken))
	if err != nil {
		return err
	}
	rsp.Token = serializeToken(tok)
	return nil
}

// Login looks up an account using an email and password
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	// Generate a token
	tok, err := h.auth.Token(auth.WithCredentials(req.Email, req.Password))
	if err != nil {
		return err
	}

	// Lookup the user
	uRsp, err := h.users.Read(ctx, &users.ReadRequest{Email: req.Email}, client.WithServiceToken())
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = serializeToken(tok)
	return nil
}

// Signup creates an account using an email and password
func (h *Handler) Signup(ctx context.Context, req *pb.SignupRequest, rsp *pb.SignupResponse) error {
	// Validate the user can be created
	_, err := h.users.Create(ctx, &users.CreateRequest{
		User:         &users.User{Email: req.Email},
		ValidateOnly: true,
	}, client.WithServiceToken())
	if err != nil {
		return err
	}

	// Generate an account and token
	acc, err := h.auth.Generate(req.Email, auth.WithRoles("user"), auth.WithSecret(req.Password))
	if err != nil {
		return err
	}
	tok, err := h.auth.Token(auth.WithCredentials(acc.ID, acc.Secret), auth.WithExpiry(time.Hour*24))
	if err != nil {
		return err
	}

	// Create the user
	uRsp, err := h.users.Create(ctx, &users.CreateRequest{User: &users.User{Email: req.Email}}, client.WithServiceToken())
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = serializeToken(tok)
	return nil
}
