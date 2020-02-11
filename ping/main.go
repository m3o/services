package main

import (
	"fmt"
	"time"

	"context"

	"github.com/micro/go-micro/v2"
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
	request := g.service.Client().NewRequest("go.micro.pong", "Pong.Pong", &pongproto.Request{})
	response := &pongproto.Response{}
	time.Sleep(500 * time.Millisecond)
	if err := g.service.Client().Call(ctx, request, response); err != nil {
		return err
	}
	rsp.Ping = "Ping service called Pong and that responded: " + response.GetPong()
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("go.micro.ping"),
		micro.Version("latest"),
	)

	service.Init()

	// Register handler
	proto.RegisterPingHandler(service.Server(), &Ping{service})

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
