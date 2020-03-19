package main

import (
	"log"

	"github.com/micro/go-micro/v2"
	proto "github.com/micro/services/helloworld/api/proto/hello"
	hello "github.com/micro/services/helloworld/proto/helloworld"

	"context"
)

type Helloworld struct {
	Client hello.HelloworldService
}

func (g *Helloworld) Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Print("Received Helloworld.Call API request")

	// make the request
	response, err := g.Client.Call(ctx, &hello.Request{Name: req.Name})
	if err != nil {
		return err
	}

	// set api response
	rsp.Msg = response.Msg
	return nil
}

func main() {
	// Create service
	service := micro.NewService(
		micro.Name("go.micro.api.helloworld"),
	)

	// Init to parse flags
	service.Init()

	// Register Handlers
	proto.RegisterHelloworldHandler(service.Server(), &Helloworld{
		// Create Service Client
		Client: hello.NewHelloworldService("go.micro.service.helloworld", service.Client()),
	})

	// for handler use

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
