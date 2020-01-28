package worldtradingdata

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Quote returns the latest price for a stock
func (h Handler) Quote(symbol string) (*Quote, error) {
	rsp, err := h.Get(fmt.Sprintf("stock?symbol=%v", symbol))
	if err != nil {
		return nil, err
	}

	var data StockResponse
	if err := json.Unmarshal(rsp, &data); err != nil {
		return nil, err
	}

	if len(data.Quotes) == 0 {
		return nil, nil
	}

	return data.Quotes[0], nil
}

// Quote is a live price for a stock
type Quote struct {
	Symbol             string `json:"symbol"`               // Example: "^IXIC",
	Name               string `json:"name"`                 // Example: "Nasdaq Composite",
	Currency           string `json:"currency"`             // Example: "N/A",
	Price              string `json:"price"`                // Example: "8639.77",
	PriceOpen          string `json:"price_open"`           // Example: "8623.56",
	DayHigh            string `json:"day_high"`             // Example: "8650.76",
	DayLow             string `json:"day_low"`              // Example: "8600.82",
	OneYearHigh        string `json:"52_week_high"`         // Example: "8705.91",
	OneYearLow         string `json:"52_week_low"`          // Example: "6190.17",
	DayChange          string `json:"day_change"`           // Example: "17.94",
	ChangePercent      string `json:"change_pct"`           // Example: "0.21",
	CloseYesterday     string `json:"close_yesterday"`      // Example: "8621.83",
	MarketCap          string `json:"market_cap"`           // Example: "N/A",
	Volume             string `json:"volume"`               // Example: "688352755",
	VolumeAverage      string `json:"volume_avg"`           // Example: "N/A",
	Shares             string `json:"shares"`               // Example: "N/A",
	StockExchangeLong  string `json:"stock_exchange_long"`  // Example: "",
	StockExchangeShort string `json:"stock_exchange_short"` // Example: "INDEXNASDAQ",
	Timezone           string `json:"timezone"`             // Example: "EST",
	TimezoneName       string `json:"timezone_name"`        // Example: "America/New_York",
	GMTOffset          string `json:"gmt_offset"`           // Example: "-18000",
	LastTradeTime      string `json:"last_trade_time"`      // Example: "2019-12-10 11:28:33",
	PE                 string `json:"pe"`                   // Example: "N/A",
	EPS                string `json:"eps"`                  // Example: "N/A"
}

// LastUpdated is the timezone adjusted last updated time
func (q Quote) LastUpdated() (time.Time, error) {
	format := "2006-01-02 15:04:05 MST"
	value := fmt.Sprintf("%v %v", q.LastTradeTime, q.Timezone)

	t, err := time.Parse(format, value)
	if err != nil {
		return t, err
	}

	seconds, err := strconv.Atoi(q.GMTOffset)
	if err != nil {
		return t, err
	}
	return t.Add(time.Second * -time.Duration(seconds)), nil
}

// StockResponse is the datatype returned by the Stock Endpoint
type StockResponse struct {
	SymbolsRequested int64    `json:"symbols_requested"`
	SymbolsReturned  int64    `json:"symbols_returned"`
	Quotes           []*Quote `json:"data"`
}
