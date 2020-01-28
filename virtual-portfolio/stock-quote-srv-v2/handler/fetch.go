package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kytra-app/helpers/unique"
	insights "github.com/kytra-app/insights-srv/proto"
	storage "github.com/kytra-app/stock-quote-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
)

// FetchLivePrices retrieves the intraday prices from IEX for assets
// which are being actively trades on the platform
func (h *Handler) FetchLivePrices() {
	fmt.Println("Fetching live prices")

	// Step 1. Fetch the assets being traded or with insights
	tRsp, err := h.trades.AllAssets(context.Background(), &trades.AllAssetsRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	iRsp, err := h.insights.ListAssets(context.Background(), &insights.ListAssetsRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 2. Fetch the symbols for these assets
	uuids := []string{}
	for _, a := range tRsp.GetAssets() {
		uuids = append(uuids, a.Uuid)
	}
	for _, a := range iRsp.GetAssets() {
		uuids = append(uuids, a.Uuid)
	}
	uuids = unique.Strings(uuids)

	sRsp, err := h.stocks.List(context.Background(), &stocks.ListRequest{Uuids: uuids})
	if err != nil {
		fmt.Println(err)
		return
	}

	symbolForUUID := make(map[string]string, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		symbolForUUID[s.Uuid] = s.Symbol
	}
	// Step 3. Fetch the prices
	for _, uuid := range uuids {
		symbol := symbolForUUID[uuid]
		quote, err := h.iex.Quote(symbol)

		if err != nil {
			fmt.Println(err)
			return
		}

		var q storage.Quote
		if quote.USMarketOpen {
			q, err = h.db.Create(storage.Quote{
				Price:         int32(quote.LatestPrice * 100),
				StockUUID:     uuid,
				CreatedAt:     time.Unix(int64(quote.LatestUpdate/1000), 0),
				ChangePercent: quote.ChangePercent * 100,
				MarketOpen:    true,
			})
		} else {
			// Not all stocks have extended prices
			if quote.ExtendedPrice == 0 {
				continue
			}

			q, err = h.db.Create(storage.Quote{
				Price:         int32(quote.ExtendedPrice * 100),
				StockUUID:     uuid,
				CreatedAt:     time.Unix(int64(quote.ExtendedPriceTime/1000), 0),
				ChangePercent: quote.ExtendedChangePercent * 100,
				MarketOpen:    false,
			})
		}

		if err != nil {
			fmt.Println(err)
			fmt.Println(quote)
			fmt.Println(q)
			return
		}
	}
}

// FetchEndOfDayPrices retrieves the end of day prices from IEX for
// all assets stored in the platform
func (h *Handler) FetchEndOfDayPrices() {
	fmt.Println("Fetching end of day prices")

	// Step 1. Fetch all stocks
	sRsp, err := h.stocks.All(context.Background(), &stocks.AllRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 2. Fetch the prices
	for _, stock := range sRsp.GetStocks() {
		go func(stock *stocks.Stock) {
			quote, err := h.iex.Quote(stock.Symbol)

			if err != nil {
				fmt.Println(err)
				return
			}

			q, err := h.db.Create(storage.Quote{
				StockUUID:     stock.Uuid,
				Price:         int32(quote.LatestPrice * 100),
				ChangePercent: quote.ChangePercent * 100,
				MarketOpen:    true,
			})

			if err != nil {
				fmt.Println(err)
				fmt.Println(quote)
				fmt.Println(q)
				return
			}
		}(stock)

		time.Sleep(time.Second / 15)
	}
}

// FetchHistoricPrices gets the end of day price for all stocks for the past 5yrs
func (h *Handler) FetchHistoricPrices() {
	sRsp, err := h.stocks.All(context.Background(), &stocks.AllRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, s := range sRsp.GetStocks() {
		go func(s *stocks.Stock) {
			prices, err := h.iex.HistoricalPrices(s.Symbol, "5y", false)
			if err != nil {
				fmt.Println(err)
			}

			for _, p := range prices {
				t, _ := time.Parse("2006-01-02", p.Date)
				h.db.Create(storage.Quote{
					StockUUID: s.Uuid,
					Price:     int32(p.Close * 100),
					CreatedAt: t.Add(time.Hour * 23),
				})
			}

		}(s)

		time.Sleep(time.Second / 2)
	}
}

// FetchIndexPrices gets the index prices (e.g. the Nasdaq Composite) from WorldTradeData
func (h *Handler) FetchIndexPrices() {
	symbol := "^IXIC"

	quote, err := h.wtd.Quote(symbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	price, err := strconv.ParseFloat(quote.Price, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	pct, err := strconv.ParseFloat(quote.ChangePercent, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	createdAt, err := quote.LastUpdated()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = h.db.Create(storage.Quote{
		StockUUID:     symbol,
		Price:         int32(price * 100),
		ChangePercent: float32(pct),
		MarketOpen:    true,
		CreatedAt:     createdAt,
	})
}
