package iex

import "testing"

func TestKeyStats(t *testing.T) {
	tt := []struct {
		Name   string
		Symbol string
		Error  error
	}{
		{Name: "Valid symbol", Symbol: "AAPL", Error: nil},
		{Name: "Invalid symbol", Symbol: "INVALID", Error: ErrNotFound},
	}

	service, err := validTestService()
	if err != nil {
		t.Fatal("Could not generate a valid test service")
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := service.KeyStats(tc.Symbol)
			if err != tc.Error {
				t.Errorf("Incorrect error returned. Expected '%v', got '%v'", tc.Error, err)
			}
		})
	}
}

func TestPreviousDayPrice(t *testing.T) {
	tt := []struct {
		Name   string
		Symbol string
		Error  error
	}{
		{Name: "Valid symbol", Symbol: "AAPL", Error: nil},
		{Name: "Invalid symbol", Symbol: "INVALID", Error: ErrNotFound},
	}

	service, err := validTestService()
	if err != nil {
		t.Fatal("Could not generate a valid test service")
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := service.PreviousDayPrice(tc.Symbol)
			if err != tc.Error {
				t.Errorf("Incorrect error returned. Expected '%v', got '%v'", tc.Error, err)
			}
		})
	}
}
func TestHistoricalPrices(t *testing.T) {
	tt := []struct {
		Name   string
		Symbol string
		Range  string
		Error  error
	}{
		{Name: "Valid symbol", Symbol: "AAPL", Error: nil, Range: "5dm"},
		{Name: "Invalid symbol", Symbol: "INVALID", Error: ErrNotFound, Range: "1m"},
		{Name: "Invalid range", Symbol: "INVALID", Error: ErrNotFound, Range: "x"},
	}

	service, err := validTestService()
	if err != nil {
		t.Fatal("Could not generate a valid test service")
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := service.HistoricalPrices(tc.Symbol, tc.Range)
			if err != tc.Error {
				t.Errorf("Incorrect error returned. Expected '%v', got '%v'", tc.Error, err)
			}
		})
	}
}

func TestQuote(t *testing.T) {
	tt := []struct {
		Name   string
		Symbol string
		Error  error
	}{
		{Name: "Valid symbol", Symbol: "AAPL", Error: nil},
		{Name: "Invalid symbol", Symbol: "INVALID", Error: ErrNotFound},
	}

	service, err := validTestService()
	if err != nil {
		t.Fatal("Could not generate a valid test service")
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := service.Quote(tc.Symbol)
			if err != tc.Error {
				t.Errorf("Incorrect error returned. Expected '%v', got '%v'", tc.Error, err)
			}
		})
	}
}
