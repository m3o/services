package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/secrets/service/handler"
	pb "github.com/micro/services/secrets/service/proto"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.service.secrets"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterSecretsHandler(srv.Server(), handler.New(srv))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
