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
	c := &consumer.Consumer{
		ErrChan: make(chan error),
	}
	if err := c.Init(); err != nil {
		logger.Fatal(err)
	}
	if err := c.Run(); err != nil {
		logger.Fatal(err)
	}

	// Run the service
	go func() {
		c.ErrChan <- srv.Run()
	}()

	// Wait for either the application to be cancelled or the consumer
	// to error
	if err := <-c.ErrChan; err != nil {
		logger.Fatal(err)
	}
}
