package main

import (
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2"
	"vtolstov/handler"
	"vtolstov/subscriber"

	vtolstov "vtolstov/proto/vtolstov"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.vtolstov"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	vtolstov.RegisterVtolstovHandler(service.Server(), new(handler.Vtolstov))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.vtolstov", service.Server(), new(subscriber.Vtolstov))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
