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
	svcs := config.Get("micro", "status", "services").StringSlice(nil)
	log.Infof("Services to monitor %+v", svcs)

	// Register Handler
	srv.Handle(handler.NewStatusHandler(svcs))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
