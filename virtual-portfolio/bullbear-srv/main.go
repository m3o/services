package main

import (
	"fmt"
	"os"

	"github.com/kytra-app/bullbear-srv/handler"
	proto "github.com/kytra-app/bullbear-srv/proto"
	"github.com/kytra-app/bullbear-srv/storage/postgres"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/micro/go-micro"
	"github.com/pkg/errors"
)

func main() {
	storageService, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to the database"))
	}
	defer storageService.Close()

	service := micro.NewService(
		micro.Name("kytra-srv-v1-bullbear"),
		micro.Version("latest"),
	)
	service.Init()

	proto.RegisterBullBearHandler(service.Server(), handler.New(storageService))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
