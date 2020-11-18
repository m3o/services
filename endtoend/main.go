package main

import (
	"github.com/m3o/services/endtoend/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/robfig/cron/v3"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("endtoend"),
		service.Version("latest"),
	)

	// Register handler
	e := handler.NewEndToEnd(srv)
	srv.Handle(e)

	c := cron.New()
	c.AddFunc("0/5 * * * *", e.CronCheck)
	c.Start()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
