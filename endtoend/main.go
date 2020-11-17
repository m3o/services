package main

import (
	"github.com/m3o/services/endtoend/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("endtoend"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(handler.NewEndToEnd(srv))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
