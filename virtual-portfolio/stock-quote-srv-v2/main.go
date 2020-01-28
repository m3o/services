package main

import (
	"fmt"
	"os"
	"time"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/helpers/worldtradingdata"
	"github.com/kytra-app/stock-quote-srv/handler"
	proto "github.com/kytra-app/stock-quote-srv/proto"
	"github.com/kytra-app/stock-quote-srv/storage/postgres"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/robfig/cron/v3"
)

func main() {
	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v2-stock-quote"),
		micro.Version("latest"),
	)
	service.Init()

	// Initialize IEX Package
	iex, err := iex.New(os.Getenv("IEX_TOKEN"))
	if err != nil {
		fmt.Printf("Could not initiate iex package: %v\n.", err)
		os.Exit(2)
	}

	// Initialize worldtradingdata Package
	wtd, err := worldtradingdata.New(os.Getenv("WORLD_TRADE_DATA_TOKEN"))
	if err != nil {
		fmt.Printf("Could not initiate world trading data package: %v\n.", err)
		os.Exit(2)
	}

	// Connect to the Database (Postgres)
	db, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		fmt.Printf("Could not connect to DB: %v\n.", err)
		os.Exit(2)
	}
	defer db.Close()

	// Register to Service Discovery
	h := handler.New(wtd, iex, db, service.Client())
	proto.RegisterStockQuoteHandler(service.Server(), h)

	// Setup CRON job to fetch data
	c := cron.New(cron.WithLocation(time.UTC))
	c.AddFunc("*/10 * * * *", h.FetchLivePrices)   // Every 10 mins
	c.AddFunc("*/30 * * * *", h.FetchIndexPrices)  // Every 30 mins
	c.AddFunc("0 23 * * *", h.FetchEndOfDayPrices) // Every night at 11pm
	c.Start()
	defer c.Stop()

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
