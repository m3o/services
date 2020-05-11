package main

import (
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/events/api/handler"
	pb "github.com/micro/services/events/api/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.events"),
		micro.Version("latest"),
	)

	service.Init()

	// generate an elevated service account for the events api so
	// it can interact with the runtime and override the contexts
	// namespace. TODO: refactor once identity has been released.
	accName := fmt.Sprintf("%v-%v", service.Name(), service.Server().Options().Id)
	acc, err := service.Options().Auth.Generate(accName, auth.WithRoles("admin"))
	if err != nil {
		logger.Fatalf("Error generating elevated service account: %v", err)
	}
	service.Options().Auth.Init(auth.Credentials(acc.ID, acc.Secret))

	// register the handler
	pb.RegisterEventsHandler(service.Server(), handler.New(service))

	// if err := service.Run(); err != nil {
	// 	logger.Fatal(err)
	// }
}
