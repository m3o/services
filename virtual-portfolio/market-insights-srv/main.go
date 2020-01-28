package main

import (
	"fmt"
	"os"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/market-insights-srv/generator"
	"github.com/kytra-app/market-insights-srv/handler"
	proto "github.com/kytra-app/market-insights-srv/proto"
	"github.com/kytra-app/market-insights-srv/storage/postgres"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-market-insights"),
		micro.Version("latest"),
	)
	service.Init()

	// Initialize IEX Package
	iex, err := iex.New(os.Getenv("IEX_TOKEN"))
	if err != nil {
		fmt.Printf("Could not initiate iex package: %v\n.", err)
		os.Exit(2)
	}

	// Connect to the DB
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

	// TODO: MOVE TO CRON
	g := generator.New(iex, storageService, service.Client())
	go g.CreateDailyInsights()

	// Register to Service Discovery
	hander := handler.New(storageService)
	proto.RegisterMarketInsightsHandler(service.Server(), hander)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
