package main

import (
	"github.com/m3o/services/status/handler"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("status"),
	)

	// grab services to monitor
	val, err := config.Get("micro.status.services")
	if err != nil {
		log.Warnf("Error loading config: %v", err)
	}

	services := val.StringSlice(nil)
	log.Infof("Services to monitor %+v", services)

	// Register Handler
	srv.Handle(handler.NewStatusHandler(services))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
