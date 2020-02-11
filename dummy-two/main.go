package main

import (
	"fmt"

	"context"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	proto "github.com/micro/services/dummy-two/proto"
)

/*
Example usage of top level service initialisation
*/

type DummyTwo struct{}

func (g *DummyTwo) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Dummy = "Dummy 2"
	return nil
}

// Setup and the client
func runClient(service micro.Service) {
	// Create new DummyTwo client
	DummyTwo := proto.NewDummyTwoService("DummyTwo", service.Client())

	// Call the DummyTwo
	rsp, err := DummyTwo.Hello(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	fmt.Println(rsp.Dummy)
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("go.micro.dummy-two"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),

		// Setup some flags. Specify --run_client to run the client

		// Add runtime flags
		// We could do this below too
		micro.Flags(&cli.BoolFlag{
			Name:  "run_client",
			Usage: "Launch the client",
		}),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init()

	// By default we'll run the server unless the flags catch us

	// Setup the server

	// Register handler
	proto.RegisterDummyTwoHandler(service.Server(), new(DummyTwo))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
