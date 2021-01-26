package main

import (
	"github.com/m3o/services/quota/handler"
	pb "github.com/m3o/services/quota/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("quota"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterQuotaHandler(srv.Server(), handler.New(srv.Client()))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
