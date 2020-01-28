package main

import (
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config/cmd"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	"github.com/kytra-app/insights-srv/api"
	"github.com/kytra-app/insights-srv/consumer"
	proto "github.com/kytra-app/insights-srv/proto"
	"github.com/kytra-app/insights-srv/storage/postgres"
	"github.com/micro/go-micro"
)

func main() {
	cmd.Init()

	// Create the service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-insights"),
		micro.Version("latest"),
	)
	service.Init()

	// Connect to the Database (Postgres)
	storageService, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer storageService.Close()

	// Connect to Message Broker (RabbitMQ)
	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	// Create the consumer
	consumer := consumer.New(storageService, service.Client())
	consumer.Subscribe()
	defer consumer.Unsubscribe()

	// Register to Service Discovery
	api := api.New(storageService, service.Client())
	proto.RegisterInsightsHandler(service.Server(), api)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
