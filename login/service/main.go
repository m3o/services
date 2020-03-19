package main

import (
	"github.com/micro/services/login/service/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"

	pb "github.com/micro/services/login/service/proto/login"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.login"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	h := handler.NewHandler(service)
	pb.RegisterLoginHandler(service.Server(), h)

	// Subscribe to user update events
	micro.RegisterSubscriber("go.micro.service.users", service.Server(), h.HandleUserEvent, server.SubscriberQueue("queue.login"))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
