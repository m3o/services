package iex

import (
	"encoding/json"
	"fmt"
)

// KeyStats calls the Key Stats Endpoint
func (h Handler) KeyStats(symbol string) (KeyStatsResponse, error) {
	rsp, err := h.Get(fmt.Sprintf("stock/%v/stats", symbol))

	if err != nil {
		return KeyStatsResponse{}, err
	}

	var data KeyStatsResponse
	if err := json.Unmarshal(rsp, &data); err != nil {
		return data, err
	}

	return data, nil
}

// GetMarketCap returns the market cap for a given stock
func (h Handler) GetMarketCap(symbol string) (float32, error) {
	rsp, err := h.Get(fmt.Sprintf("stock/%v/stats/marketcap", symbol))
	if err != nil {
		return 0, err
	}

	var result float32
	if err := json.Unmarshal(rsp, &result); err != nil {
		return 0, err
	}

	return result, nil
}

// GetPriceTarget returns the latest price targets for a stock
func (h Handler) GetPriceTarget(symbol string) (PriceTargetResponse, error) {
	rsp, err := h.Get(fmt.Sprintf("stock/%v/price-target", symbol))
	if err != nil {
		return PriceTargetResponse{}, err
	}

	var data PriceTargetResponse
	if err := json.Unmarshal(rsp, &data); err != nil {
		return data, err
	}

	return data, nil
}

// Quote returns the latest price for a stock
func (h Handler) Quote(symbol string) (QuoteResponse, error) {
	rsp, err := h.Get(fmt.Sprintf("stock/%v/quote", symbol))
	if err != nil {
		return QuoteResponse{}, err
	}

	var data QuoteResponse
	if err := json.Unmarshal(rsp, &data); err != nil {
		return data, err
	}

	return data, nil
}

// PreviousDayPrice returns previous day adjusted price data for a stock
func (h Handler) PreviousDayPrice(symbol string) (PriceResponse, error) {
	rsp, err := h.Get(fmt.Sprintf("stock/%v/previous", symbol))
	if err != nil {
		return PriceResponse{}, err
	}

	var data PriceResponse
	if err := json.Unmarshal(rsp, &data); err != nil {
		return data, err
	}

	return data, nil
}

// ListPreviousDayPrices returns previous day adjusted price data for all stocks
func (h Handler) ListPreviousDayPrices() ([]PriceResponse, error) {
	url := fmt.Sprintf("stock/market/previous")
	fmt.Println(url)
	rsp, err := h.Get(url)
	if err != nil {
		return []PriceResponse{}, err
	}

	var data []PriceResponse
	err = json.Unmarshal(rsp, &data)
	return data, err
}

// ListUpcomingEarnings returns previous day adjusted price data for all stocks
func (h Handler) ListUpcomingEarnings() ([]EventResponse, error) {
	url := fmt.Sprintf("stock/market/upcoming-earnings")
	fmt.Println(url)
	rsp, err := h.Get(url)
	if err != nil {
		return []EventResponse{}, err
	}

	var data []EventResponse
	err = json.Unmarshal(rsp, &data)
	return data, err
}

// ListUpcomingEarningsForStock returns previous day adjusted price data for the given stock
func (h Handler) ListUpcomingEarningsForStock(symbol string) ([]EventResponse, error) {
	url := fmt.Sprintf("stock/%v/upcoming-earnings", symbol)
	fmt.Println(url)
	rsp, err := h.Get(url)
	if err != nil {
		return []EventResponse{}, err
	}

	var data []EventResponse
	err = json.Unmarshal(rsp, &data)
	return data, err
}

// HistoricalPrices returns adjusted and unadjusted historical data for up to 15 years. Useful for building charts.
func (h Handler) HistoricalPrices(symbol, rng string, closeOnly bool) ([]PriceResponse, error) {
	url := fmt.Sprintf("stock/%v/chart/%v?chartCloseOnly=%v", symbol, rng, closeOnly)
	fmt.Println(url)
	rsp, err := h.Get(url)
	if err != nil {
		return []PriceResponse{}, err
	}

	var data []PriceResponse
	if err := json.Unmarshal(rsp, &data); err != nil {
		return data, err
	}

	// The IEX API returns the 5dm dataset in reverse order
	if rng == "5dm" {
		reverse := make([]PriceResponse, len(data))
		for i := range data {
			reverse[i] = data[len(data)-i-1]
		}
		data = reverse
	}

	return data, nil
}

// EventResponse is the datatype returned by the Events Endpoints
type EventResponse struct {
	Symbol string `json:"symbol"`     // Example: "RESN",
	Date   string `json:"reportDate"` // Example: "2020-02-29"
}

// KeyStatsResponse is the datatype returned by the Key Stats Endpoint
type KeyStatsResponse struct {
	CompanyName         string  `json:"companyName"`         // Example: "Apple Inc.",
	MarketCap           float32 `json:"marketcap"`           // Example: 760334287200,
	Week52High          float32 `json:"week52high"`          // Example: 156.65,
	Week52Low           float32 `json:"week52low"`           // Example: 93.63,
	Week52Change        float32 `json:"week52change"`        // Example: 58.801903,
	SharesOutstanding   float32 `json:"sharesOutstanding"`   // Example: 5213840000,
	Float               float32 `json:"float"`               // Example: 5203997571,
	Avg10Volume         float32 `json:"avg10Volume"`         // Example: 2774000,
	Avg30Volume         float32 `json:"avg30Volume"`         // Example: 12774000,
	Day200MovingAvg     float32 `json:"day200MovingAvg"`     // Example: 140.60541,
	Day50MovingAvg      float32 `json:"day50MovingAvg"`      // Example: 156.49678,
	Employees           float32 `json:"employees"`           // Example: 120000,
	TTMEPS              float32 `json:"ttmEPS"`              // Example: 16.5,
	TTMDividendRate     float32 `json:"ttmDividendRate"`     // Example: 2.25,
	DividendYield       float32 `json:"dividendYield"`       // Example: .021,
	NextDividendDate    string  `json:"nextDividendDate"`    // Example: '2019-03-01',
	ExDividendDate      string  `json:"exDividendDate"`      // Example: '2019-02-08',
	NextEarningsDate    string  `json:"nextEarningsDate"`    // Example: '2019-01-01',
	PERatio             float32 `json:"peRatio"`             // Example: 14,
	Beta                float32 `json:"beta"`                // Example: 1.25,
	MaxChangePercent    float32 `json:"maxChangePercent"`    // Example: 153.021,
	Year5ChangePercent  float32 `json:"year5ChangePercent"`  // Example: 0.5902546932200027,
	Year2ChangePercent  float32 `json:"year2ChangePercent"`  // Example: 0.3777449874142869,
	Year1ChangePercent  float32 `json:"year1ChangePercent"`  // Example: 0.39751716851558366,
	YtdChangePercent    float32 `json:"ytdChangePercent"`    // Example: 0.36659492036160124,
	Month6ChangePercent float32 `json:"month6ChangePercent"` // Example: 0.12208398133748043,
	Month3ChangePercent float32 `json:"month3ChangePercent"` // Example: 0.08466584665846649,
	Month1ChangePercent float32 `json:"month1ChangePercent"` // Example: 0.009668596145283263,
	Day30ChangePercent  float32 `json:"day30ChangePercent"`  // Example: -0.002762605699968781,
	Day5ChangePercent   float32 `json:"day5ChangePercent"`   // Example: -0.005762605699968781
}

// PriceResponse is the datatype returned by the Previous Day Price Endpoint, among others
type PriceResponse struct {
	Date             string  `json:"date"`           // Example: "2019-03-25",
	Minute           string  `json:"minute"`         // Example: "13:30",
	Open             float32 `json:"open"`           // Example: 191.51,
	Close            float32 `json:"close"`          // Example: 188.74,
	High             float32 `json:"high"`           // Example: 191.98,
	Low              float32 `json:"low"`            // Example: 186.6,
	Volume           float32 `json:"volume"`         // Example: 43845293,
	UnadjustedOpen   float32 `json:"uOpen"`          // Example: 191.51,
	UnadjustedClose  float32 `json:"uClose"`         // Example: 188.74,
	UnadjustedHigh   float32 `json:"uHigh"`          // Example: 191.98,
	UnadjustedLow    float32 `json:"uLow"`           // Example: 186.6,
	UnadjustedVolume float32 `json:"uVolume"`        // Example: 43845293,
	Change           float32 `json:"change"`         // Example: 0,
	ChangePercent    float32 `json:"changePercent"`  // Example: 0,
	ChangeOverTime   float32 `json:"changeOverTime"` // Example: 0,
	Symbol           string  `json:"symbol"`         // Example: "AAPL"
}

// PriceTargetResponse is returned by the Price Target Endpoiunt
type PriceTargetResponse struct {
	Symbol             string  `json:"symbol"`             // Example: "AAPL",
	UpdatedDate        string  `json:"updatedDate"`        // Example: "2019-01-30",
	PriceTargetAverage float32 `json:"priceTargetAverage"` // Example: 178.59,
	PriceTargetHigh    float32 `json:"priceTargetHigh"`    // Example: 245,
	PriceTargetLow     float32 `json:"priceTargetLow"`     // Example: 140,
	NumberOfAnalysts   float32 `json:"numberOfAnalysts"`   // Example: 34
}

// QuoteResponse is the datatype returned by the Quote Endpoint
type QuoteResponse struct {
	Symbol                string  `json:"symbol"`                // Example: "AAPL"
	CompanyName           string  `json:"companyName"`           // Example: "Apple Inc."
	CalculationPrice      string  `json:"calculationPrice"`      // Example: "tops"
	Open                  float32 `json:"open"`                  // Example: 154
	OpenTime              float32 `json:"openTime"`              // Example: 1506605400394
	Close                 float32 `json:"close"`                 // Example: 153.28
	CloseTime             float32 `json:"closeTime"`             // Example: 1506605400394
	High                  float32 `json:"high"`                  // Example: 154.80
	Low                   float32 `json:"low"`                   // Example: 153.25
	LatestPrice           float32 `json:"latestPrice"`           // Example: 158.73
	LatestSource          string  `json:"latestSource"`          // Example: "Previous close"
	LatestTime            string  `json:"latestTime"`            // Example: "September 19 2017"
	LatestUpdate          float32 `json:"latestUpdate"`          // Example: 1505779200000
	LatestVolume          float32 `json:"latestVolume"`          // Example: 20567140
	IEXRealtimePrice      float32 `json:"iexRealtimePrice"`      // Example: 158.71
	IEXRealtimeSize       float32 `json:"iexRealtimeSize"`       // Example: 100
	IEXLastUpdated        float32 `json:"iexLastUpdated"`        // Example: 1505851198059
	DelayedPrice          float32 `json:"delayedPrice"`          // Example: 158.71
	DelayedPriceTime      float32 `json:"delayedPriceTime"`      // Example: 1505854782437
	ExtendedPrice         float32 `json:"extendedPrice"`         // Example: 159.21
	ExtendedChange        float32 `json:"extendedChange"`        // Example: -1.68
	ExtendedChangePercent float32 `json:"extendedChangePercent"` // Example: -0.0125
	ExtendedPriceTime     float32 `json:"extendedPriceTime"`     // Example: 1527082200361
	PreviousClose         float32 `json:"previousClose"`         // Example: 158.73
	Change                float32 `json:"change"`                // Example: -1.67
	ChangePercent         float32 `json:"changePercent"`         // Example: -0.01158
	IEXMarketPercent      float32 `json:"iexMarketPercent"`      // Example: 0.00948
	IEXVolume             float32 `json:"iexVolume"`             // Example: 82451
	AvgTotalVolume        float32 `json:"avgTotalVolume"`        // Example: 29623234
	IEXBidPrice           float32 `json:"iexBidPrice"`           // Example: 153.01
	IEXBidSize            float32 `json:"iexBidSize"`            // Example: 100
	IEXAskPrice           float32 `json:"iexAskPrice"`           // Example: 158.66
	IEXAskSize            float32 `json:"iexAskSize"`            // Example: 100
	MarketCap             float32 `json:"marketCap"`             // Example: 751627174400
	Week52High            float32 `json:"week52High"`            // Example: 159.65
	Week52Low             float32 `json:"week52Low"`             // Example: 93.63
	YTDChange             float32 `json:"ytdChange"`             // Example: 0.3665
	PERatio               float32 `json:"peRatio"`               // Example: 17.18
	USMarketOpen          bool    `json:"isUSMarketOpen"`        // Example: false
}
