package handler

import (
	"testing"
	"time"
)

func TestSetPriceInCache(t *testing.T) {
	h := Handler{quotes: make(map[string]quote)}
	h.setPriceInCache("AAPL", 10)

	res, ok := h.quotes["AAPL"]
	if !ok {
		t.Fatal("Expected the symbol to be set in the cache but none was")
	}

	if time.Now().Sub(res.cachedAt).Seconds() > 1 {
		t.Errorf("Expected the cachedAt to be %v but was actually %v", time.Now(), res.cachedAt)
	}

	if res.price != 10 {
		t.Errorf("Expected the cachedAt to be %v but was actually %v", 10, res.price)
	}
}
func TestLoadPriceFromCache(t *testing.T) {
	tt := []struct {
		Name   string
		Cache  map[string]quote
		Symbol string
		Error  error
		Quote  int32
	}{
		{Name: "No quotes",
			Symbol: "AAPL",
			Cache:  map[string]quote{},
			Error:  errNoCacheResult,
			Quote:  0,
		},
		{
			Name:   "Valid quote",
			Symbol: "AAPL",
			Cache: map[string]quote{
				"AAPL": quote{price: 10, cachedAt: time.Now()},
			},
			Error: nil,
			Quote: 10,
		},
		{
			Name:   "Expired quote",
			Symbol: "AAPL",
			Cache: map[string]quote{
				"AAPL": quote{price: 10, cachedAt: time.Now().Add(time.Hour * -1)},
			},
			Error: errNoCacheResult,
			Quote: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			h := Handler{quotes: tc.Cache}
			quote, err := h.loadPriceFromCache(tc.Symbol)

			if err != tc.Error {
				t.Errorf("Expected the error to be %v but was actually %v", tc.Error, err)
			}

			if quote != tc.Quote {
				t.Errorf("Expected the quote to be %v but was actually %v", tc.Quote, quote)
			}
		})
	}
}
