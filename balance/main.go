package main

import (
	"github.com/m3o/services/balance/handler"
	pb "github.com/m3o/services/balance/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("balance"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterBalanceHandler(srv.Server(), new(handler.Balance))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
