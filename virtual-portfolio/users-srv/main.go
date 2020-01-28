package main

import (
	"fmt"
	"os"

	"github.com/kytra-app/users-srv/handler"
	proto "github.com/kytra-app/users-srv/proto"
	"github.com/kytra-app/users-srv/storage/postgres"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
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
		micro.Name("kytra-srv-v1-users"),
		micro.Version("latest"),
	)
	service.Init()

	proto.RegisterUsersHandler(service.Server(), handler.New(storageService))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
