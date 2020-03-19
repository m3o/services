package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	handler "subscribe/handler"
	subscribe "subscribe/proto/subscribe"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.subscribe"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	subscribe.RegisterSubscribeHandler(service.Server(), handler.NewSubscribe())

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
