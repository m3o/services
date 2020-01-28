package handler

import (
	"context"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	iex "github.com/kytra-app/helpers/iex-cloud"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	proto "github.com/kytra-app/trades-api/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth       auth.Authenticator
	iex        iex.Service
	stocks     stocks.StocksService
	portfolios portfolios.PortfoliosService
	trades     trades.TradesService
	followers  followers.FollowersService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, iex iex.Service, client client.Client) Handler {
	return Handler{
		auth:       auth,
		iex:        iex,
		stocks:     stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades:     trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		followers:  followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
	}
}

// GetAsset retrieves the trades for a given users portfolio, relating to the requested asset.
func (h Handler) GetAsset(ctx context.Context, req *proto.Asset, rsp *proto.Asset) error {
	portfolioUUID, err := h.portfolioUUIDForUser(ctx)
	if err != nil {
		return err
	}

	// Lookup the stock
	stockRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: req.Uuid})
	if err != nil {
		return err
	}

	// Lookup the current price
	quote, err := h.iex.Quote(stockRsp.Stock.Symbol)
	if err != nil {
		return err
	}
	currentPrice := int64(quote.LatestPrice * 100)

	// Lookup the trades
	tradesRsp, err := h.trades.ListTradesForPosition(ctx, &trades.ListRequest{
		IncludeMetadata: true,
		PortfolioUuid:   portfolioUUID,
		Asset:           &trades.Asset{Uuid: req.Uuid, Type: req.Type},
	})
	if err != nil {
		return err
	}

	// Serialize the data
	trades := make([]*proto.Trade, len(tradesRsp.Trades))
	for i, trade := range tradesRsp.Trades {
		trades[i] = serializeTrade(trade)
	}

	// Return the result
	*rsp = proto.Asset{
		Uuid:            req.Uuid,
		Type:            req.Type,
		Trades:          trades,
		BookCost:        tradesRsp.BookCost,
		CurrentQuantity: tradesRsp.Quantity,
		CurrentValue:    tradesRsp.Quantity * currentPrice,
	}
	return nil
}

// CreateTrade executes a trade in the users portfolio
func (h Handler) CreateTrade(ctx context.Context, req *proto.Trade, rsp *proto.Trade) error {
	// Validate the request
	if req.Asset == nil || req.Asset.Type != "Stock" {
		return errors.BadRequest("INVALID_ASSET", "The asset requested is invalid")

	}

	// Authorise the user
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	// Find the users portfolio
	portfolioUUID, err := h.portfolioUUIDForUser(ctx)
	if err != nil {
		return err
	}

	// Determine the trade type
	var tradeType trades.TradeType
	switch req.Type {
	case "BUY":
		tradeType = trades.TradeType_BUY
	case "SELL":
		tradeType = trades.TradeType_SELL
	}

	// Construct the trade
	trade := &trades.Trade{
		Asset:         &trades.Asset{Uuid: req.Asset.Uuid, Type: req.Asset.Type},
		Quantity:      req.Quantity,
		PortfolioUuid: portfolioUUID,
		Type:          tradeType,
		ClientUuid:    req.ClientUuid,
	}

	// Execute the trade
	if trade, err = h.trades.CreateTrade(ctx, trade); err != nil {
		return err
	}

	// Follow the stock
	h.followers.Follow(ctx, &followers.Request{
		Followee: &followers.Resource{Uuid: req.Asset.Uuid, Type: req.Asset.Type},
		Follower: &followers.Resource{Uuid: user.UUID, Type: "User"},
	})

	// Serialize the result
	*rsp = *serializeTrade(trade)
	return nil
}

// SetTradeMetadata updates the metadata for the given trade
func (h Handler) SetTradeMetadata(ctx context.Context, req *proto.Trade, rsp *proto.Trade) error {
	portfolioUUID, err := h.portfolioUUIDForUser(ctx)
	if err != nil {
		return err
	}

	// Ensure the user has permissions on this trade
	trade, err := h.trades.GetTrade(ctx, &trades.Trade{Uuid: req.Uuid})
	if err != nil {
		return err
	} else if trade.PortfolioUuid != portfolioUUID {
		return errors.Forbidden("FORBIDDEN_TRADE", "This trade does not belong to your portfolio")
	}

	// Execute the update
	query := trades.Trade{Uuid: req.Uuid, TargetPrice: req.TargetPrice, Notes: req.Notes}
	trade, err = h.trades.SetTradeMetadata(ctx, &query)
	if err != nil {
		return err
	}

	// Serialize the result
	*req = *serializeTrade(trade)
	return nil
}

func (h Handler) portfolioUUIDForUser(ctx context.Context) (string, error) {
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return "", errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	porfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.UUID})
	if err != nil {
		return "", err
	}

	return porfolio.Uuid, nil
}

func serializeTrade(trade *trades.Trade) *proto.Trade {
	return &proto.Trade{
		Uuid:        trade.Uuid,
		Quantity:    trade.Quantity,
		Type:        trade.Type.String(),
		UnitPrice:   trade.UnitPrice,
		TargetPrice: trade.TargetPrice,
		Notes:       trade.Notes,
		CreatedAt:   time.Unix(trade.CreatedAt, 0).String(),
	}
}
