package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/users/service/handler"
	pb "github.com/micro/services/users/service/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.users"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	h, err := handler.NewHandler(service)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterUsersHandler(service.Server(), h)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
