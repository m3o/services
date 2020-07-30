package main

import (
	log "github.com/micro/go-micro/v3/logger"

	"github.com/m3o/services/status/handler"
	status "github.com/m3o/services/status/proto/status"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.status"),
	)

	// grab services to monitor
	svcs := config.Get("micro", "status", "services").StringSlice(nil)
	log.Infof("Services to monitor %+v", svcs)
	// Register Handler
	status.RegisterStatusHandler(handler.NewStatusHandler(svcs))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
