package handler

import (
	"context"
	"fmt"
	"time"

	proto "github.com/kytra-app/charts-api/proto"
	auth "github.com/kytra-app/helpers/authentication"
	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/helpers/microtime"
	portfolioTracking "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	iex               iex.Service
	auth              auth.Authenticator
	stocks            stocks.StocksService
	portfolios        portfolios.PortfoliosService
	portfolioTracking portfolioTracking.PortfolioValueTrackingService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, iex iex.Service, client client.Client) Handler {
	return Handler{
		iex:               iex,
		auth:              auth,
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

	if req.Range == "1d" {
		today := time.Now().Truncate(time.Hour * 24)
		rsp.MinTime = today.Add(14 * time.Hour).Add(30 * time.Minute).Unix()
		rsp.MaxTime = today.Add(21 * time.Hour).Unix()

		query := "date/" + today.Format("20060102")
		prices, err := h.iex.HistoricalPrices(sRsp.Stock.Symbol, query, false)
		if err != nil {
			return err
		}

		rsp.Points = []*proto.Point{}

		for _, p := range prices {
			timeStr := fmt.Sprintf("%v %v -0500", p.Date, p.Minute)
			t, err := time.Parse("2006-01-02 15:04 -0700", timeStr)

			if err != nil {
				fmt.Printf("Err parsing time %v\n", p)
				return err
			}

			if p.Close != 0 {
				rsp.Points = append(rsp.Points, &proto.Point{
					Time:  t.Unix(),
					Value: p.Close,
				})
			}
		}

		return nil
	}

	prices, err := h.iex.HistoricalPrices(sRsp.Stock.Symbol, req.Range, true)
	if err != nil {
		return err
	}

	rsp.Points = []*proto.Point{}
	for _, p := range prices {
		t, err := time.Parse("2006-01-02", p.Date)
		if err != nil {
			fmt.Printf("Err parsing time %v\n", p)
			continue
		}

		rsp.Points = append(rsp.Points, &proto.Point{
			Time:  t.Add(time.Hour * 5).Unix(),
			Value: p.Close,
		})

		if rsp.MinTime == 0 || rsp.MinTime > t.Unix() {
			rsp.MinTime = t.Unix()
		}

		if rsp.MaxTime < t.Unix() {
			rsp.MaxTime = t.Unix()
		}
	}

	return nil
}

// GetPortfolio takes a Portfolio UUID, and returns a chart of it's value
func (h Handler) GetPortfolio(ctx context.Context, req *proto.Request, rsp *proto.Chart) error {
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}
	date = date.Truncate(time.Hour * 24)

	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return err
	}

	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.UUID})
	if err != nil {
		return err
	}

	var p *portfolioTracking.Portfolio
	if req.Range == "1d" {
		p, err = h.portfolioTracking.GetIntradayHistory(ctx, &portfolioTracking.Portfolio{Uuid: portfolio.Uuid})
	} else {
		p, err = h.portfolioTracking.GetDailyHistory(ctx, &portfolioTracking.Portfolio{Uuid: portfolio.Uuid})
	}
	if err != nil {
		return err
	}

	if req.Range == "1d" {
		rsp.MinTime = date.Add(14 * time.Hour).Add(30 * time.Minute).Unix()
		rsp.MaxTime = date.Add(21 * time.Hour).Unix()
	} else {
		for _, p := range p.GetHistory() {
			if rsp.MinTime == 0 || p.Date < rsp.MinTime {
				rsp.MinTime = p.Date
			}
			if rsp.MaxTime < p.Date {
				rsp.MaxTime = p.Date
			}
		}
	}

	history := []*portfolioTracking.Valuation{}
	for _, h := range p.History {
		if h.Date > rsp.MaxTime {
			continue
		}

		if h.Date < rsp.MinTime {
			continue
		}

		history = append(history, h)
	}

	rsp.Points = make([]*proto.Point, len(history))
	for i, p := range history {
		rsp.Points[i] = &proto.Point{
			Time:  p.Date,
			Value: float32(p.Amount) / 100,
		}
	}

	return nil
}

// GetInvestor returns the chart for a investors portfolio.
// TODO: Once live portfolios are enabled, this API will need to be authenticated
// to prevent users portfolio balances being exposed.
func (h Handler) GetInvestor(ctx context.Context, req *proto.Request, rsp *proto.Chart) error {
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}
	date = date.Truncate(time.Hour * 24)

	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: req.Uuid})
	if err != nil {
		return err
	}

	var p *portfolioTracking.Portfolio
	if req.Range == "1d" {
		p, err = h.portfolioTracking.GetIntradayHistory(ctx, &portfolioTracking.Portfolio{Uuid: portfolio.Uuid})
	} else {
		p, err = h.portfolioTracking.GetDailyHistory(ctx, &portfolioTracking.Portfolio{Uuid: portfolio.Uuid})
	}
	if err != nil {
		return err
	}

	if req.Range == "1d" {
		rsp.MinTime = date.Add(14 * time.Hour).Add(30 * time.Minute).Unix()
		rsp.MaxTime = date.Add(21 * time.Hour).Unix()
	} else {
		for _, p := range p.GetHistory() {
			if rsp.MinTime == 0 || p.Date < rsp.MinTime {
				rsp.MinTime = p.Date
			}
			if rsp.MaxTime < p.Date {
				rsp.MaxTime = p.Date
			}
		}
	}

	history := []*portfolioTracking.Valuation{}
	for _, h := range p.History {
		if h.Date > rsp.MaxTime {
			continue
		}

		if h.Date < rsp.MinTime {
			continue
		}

		history = append(history, h)
	}

	rsp.Points = make([]*proto.Point, len(history))
	for i, p := range history {
		rsp.Points[i] = &proto.Point{
			Time:  p.Date,
			Value: float32(p.Amount) / 100,
		}
	}

	return nil
}
