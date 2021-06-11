package main

import (
	"github.com/m3o/services/customers/handler"
	"github.com/m3o/services/pkg/tracing"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("customers"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(handler.New(srv))
	traceCloser := tracing.SetupOpentracing("customers")
	defer traceCloser.Close()
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
