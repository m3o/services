package main

import (
	"github.com/m3o/services/login/service/handler"

	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/micro/v3/service"

	pb "github.com/m3o/services/login/service/proto/login"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.login"),
		service.Version("latest"),
	)

	// Register Handler
	h := handler.NewHandler(srv)
	pb.RegisterLoginHandler(h)

	// Subscribe to user update events
	service.RegisterSubscriber("go.micro.service.users", h.HandleUserEvent, server.SubscriberQueue("queue.login"))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
