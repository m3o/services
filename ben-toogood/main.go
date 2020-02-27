package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"ben-toogood/handler"
	"ben-toogood/subscriber"

	ben "ben-toogood/proto/ben"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.ben"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	ben.RegisterBenHandler(service.Server(), new(handler.Ben))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.ben", service.Server(), new(subscriber.Ben))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
