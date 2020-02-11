package main

import (
	"fmt"

	"context"

	"github.com/micro/go-micro/v2"
	proto "github.com/micro/services/pong/proto"
)

/*
Example usage of top level service initialisation
*/

type Pong struct{}

func (g *Pong) Pong(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Pong = "Pong"
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("go.micro.pong"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init()

	// By default we'll run the server unless the flags catch us

	// Setup the server

	// Register handler
	proto.RegisterPongHandler(service.Server(), new(Pong))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
