package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/services/account/invite/handler"
	pb "github.com/micro/services/account/invite/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.account.invite"),
		micro.Version("latest"),
	)

	service.Init()

	h := handler.NewHandler(service)
	pb.RegisterInviteHandler(service.Server(), h)

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
