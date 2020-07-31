package main

import (
	"fmt"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	mauth "github.com/micro/micro/v3/service/auth"

	"github.com/m3o/services/events/api/handler"
	pb "github.com/m3o/services/events/api/proto"
)

func main() {
	srv := service.New(
		service.Name("go.micro.api.events"),
		service.Version("latest"),
	)

	// generate an elevated service account for the events api so
	// it can interact with the runtime and override the contexts
	// namespace. TODO: refactor once identity has been released.
	accName := fmt.Sprintf("%v-%v", srv.Name(), srv.Server().Options().Id)
	acc, err := mauth.Generate(accName, auth.WithScopes("admin", "service"))
	if err != nil {
		logger.Fatalf("Error generating elevated service account: %v", err)
	}
	mauth.DefaultAuth.Init(auth.Credentials(acc.ID, acc.Secret))

	// register the handler
	pb.RegisterEventsHandler(handler.New(srv))

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
