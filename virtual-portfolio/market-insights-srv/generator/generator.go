package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/market-insights-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/client"
)

// New returns an instance of Generator
func New(iex iex.Service, db storage.Service, client client.Client) *Generator {
	return &Generator{
		iex:    iex,
		db:     db,
		stocks: stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
	}
}

// Generator is responsible for generating the daily insights.
// A CRON job will be responsible for calling the CreateDailyInsights method.
type Generator struct {
	iex    iex.Service
	db     storage.Service
	stocks stocks.StocksService
}

// CreateDailyInsights creates an insight object in the database for each
// stock in the stocks-srv database.
func (g *Generator) CreateDailyInsights() {
	// Step 1. Fetch the stocks
	sRsp, err := g.stocks.All(context.Background(), &stocks.AllRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 2. For each stock, create yesterdays insight
	stocks := sRsp.GetStocks()
	for i, stock := range stocks {
		fmt.Printf("[Generator] Creating insight %v/%v: %v\n", i+1, len(stocks), stock.Symbol)
		go g.createTodaysStockInsight(stock)
		time.Sleep(time.Second / 25)
	}
}

func (g *Generator) createTodaysStockInsight(stock *stocks.Stock) {
	date := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)

	// Step 1. Get the previous days price change
	priceRsp, err := g.iex.PreviousDayPrice(stock.Symbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 2. Check the previous price was recent
	priceDate, err := time.Parse("2006-01-02", priceRsp.Date)
	if err != nil {
		fmt.Println(err)
		return
	}
	if priceDate.Unix() < date.Add(time.Hour*-3*24).Unix() {
		return
	}

	// Step 3. Fetch the market cap
	marketCap, err := g.iex.GetMarketCap(stock.Symbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 4. Fetch the events
	events, err := g.iex.ListUpcomingEarningsForStock(stock.Symbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 5. Detetemine if any of those events are happening today
	earningsToday := false
	for _, e := range events {
		eDate, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if date.Unix() != eDate.Unix() {
			continue
		}
		earningsToday = true
		break
	}

	// Step 6. Calculate the score
	var score float32

	// Step 6.1 Increment score if there was a big movement
	if priceRsp.ChangePercent > 5 {
		score++
	}

	// Step 6.2 Increment score if there is an event today
	if earningsToday {
		score++
	}

	// Step 6.3 Increment score if large-cap stock (>10bn)
	if marketCap > 10000000000 {
		score++
	}

	// Step 6.4 Decrement score if small-cap stock (<2bn)
	if marketCap < 2000000000 {
		score--
	}

	// Step 7. Write the insight to the database
	_, err = g.db.Create(storage.Insight{
		Date:                         date,
		AssetType:                    "Stock",
		AssetUUID:                    stock.Uuid,
		Score:                        score,
		EarningsToday:                earningsToday,
		PrevDayPriceChangePercentage: priceRsp.ChangePercent,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
