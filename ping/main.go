package main

import (
	"fmt"

	"context"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	proto "github.com/micro/services/ping/proto"
	pongproto "github.com/micro/services/pong/proto"
)

/*
Example usage of top level service initialisation
*/

type Ping struct {
	service micro.Service
}

func (g *Ping) Ping(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	request := client.NewRequest("go.micro.pong", "Pong.Pong", &pongproto.Request{})
	response := &pongproto.Response{}
	if err := g.service.Client().Call(ctx, request, response); err != nil {
		return err
	}
	rsp.Ping = "Ping service called Pong and that responded: " + response.GetPong()
	return nil
}

// Setup and the client
func runClient(service micro.Service) {
	// Create new Ping client
	Ping := proto.NewPingService("Ping", service.Client())

	// Call the Ping
	rsp, err := Ping.Ping(context.TODO(), &proto.Request{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	fmt.Println(rsp.Ping)
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("go.micro.ping"),
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
	proto.RegisterPingHandler(service.Server(), &Ping{service})

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
