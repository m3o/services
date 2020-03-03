package main

import (
	"api/handler"
	client "api/proto/client"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/api"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.client"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	client.RegisterClientHandler(service.Server(), &handler.Client{service.Client()}, api.WithEndpoint(
		// TODO: remove when api supports Call method as default for /foo singular paths
		&api.Endpoint{
			Name:    "Client.Call",
			Path:    []string{"^/client?$"},
			Method:  []string{"GET", "POST"},
			Handler: "rpc",
		},
	))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
