package main

import (
	"github.com/m3o/services/alert/handler"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	alert "github.com/m3o/services/alert/proto/alert"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("alert"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	alert.RegisterAlertHandler(service.Server(), handler.NewAlert(service.Options().Store))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
