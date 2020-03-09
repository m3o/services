package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/notes/api/handler"
	pb "github.com/micro/services/notes/api/proto/notes"
)

func main() {
	// New Service
	srv := micro.NewService(
		micro.Name("go.micro.api.notes"),
		micro.Version("latest"),
	)

	// Initialise service
	srv.Init()

	// Register Handler
	pb.RegisterNotesHandler(srv.Server(), handler.NewHandler(srv))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
