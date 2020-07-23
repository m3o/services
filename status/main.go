package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/status/handler"
	status "github.com/micro/services/status/proto/status"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.status"),
	)

	// Initialise service
	service.Init()
	// grab services to monitor
	svcs := service.Options().Config.Get("micro", "status", "services").StringSlice(nil)
	log.Infof("Services to monitor %+v", svcs)
	// Register Handler
	status.RegisterStatusHandler(service.Server(), handler.NewStatusHandler(svcs))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
