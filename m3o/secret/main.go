package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/m3o/secret/handler"
	pb "github.com/micro/services/m3o/secret/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.m3o.secret"),
		micro.Version("latest"),
	)

	service.Init()

	pb.RegisterSecretServiceHandler(service.Server(), handler.NewSecret(service))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
