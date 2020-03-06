package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"subscribe/handler"
	"subscribe/subscriber"

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

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.subscribe", service.Server(), new(subscriber.Subscribe))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
