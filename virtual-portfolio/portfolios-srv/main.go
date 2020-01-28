package main

import (
	"fmt"
	"os"

	"github.com/kytra-app/portfolios-srv/handler"
	proto "github.com/kytra-app/portfolios-srv/proto"
	"github.com/kytra-app/portfolios-srv/storage/postgres"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/micro/go-micro"
	"github.com/pkg/errors"
)

func main() {
	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-portfolios"),
		micro.Version("latest"),
	)
	service.Init()

	// Connect to the Database (Postgres)
	db, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to the database"))
	}
	defer db.Close()

	// Register to Service Discovery
	hander := handler.New(db)
	proto.RegisterPortfoliosHandler(service.Server(), hander)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
