package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/platform/service/handler"
	pb "github.com/micro/services/platform/service/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.service.platform"),
	)

	service.Init()

	h := handler.NewHandler(service)
	pb.RegisterPlatformHandler(service.Server(), h)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
