package handler

import (
	"context"

	proto "github.com/kytra-app/portfolio-allocation-srv/proto"
	valuation "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(client client.Client) *Handler {
	return &Handler{
		stocks:     stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades:     trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		quotes:     quotes.NewStockQuoteService("kytra-srv-v2-stock-quote:8080", client),
		valuation:  valuation.NewPortfolioValueTrackingService("kytra-srv-v1-portfolio-value-tracking:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	stocks     stocks.StocksService
	trades     trades.TradesService
	quotes     quotes.StockQuoteService
	valuation  valuation.PortfolioValueTrackingService
	portfolios portfolios.PortfoliosService
}

// List returns a list of portfolios and their allocations
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) (err error) {
	if len(req.Uuids) == 0 && len(req.UserUuids) == 0 {
		return errors.BadRequest("MISSING_UUIDS", "One or more uuids or user_uuids are required")
	}

	pRsp, err := h.portfolios.List(ctx, &portfolios.ListRequest{
		Uuids: req.Uuids, UserUuids: req.UserUuids,
	})
	if err != nil {
		return err
	}

	rsp.Portfolios, err = h.serializePortfolios(ctx, pRsp.GetPortfolios())
	if err != nil {
		return err
	}

	return nil
}

// Get returns a single portfolio and its allocation
func (h *Handler) Get(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if req.Uuid == "" && req.UserUuid == "" {
		return errors.BadRequest("MISSING_UUID", "A uuid is required")
	}

	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{
		Uuid: req.Uuid, UserUuid: req.UserUuid,
	})
	if err != nil {
		return err
	}

	data, err := h.serializePortfolios(ctx, []*portfolios.Portfolio{portfolio})
	if err != nil {
		return err
	}
	if len(data) != 1 {
		return errors.InternalServerError("MISSING_DATA", "An invalid number of portfolios were returned")
	}

	*rsp = *data[0]
	return nil
}
