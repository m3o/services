package main

import (
	"fmt"
	"os"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/stock-target-price-srv/handler"
	proto "github.com/kytra-app/stock-target-price-srv/proto"
	"github.com/kytra-app/stock-target-price-srv/storage/postgres"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	cmd.Init()

	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-stock-target-price"),
		micro.Version("latest"),
	)
	service.Init()

	// Initialize IEX Package
	iex, err := iex.New(os.Getenv("IEX_TOKEN"))
	if err != nil {
		fmt.Printf("Could not initiate iex package: %v\n.", err)
		os.Exit(2)
	}

	// Connect to the Database (Postgres)
	db, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		fmt.Printf("Could not connect to DB: %v\n.", err)
		os.Exit(2)
	}
	defer db.Close()

	// Connect to Message Broker (RabbitMQ)
	brkr := rabbitmq.NewBroker(
		broker.Addrs(os.Getenv("MICRO_BROKER_ADDRESS")),
	)
	if err := brkr.Connect(); err != nil {
		panic(errors.Wrap(err, "Could not connect to the message broker"))
	}

	// Register to Service Discovery
	hander := handler.New(iex, db, service.Client(), brkr)
	proto.RegisterStockTargetPriceHandler(service.Server(), hander)

	// Consume messages
	sub, err := brkr.Subscribe(
		"kytra-srv-v1-insights-insight-created",
		hander.HandleNewInsight,
		broker.Queue("stock-target-price-insight-created"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not subscrube to the queue"))
	}
	defer sub.Unsubscribe()

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
