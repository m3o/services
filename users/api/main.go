package main

import (
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"

	"github.com/m3o/services/users/api/handler"
	pb "github.com/m3o/services/users/api/proto"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.api.users"),
		service.Version("latest"),
	)

	// Register Handler
	h := handler.NewHandler(srv)
	pb.RegisterUsersHandler(h)

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
