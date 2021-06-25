package main

import (
	"github.com/m3o/services/onboarding/handler"
	"github.com/m3o/services/pkg/tracing"
	"github.com/micro/micro/v3/service"
	mauth "github.com/micro/micro/v3/service/auth/client"
	log "github.com/micro/micro/v3/service/logger"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("onboarding1"),
	)

	// passing in auth because the DefaultAuth is the one used to set up the service
	auth := mauth.NewAuth()

	// Register Handler
	srv.Handle(handler.NewSignup(srv, auth))

	// kick off event consumption
	handler.NewOnboarding(srv)
	traceCloser := tracing.SetupOpentracing("onboarding")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
