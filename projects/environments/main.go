package main

import (
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"

	"github.com/m3o/services/projects/environments/handler"
	pb "github.com/m3o/services/projects/environments/proto"
)

func main() {
	srv := service.New(
		service.Name("go.micro.service.projects.environments"),
		service.Version("latest"),
	)

	h := handler.NewEnvironments(srv)
	pb.RegisterEnvironmentsHandler(h)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
