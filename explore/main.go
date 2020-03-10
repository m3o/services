package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"explore/handler"
	"explore/subscriber"

	explore "explore/proto/explore"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.explore"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	explore.RegisterExploreHandler(service.Server(), new(handler.Explore))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.explore", service.Server(), new(subscriber.Explore))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
