package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"sumo/handler"
	"sumo/subscriber"

	sumo "sumo/proto/sumo"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.sumo"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	sumo.RegisterSumoHandler(service.Server(), new(handler.Sumo))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.sumo", service.Server(), new(subscriber.Sumo))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
