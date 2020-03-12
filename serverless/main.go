package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"serverless/handler"
	"serverless/subscriber"

	serverless "serverless/proto/serverless"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.serverless"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	serverless.RegisterServerlessHandler(service.Server(), new(handler.Serverless))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.serverless", service.Server(), new(subscriber.Serverless))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
