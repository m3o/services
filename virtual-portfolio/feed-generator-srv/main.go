package main

import (
	"fmt"
	"log"

	"github.com/kytra-app/feed-generator-srv/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config/cmd"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	cmd.Init()

	// Setup the service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-feed-generator"),
		micro.Version("latest"),
	)
	service.Init()

	// Setup the broker (RabbitMQ)
	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	// Create a Handler and subscribe to events
	h := handler.New(service.Client())
	broker.Subscribe("kytra-srv-v1-posts-post-created", h.HandleNewPost, broker.Queue("feed-generator-post-created"))
	broker.Subscribe("kytra-srv-v1-followers-new-follow", h.HandleFollow, broker.Queue("feed-generator-new-follow"))
	broker.Subscribe("kytra-srv-v1-followers-new-unfollow", h.HandleUnfollow, broker.Queue("feed-generator-new-unfollow"))

	// Run the service
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
