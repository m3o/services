package main

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/auth/api/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.auth"),
		micro.Version("latest"),
	)

	service.Init()

	handler := &Handler{auth: service.Options().Auth}
	pb.RegisterAuthHandler(service.Server(), handler)

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

type Handler struct {
	auth auth.Auth
}

// Login exchanges auth credentials for a short lived token to be used when calling other apis
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	token, err := h.auth.Token(
		auth.WithCredentials(req.Id, req.Secret),
		auth.WithExpiry(time.Minute*5),
	)

	if err != nil {
		return err
	}

	rsp.Token = token.AccessToken
	return nil
}
