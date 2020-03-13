package main

import (
	"github.com/micro/services/apps/service/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	pb "github.com/micro/services/apps/service/proto/apps"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.apps"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	h := handler.NewHandler(service)
	pb.RegisterAppsHandler(service.Server(), h)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
