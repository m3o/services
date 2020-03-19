package main

import (
	"log"

	"github.com/micro/go-micro/v2"
	"github.com/micro/services/location/handler"
	"github.com/micro/services/location/ingester"
	proto "github.com/micro/services/location/proto/location"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.location"),
		micro.Version("latest"),
	)

	service.Init()

	proto.RegisterLocationHandler(service.Server(), new(handler.Location))

	micro.RegisterSubscriber(ingester.Topic, service.Server(), new(ingester.Geo))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
