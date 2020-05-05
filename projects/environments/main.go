package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/projects/environments/handler"
	pb "github.com/micro/services/projects/environments/proto"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.service.projects.environments"),
		micro.Version("latest"),
	)

	srv.Init()

	h := handler.NewEnvironments(srv)
	pb.RegisterEnvironmentsHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
