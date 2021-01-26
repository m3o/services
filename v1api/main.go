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
		service.Name("v1"),
		service.Version("latest"),
	)

	srv.Server().Handle(
		srv.Server().NewHandler(
			new(handler.V1Api),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1Api.Endpoint",
					Path:    []string{"^/v1/.*$"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "api",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1Api.GenerateKey",
					Path:    []string{"/v1/generatekey"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1Api.RevokeKey",
					Path:    []string{"/v1/revokekey"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1Api.UpdateAllowedPaths",
					Path:    []string{"/v1/updateallowedpaths"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1Api.ListKeys",
					Path:    []string{"/v1/listkeys"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				},
			)))
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
