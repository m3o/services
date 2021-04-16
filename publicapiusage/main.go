package main

import (
	"github.com/m3o/services/publicapiusage/handler"
	pb "github.com/m3o/services/publicapiusage/proto"
	"github.com/robfig/cron/v3"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("publicapiusage"),
		service.Version("latest"),
	)

	p := handler.NewHandler(srv)
	// Register handler
	pb.RegisterPublicapiusageHandler(srv.Server(), p)

	c := cron.New()
	c.AddFunc("1 0 * * *", p.UsageCron)
	c.Start()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
