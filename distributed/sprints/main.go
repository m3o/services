package main

import (
	"github.com/micro/services/distributed/sprints/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/distributed/sprints/proto/sprints"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.distributed.sprints"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	h := handler.NewHandler(service)
	pb.RegisterSprintsHandler(service.Server(), h)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
