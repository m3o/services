package main

import (
	"github.com/m3o/services/notifications/service/handler"
	"github.com/m3o/services/notifications/service/subscriber"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"

	notifications "github.com/m3o/services/notifications/service/proto/notifications"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.notifications"),
		service.Version("latest"),
	)

	// Register Handler
	notifications.RegisterNotificationsHandler(new(handler.Notifications))

	// Register Struct as Subscriber
	service.RegisterSubscriber("go.micro.service.notifications", new(subscriber.Subscriber))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
