package main

import (
	"github.com/m3o/services/alert/handler"
	"github.com/m3o/services/pkg/tracing"
	"github.com/micro/micro/v3/service"
	log "github.com/micro/micro/v3/service/logger"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("alert"),
	)

	// Register Handler
	srv.Handle(handler.NewAlert())

	traceCloser := tracing.SetupOpentracing("alert")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
