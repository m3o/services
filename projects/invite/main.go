package main

import (
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"

	"github.com/m3o/services/projects/invite/handler"
	pb "github.com/m3o/services/projects/invite/proto"
)

func main() {
	srv := service.New(
		service.Name("go.micro.service.projects.invite"),
		service.Version("latest"),
	)

	pb.RegisterInviteServiceHandler(handler.New(srv))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
