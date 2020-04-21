package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/teams/invites/handler"
	pb "github.com/micro/services/teams/invites/proto/invites"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.teams.invites"),
		micro.Version("latest"),
	)

	service.Init()

	pb.RegisterInvitesHandler(service.Server(), handler.New(service))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
