package handler

import (
	"context"
	"fmt"
	"time"

	proto "github.com/kytra-app/charts-api/proto"
	iex "github.com/kytra-app/helpers/iex-cloud"
	portfolioTracking "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	iex               iex.Service
	stocks            stocks.StocksService
	portfolios        portfolios.PortfoliosService
	portfolioTracking portfolioTracking.PortfolioValueTrackingService
}

// New returns an instance of Handler
func New(iex iex.Service, client client.Client) Handler {
	return Handler{
		iex:               iex,
		stocks:            stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		portfolios:        portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		portfolioTracking: portfolioTracking.NewPortfolioValueTrackingService("kytra-srv-v1-portfolio-value-tracking:8080", client),
	}
}

// GetStock takes a Stock UUID, and returns a chart for that period
func (h Handler) GetStock(ctx context.Context, req *proto.Request, rsp *proto.Chart) error {
	if req.Uuid == "" {
		return errors.BadRequest("INVALID_UUID", "A UUID is required")
	}

	sRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: req.Uuid})
	if err != nil {
		return err
	}

	r := req.Range
	if r == "1d" {
		r = "date/" + time.Now().Format("20060102")
	}

	prices, err := h.iex.HistoricalPrices(sRsp.Stock.Symbol, r)
	if err != nil {
		return err
	}

	rsp.Points = make([]*proto.Point, len(prices))
	for i, p := range prices {
		rsp.Points[i] = &proto.Point{
			Date:   fmt.Sprintf("%v %v", p.Date, p.Minute),
			Volume: p.Volume,
			Close:  p.Close,
		}
	}

	return nil
}

// GetPortfolio takes a Portfolio UUID, and returns a chart of it's value
// TODO: Remove the dynamic UUID lookup, only return the active users portfolio.
func (h Handler) GetPortfolio(ctx context.Context, req *proto.Request, rsp *proto.Chart) error {
	var err error
	var p *portfolioTracking.Portfolio

	if req.Range == "1d" {
		p, err = h.portfolioTracking.GetIntradayHistory(ctx, &portfolioTracking.Portfolio{Uuid: req.Uuid})
	} else {
		p, err = h.portfolioTracking.GetDailyHistory(ctx, &portfolioTracking.Portfolio{Uuid: req.Uuid})
	}

	if err != nil {
		return err
	}

	rsp.Points = make([]*proto.Point, len(p.History))
	for i, p := range p.History {
		rsp.Points[i] = &proto.Point{
			Date:  time.Unix(p.Date, 0).String(),
			Value: float32(p.Amount) / 100,
		}
	}

	return nil
}

// GetInvestor returns the chart for a investors portfolio.
// TODO: Once live portfolios are enabled, this API will need to be authenticated
// to prevent users portfolio balances being exposed.
func (h Handler) GetInvestor(ctx context.Context, req *proto.Request, rsp *proto.Chart) error {
	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: req.Uuid})
	if err != nil {
		return err
	}

	p, err := h.portfolioTracking.GetDailyHistory(ctx, &portfolioTracking.Portfolio{Uuid: portfolio.Uuid})
	if err != nil {
		return err
	}

	rsp.Points = make([]*proto.Point, len(p.History))
	for i, p := range p.History {
		rsp.Points[i] = &proto.Point{
			Date:  time.Unix(p.Date, 0).String(),
			Value: float32(p.Amount) / 100,
		}
	}

	return nil
}
