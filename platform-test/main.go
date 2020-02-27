package main

import (
	"fmt"
	"platform-test/handler"
	pb "platform-test/proto"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.platform.test"),
	)
	service.Init()

	fmt.Println("Foo")

	h := handler.NewHandler(service)
	pb.RegisterPlatformTestHandler(service.Server(), h)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
