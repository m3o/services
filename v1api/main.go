package main

import (
	"github.com/m3o/services/v1api/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/metrics"
	"github.com/micro/micro/v3/service/metrics/logging"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("v1"),
		service.Version("latest"),
	)

	srv.Server().Handle(
		srv.Server().NewHandler(
			new(handler.V1api),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1api.Endpoint",
					Path:    []string{"^/v1/.*$"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "api",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1api.Generate",
					Path:    []string{"/v1/generate"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1api.Revoke",
					Path:    []string{"/v1/revoke"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1api.ListKeys",
					Path:    []string{"/v1/listkeys"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				},
			)))
	metrics.DefaultMetricsReporter = logging.New()
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
