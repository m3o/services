package main

import (
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"

	"github.com/m3o/services/users/service/handler"
	pb "github.com/m3o/services/users/service/proto"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.users"),
		service.Version("latest"),
	)

	// Register Handler
	h, err := handler.NewHandler(srv)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterUsersHandler(h)

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
