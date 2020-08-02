package main

import (
	"github.com/m3o/services/alert/handler"
	log "github.com/micro/go-micro/v3/logger"

	alert "github.com/m3o/services/alert/proto/alert"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/store"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("alert"),
	)

	// Register Handler
	alert.RegisterAlertHandler(handler.NewAlert(store.DefaultStore))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
