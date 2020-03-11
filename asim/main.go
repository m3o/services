package main

import (
	"asim/handler"
	"asim/subscriber"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	asim "asim/proto/asim"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.asim"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	asim.RegisterAsimHandler(service.Server(), new(handler.Asim))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.asim", service.Server(), new(subscriber.Asim))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
