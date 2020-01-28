package main

import (
	"fmt"
	"os"

	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/rabbitmq"

	"github.com/kytra-app/feeditems-srv/handler"
	proto "github.com/kytra-app/feeditems-srv/proto"
	"github.com/kytra-app/feeditems-srv/storage/postgres"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	// Connect to the Database (Postgres)
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
	if err := broker.Init(); err != nil {
		panic(errors.Wrap(err, "Could not connect to the message broker"))
	}
	if err := broker.Connect(); err != nil {
		panic(errors.Wrap(err, "Could not connect to the message broker"))
	}

	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-feed-items"),
		micro.Version("latest"),
	)
	service.Init()

	// Register to Service Discovery (Consul)
	hander := handler.New(storageService)
	proto.RegisterFeedItemsHandler(service.Server(), hander)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
