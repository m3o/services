package main

import (
	"github.com/m3o/services/v1api/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("v1api"),
		service.Version("latest"),
	)

	srv.Server().Handle(
		srv.Server().NewHandler(
			new(handler.V1api),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1api.Endpoint",
					Path:    []string{"^/v1api/.*$"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "api",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1api.Generate",
					Path:    []string{"/v1api/generate"},
					Method:  []string{"POST", "OPTIONS", "HEAD"},
					Handler: "rpc",
				},
			)))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
