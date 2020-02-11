package main

import (
	"fmt"

	"context"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	proto "github.com/micro/services/dummy-one/proto"
	twoproto "github.com/micro/services/dummy-two/proto"
)

/*
Example usage of top level service initialisation
*/

type DummyOne struct {
	service micro.Service
}

func (g *DummyOne) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	request := client.NewRequest("go.micro.dummy-two", "Hello", &twoproto.Request{})
	response := &twoproto.Response{}
	if err := g.service.Client().Call(ctx, request, response); err != nil {
		return err
	}
	rsp.Dummy = "Dummy 1 -> " + response.GetDummy()
	return nil
}

// Setup and the client
func runClient(service micro.Service) {
	// Create new DummyOne client
	DummyOne := proto.NewDummyOneService("DummyOne", service.Client())

	// Call the DummyOne
	rsp, err := DummyOne.Hello(context.TODO(), &proto.Request{})
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
		micro.Name("go.micro.dummy-one"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),

		micro.Flags(&cli.BoolFlag{
			Name:  "run_client",
			Usage: "Launch the client",
		}),
	)

	service.Init()

	// Register handler
	proto.RegisterDummyOneHandler(service.Server(), new(DummyOne))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
