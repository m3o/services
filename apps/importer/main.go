package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/apps/importer/handler"
	pb "github.com/micro/services/apps/importer/proto/importer"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.apps.importer"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register the handler
	h := handler.NewHandler(service)
	pb.RegisterImporterHandler(service.Server(), h)

	// Run the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
