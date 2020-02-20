package main

import (
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2"
	"peartree/handler"
	"peartree/subscriber"

	peartree "peartree/proto/peartree"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.peartree"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	peartree.RegisterPeartreeHandler(service.Server(), new(handler.Peartree))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.peartree", service.Server(), new(subscriber.Peartree))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
