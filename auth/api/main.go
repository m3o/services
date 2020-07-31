package main

import (
	"context"
	"time"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	mauth "github.com/micro/micro/v3/service/auth"

	pb "github.com/m3o/services/auth/api/proto"
)

func main() {
	srv := service.New(
		service.Name("go.micro.api.auth"),
		service.Version("latest"),
	)

	pb.RegisterAuthHandler(new(Handler))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

type Handler struct {
}

// Login exchanges auth credentials for a short lived token to be used when calling other apis
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	token, err := mauth.Token(
		auth.WithCredentials(req.Id, req.Secret),
		auth.WithExpiry(time.Minute*5),
	)

	if err != nil {
		return err
	}

	rsp.Token = token.AccessToken
	return nil
}
