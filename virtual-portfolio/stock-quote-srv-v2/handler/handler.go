package handler

import (
	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/helpers/worldtradingdata"
	insights "github.com/kytra-app/insights-srv/proto"
	storage "github.com/kytra-app/stock-quote-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	"github.com/micro/go-micro/client"
)

// New returns an instance of Handler
func New(wtd worldtradingdata.Service, iex iex.Service, db storage.Service, client client.Client) *Handler {
	return &Handler{
		db:       db,
		wtd:      wtd,
		iex:      iex,
		trades:   trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		stocks:   stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		insights: insights.NewInsightsService("kytra-srv-v1-insights:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	db       storage.Service
	iex      iex.Service
	wtd      worldtradingdata.Service
	trades   trades.TradesService
	stocks   stocks.StocksService
	insights insights.InsightsService
}
