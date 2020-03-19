package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"sokolovstas/handler"
	"sokolovstas/subscriber"

	sokolovstas "sokolovstas/proto/sokolovstas"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.sokolovstas"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	sokolovstas.RegisterSokolovstasHandler(service.Server(), new(handler.Sokolovstas))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.sokolovstas", service.Server(), new(subscriber.Sokolovstas))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
