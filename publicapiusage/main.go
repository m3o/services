package main

import (
	"github.com/m3o/services/publicapiusage/handler"
	pb "github.com/m3o/services/publicapiusage/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("publicapiusage"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPublicapiusageHandler(srv.Server(), handler.NewHandler(srv))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
