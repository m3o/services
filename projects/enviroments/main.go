package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/projects/enviroments/handler"
	pb "github.com/micro/services/projects/enviroments/proto"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.service.projects.enviroments"),
		micro.Version("latest"),
	)

	srv.Init()

	h := handler.NewEnviroments(srv)
	pb.RegisterEnviromentsHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
