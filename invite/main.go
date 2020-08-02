package main

import (
	"github.com/m3o/services/invite/handler"
	pb "github.com/m3o/services/invite/proto"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
)

func main() {
	srv := service.New(
		service.Name("go.micro.service.invite"),
	)

	h := handler.NewHandler(srv)
	pb.RegisterInviteHandler(h)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
