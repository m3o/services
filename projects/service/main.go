package main

import (
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"

	"github.com/m3o/services/projects/service/handler"
	pb "github.com/m3o/services/projects/service/proto"
)

func main() {
	srv := service.New(
		service.Name("go.micro.service.projects"),
		service.Version("latest"),
	)

	pb.RegisterProjectsHandler(handler.New(srv))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
