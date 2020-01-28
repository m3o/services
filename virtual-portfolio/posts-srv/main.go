package main

import (
	"fmt"
	"os"

	"github.com/kytra-app/posts-srv/handler"
	proto "github.com/kytra-app/posts-srv/proto"
	"github.com/kytra-app/posts-srv/storage/postgres"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-posts"),
		micro.Version("latest"),
	)
	service.Init()

	// Connect to the Database (Postgres)
	db, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to the database"))
	}
	defer db.Close()

	// Connect to Message Broker (RabbitMQ)
	b := rabbitmq.NewBroker(
		broker.Addrs(os.Getenv("MICRO_BROKER_ADDRESS")),
	)
	if err := b.Connect(); err != nil {
		panic(errors.Wrap(err, "Could not connect to the message broker"))
	}

	// Register to Service Discovery (Consul)
	hander := handler.New(db, b)
	proto.RegisterPostsHandler(service.Server(), hander)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
