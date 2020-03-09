package main

import (
	"github.com/micro/services/notes/service/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	pb "github.com/micro/services/notes/service/proto"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.notes"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	pb.RegisterNotesHandler(service.Server(), handler.NewHandler(service))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
