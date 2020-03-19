package main

import (
	"rex/handler"
	"rex/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	rex "rex/proto/rex"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.rex"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	rex.RegisterRexHandler(service.Server(), new(handler.Rex))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.rex", service.Server(), new(subscriber.Rex))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
