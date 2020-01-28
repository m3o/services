package main

import (
	"log"
	"os"
	"time"

	"github.com/kytra-app/helpers/mailer"
	"github.com/kytra-app/notifications-mailer-srv/handler"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/robfig/cron/v3"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra-srv-v1-notifications-mailer"),
		micro.Version("latest"),
	)
	service.Init()

	mailer := mailer.New(os.Getenv("MAILER_USERNAME"), os.Getenv("MAILER_PASSWORD"))

	h, err := handler.New(service.Client(), mailer)
	if err != nil {
		log.Fatal(err)
	}

	c := cron.New(cron.WithLocation(time.UTC))
	c.AddFunc("0 7 * * *", h.SendDailyEmails)
	c.Start()
	defer c.Stop()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
