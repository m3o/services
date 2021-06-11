package main

import (
	"github.com/m3o/services/pkg/tracing"
	"github.com/m3o/services/status/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"
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
	traceCloser := tracing.SetupOpentracing("status")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
