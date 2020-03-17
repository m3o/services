package main

import (
	"serverless/handler"
	serverless "serverless/proto/serverless"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	pb "github.com/micro/go-micro/v2/runtime/service/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.serverless"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	serverless.RegisterAppsHandler(service.Server(), &handler.Apps{
		Client: pb.NewRuntimeService("go.micro.runtime", service.Client()),
	})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
