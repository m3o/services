package main

import (
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"

	"github.com/m3o/services/secrets/service/handler"
	pb "github.com/m3o/services/secrets/service/proto"
)

func main() {
	srv := service.New(
		service.Name("go.micro.service.secrets"),
		service.Version("latest"),
	)

	pb.RegisterSecretsHandler(handler.New(srv))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
