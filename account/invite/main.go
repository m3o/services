package main

import (
	"github.com/m3o/services/account/invite/handler"
	pb "github.com/m3o/services/account/invite/proto"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
)

func main() {
	srv := service.New(
		service.Name("go.micro.service.account.invite"),
		service.Version("latest"),
	)

	srv.Init()

	h := handler.NewHandler(srv)
	pb.RegisterInviteHandler(h)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
