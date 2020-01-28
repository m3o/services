package handler

import (
	"context"

	followers "github.com/kytra-app/followers-srv/proto"
	proto "github.com/kytra-app/insights-summary-srv/proto"
	valuation "github.com/kytra-app/portfolio-valuation-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	earnings "github.com/kytra-app/stock-earnings-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	target "github.com/kytra-app/stock-target-price-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
)

// New returns an instance of Handler
func New(client client.Client) *Handler {
	return &Handler{
		users:      users.NewUsersService("kytra-srv-v1-users:8080", client),
		stocks:     stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades:     trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		quotes:     quotes.NewStockQuoteService("kytra-srv-v2-stock-quote:8080", client),
		targets:    target.NewStockTargetPriceService("kytra-srv-v1-stock-target-price:8080", client),
		earnings:   earnings.NewStockEarningsService("kytra-srv-v1-stock-earnings:8080", client),
		valuation:  valuation.NewPortfolioValuationService("kytra-srv-v1-portfolio-valuation:8080", client),
		followers:  followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	users      users.UsersService
	stocks     stocks.StocksService
	trades     trades.TradesService
	quotes     quotes.StockQuoteService
	targets    target.StockTargetPriceService
	earnings   earnings.StockEarningsService
	valuation  valuation.PortfolioValuationService
	followers  followers.FollowersService
	portfolios portfolios.PortfoliosService
}

// List returns a summary for the insights on a given day
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) (err error) {
	data := summaryData{
		Handler:    h,
		Context:    ctx,
		UserUUID:   req.GetUserUuid(),
		StockUUIDs: req.GetAssetUuids(),
	}
	if err := data.Fetch(); err != nil {
		return err
	}

	rsp.Assets = make([]*proto.Asset, len(req.AssetUuids))
	for i, uuid := range req.AssetUuids {
		rsp.Assets[i] = &proto.Asset{
			Uuid: uuid, Type: req.AssetType, Summary: data.Stringify(uuid),
		}
	}

	return nil
}

// Get returns the summary for a single stock on a given day
func (h *Handler) Get(ctx context.Context, req *proto.GetRequest, rsp *proto.Asset) error {
	data := summaryData{
		Handler:    h,
		Context:    ctx,
		UserUUID:   req.GetUserUuid(),
		StockUUIDs: []string{req.GetAssetUuid()},
	}
	if err := data.Fetch(); err != nil {
		return err
	}

	*rsp = proto.Asset{
		Uuid:    req.GetAssetUuid(),
		Type:    req.AssetType,
		Summary: data.Stringify(req.GetAssetUuid()),
	}

	return nil
}

type summaryData struct {
	// Input
	Handler    *Handler
	Context    context.Context
	UserUUID   string
	StockUUIDs []string

	// Processing
	AllUsers            []*users.User
	AllPositions        []*trades.Position
	AllPortfolios       []*portfolios.Portfolio
	TotalPortfolioValue int64
	SectorAllocations   map[string]int64

	// Ready for use (grouped by stock UUIID)
	StkQuote             map[string]*quotes.Quote
	StkMetadata          map[string]*stocks.Stock
	StkHasEarning        map[string]bool
	StkPriceTarget       map[string]*target.Stock
	StkUsersWithPosition map[string][]*users.User
}
