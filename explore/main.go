package main

import (
	"github.com/m3o/services/explore/handler"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/registry"
)

func main() {
	srv := service.New(
		service.Name("explore"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(handler.NewExploreHandler(registry.DefaultRegistry))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
