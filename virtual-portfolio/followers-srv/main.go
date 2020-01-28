package main

import (
	"fmt"
	"os"

	"github.com/kytra-app/followers-srv/handler"
	proto "github.com/kytra-app/followers-srv/proto"
	"github.com/kytra-app/followers-srv/storage/postgres"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-micro/broker"
	"github.com/pkg/errors"
)

func main() {
	storageService, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to the database"))
	}
	defer storageService.Close()

	// Connect to Message Broker (RabbitMQ)
	b := rabbitmq.NewBroker(
		broker.Addrs(os.Getenv("MICRO_BROKER_ADDRESS")),
	)
	if err := b.Connect(); err != nil {
		panic(errors.Wrap(err, "Could not connect to the message broker"))
	}

	service := micro.NewService(
		micro.Name("kytra-srv-v1-followers"),
		micro.Version("latest"),
	)
	service.Init()

	proto.RegisterFollowersHandler(service.Server(), handler.New(storageService, b))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
