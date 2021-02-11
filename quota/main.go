package main

import (
	"github.com/m3o/services/quota/handler"
	pb "github.com/m3o/services/quota/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	"github.com/robfig/cron/v3"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("quota"),
		service.Version("latest"),
	)

	q := handler.New(srv.Client())
	c := cron.New()
	c.AddFunc("0 0 * * *", q.ResetQuotas)
	c.Start()

	// Register handler
	pb.RegisterQuotaHandler(srv.Server(), q)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
