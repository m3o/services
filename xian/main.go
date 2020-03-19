package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"xian/handler"
	"xian/subscriber"

	xian "xian/proto/xian"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.xian"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	xian.RegisterXianHandler(service.Server(), new(handler.Xian))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.xian", service.Server(), new(subscriber.Xian))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
