package main

import (
	"fmt"
	"os"
	"time"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/stock-movers-srv/handler"
	proto "github.com/kytra-app/stock-movers-srv/proto"
	"github.com/kytra-app/stock-movers-srv/storage/postgres"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-stock-movers"),
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
	broker := rabbitmq.NewBroker(
		broker.Addrs(os.Getenv("MICRO_BROKER_ADDRESS")),
	)
	if err := broker.Connect(); err != nil {
		panic(errors.Wrap(err, "Could not connect to the message broker"))
	}

	// Register to Service Discovery
	hander := handler.New(iex, db, service.Client(), broker)

	// Setup ticker to fetch movements from IEX
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		for {
			fmt.Println("Fetching Movements")
			hander.FetchMovements()
			<-ticker.C
		}
	}()

	proto.RegisterStockMoversHandler(service.Server(), hander)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
