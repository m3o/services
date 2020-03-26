//go:generate ./hugo_build.sh
package main

import (
	"explore/handler"
	"explore/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	explore "explore/proto/explore"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.explore"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	explore.RegisterExploreHandler(service.Server(), new(handler.Explore))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.explore", service.Server(), new(subscriber.Explore))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
