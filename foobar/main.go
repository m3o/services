package main

import (
	"fmt"
	"foobar/handler"
	"foobar/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	foobar "foobar/proto/foobar"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.foobar"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	foobar.RegisterFoobarHandler(service.Server(), new(handler.Foobar))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.foobar", service.Server(), new(subscriber.Foobar))

	fmt.Println("Foo")

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
