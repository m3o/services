package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/home/api/handler"
	pb "github.com/micro/services/home/api/proto/home"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.home"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	h := handler.NewHandler(service)
	pb.RegisterHomeHandler(service.Server(), h)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
