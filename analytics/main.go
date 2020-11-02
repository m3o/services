package main

import (
	"github.com/m3o/services/analytics/consumer"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("analytics"),
	)

	// Create the consumer
	c := &consumer.Consumer{}
	if err := c.Run(); err != nil {
		logger.Fatal(err)
	}

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
