package main

import (
	"subscribe/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	subscribe "subscribe/proto/subscribe"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.subscribe"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	subscribe.RegisterSubscribeHandler(service.Server(), new(handler.Subscribe))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
