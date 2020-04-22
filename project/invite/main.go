package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/project/invite/handler"
	pb "github.com/micro/services/project/invite/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.project.invite"),
		micro.Version("latest"),
	)

	service.Init()

	pb.RegisterInviteServiceHandler(service.Server(), handler.New(service))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
