package main

import (
	"github.com/m3o/services/usage/handler"
	"github.com/robfig/cron/v3"

	nsproto "github.com/m3o/services/namespaces/proto"
	pb "github.com/micro/micro/v3/proto/auth"
	rproto "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("usage"),
		service.Version("latest"),
	)

	// Register handler
	u := handler.NewUsage(
		nsproto.NewNamespacesService("namespaces", srv.Client()),
		pb.NewAccountsService("auth", srv.Client()),
		rproto.NewRuntimeService("runtime", srv.Client()),
	)
	srv.Handle(u)

	c := cron.New()
	c.AddFunc("0 8,12,16 * * *", u.CheckUsageCron)
	c.Start()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
