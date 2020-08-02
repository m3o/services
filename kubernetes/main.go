package main

import (
	"github.com/micro/go-micro/v3/logger"

	"github.com/m3o/services/kubernetes/handler"
	pb "github.com/m3o/services/kubernetes/proto"
	"github.com/micro/micro/v3/service"
)

func main() {
	srv := service.New(
		service.Name("go.micro.service.kubernetes"),
	)

	pb.RegisterKubernetesHandler(handler.New(srv))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
