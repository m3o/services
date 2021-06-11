package main

import (
	"github.com/m3o/services/pkg/tracing"
	"github.com/m3o/services/platform/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	srv := service.New(
		service.Name("platform"),
	)

	srv.Handle(handler.New(srv))
	traceCloser := tracing.SetupOpentracing("platform")
	defer traceCloser.Close()

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
