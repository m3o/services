package main

import (
	"fmt"
	"os"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/trades-srv/handler"
	proto "github.com/kytra-app/trades-srv/proto"
	"github.com/kytra-app/trades-srv/storage/postgres"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	rabbitmq "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-trades"),
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
		fmt.Printf("Could not connect to message broker: %v\n.", err)
		os.Exit(2)
	}

	// Initialize IEX Package
	iex, err := iex.New(os.Getenv("IEX_TOKEN"))
	if err != nil {
		fmt.Printf("Could not initiate iex package: %v\n.", err)
		os.Exit(2)
	}

	// Register to Service Discovery
	hander := handler.New(iex, db, b, service.Client())
	proto.RegisterTradesHandler(service.Server(), hander)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
