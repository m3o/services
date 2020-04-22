package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/project/service/handler"
	pb "github.com/micro/services/project/service/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.project"),
		micro.Version("latest"),
	)

	service.Init()

	pb.RegisterProjectServiceHandler(service.Server(), handler.New(service))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
