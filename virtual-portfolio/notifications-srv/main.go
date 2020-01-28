package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kytra-app/notifications-srv/api"
	"github.com/kytra-app/notifications-srv/consumer"
	proto "github.com/kytra-app/notifications-srv/proto"
	"github.com/kytra-app/notifications-srv/storage/postgres"
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
		micro.Name("kytra-srv-v1-notifications"),
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
		panic(errors.New("Could not connect to the database"))
	}
	defer db.Close()

	// Setup the broker (RabbitMQ)
	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	// Create a Handler and subscribe to events
	c := consumer.New(service.Client(), db)
	broker.Subscribe("kytra-srv-v1-posts-post-created", c.ConsumeNewPost, broker.Queue("notifications-srv-post-created"))
	broker.Subscribe("kytra-srv-v1-followers-new-follow", c.ConsumeNewFollow, broker.Queue("notifications-srv-new-follow"))
	broker.Subscribe("kytra-srv-v1-comments-comment-created", c.ConsumeNewComment, broker.Queue("notifications-srv-comment-created"))

	// Run the service
	proto.RegisterNotificationsHandler(service.Server(), api.New(service.Client(), db))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
