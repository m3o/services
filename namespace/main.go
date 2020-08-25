package main

import (
	"namespace/handler"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("namespace"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(new(handler.Namespace))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
