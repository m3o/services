package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/apps/api/handler"
	pb "github.com/micro/services/apps/api/proto/apps"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.apps"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register the handler
	h := handler.NewHandler(service)
	pb.RegisterAppsHandler(service.Server(), h)

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
