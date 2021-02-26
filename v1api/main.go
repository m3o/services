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
			new(handler.V1),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1.Endpoint",
					Handler: "rpc",
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Path:    []string{"^/v1/.*$"},
					Stream:  true,
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1.GenerateKey",
					Path:    []string{"/v1/api/keys/generate"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1.RevokeKey",
					Path:    []string{"/v1/api/keys/revoke"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1.ListKeys",
					Path:    []string{"/v1/api/keys/list"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1.ListAPIs",
					Path:    []string{"/v1/api/apis/list"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "rpc",
				}),
		))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
