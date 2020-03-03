package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/users/api/handler"
	pb "github.com/micro/services/users/api/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.users"),
		micro.Version("latest"),
		micro.Auth(auth.DefaultAuth),
	)

	// Initialise service
	service.Init()

	// Register Handler
	h := handler.NewHandler(service)
	pb.RegisterUsersHandler(service.Server(), h)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
