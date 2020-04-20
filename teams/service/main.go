package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/teams/service/handler"
	pb "github.com/micro/services/teams/service/proto/teams"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.teams"),
		micro.Version("latest"),
	)

	service.Init()

	pb.RegisterTeamsHandler(service.Server(), handler.New(service))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
