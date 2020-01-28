package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/trades-srv/helpers"
	proto "github.com/kytra-app/trades-srv/proto"
	"github.com/kytra-app/trades-srv/storage"

	"github.com/kytra-app/helpers/microtime"
	ledger "github.com/kytra-app/ledger-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(iex iex.Service, storage storage.Service, broker broker.Broker, client client.Client) *Handler {
	return &Handler{
		iex:    iex,
		db:     storage,
		broker: broker,
		stocks: stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		ledger: ledger.NewLedgerService("kytra-srv-v1-ledger:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	iex    iex.Service
	db     storage.Service
	broker broker.Broker
	stocks stocks.StocksService
	ledger ledger.LedgerService
}

// CreateTrade executes a trade once deducting the cash balance form the portfolios ledger
func (h *Handler) CreateTrade(ctx context.Context, req *proto.Trade, rsp *proto.Trade) error {
	if req.Asset == nil || req.Asset.Type != "Stock" {
		return errors.BadRequest("INVALID_ASSET", "The asset type provided is invalid")
	}

	// Get Asset (currently only Stocks)
	stockRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: req.Asset.Uuid})
	if err != nil {
		return err
	}

	// Get Price from IEX
	quote, err := h.iex.Quote(stockRsp.Stock.Symbol)
	if err != nil {
		return err
	}
	unitPrice := int64(quote.LatestPrice * 100)

	// Construct the trade
	trade := storage.Trade{
		UUID:          req.Uuid,
		ClientUUID:    req.ClientUuid,
		Type:          req.Type.String(),
		AssetUUID:     req.Asset.Uuid,
		AssetType:     req.Asset.Type,
		PortfolioUUID: req.PortfolioUuid,
		Quantity:      req.Quantity,
		UnitPrice:     unitPrice,
	}

	// Pre-validate Trade (Will this result in a short position? Is the object valid?)
	if err := h.db.PrevalidateTrade(trade); err != nil {
		return err
	}

	// Create Ledger Transaction
	var ledgerTransType ledger.TransactionType
	switch req.Type {
	case proto.TradeType_BUY:
		ledgerTransType = ledger.TransactionType_BUY_ASSET
	case proto.TradeType_SELL:
		ledgerTransType = ledger.TransactionType_SELL_ASSET
	}

	transaction := &ledger.Transaction{
		PortfolioUuid: req.PortfolioUuid,
		Amount:        unitPrice * req.Quantity,
		Type:          ledgerTransType,
	}
	transaction, err = h.ledger.CreateTransaction(ctx, transaction)
	if err != nil {
		return err
	}

	// Create Trade
	trade, err = h.db.CreateTrade(trade)
	if err != nil {
		return err
	}

	// Serialize trade
	serializedTrade := h.serializeTrade(trade, false)
	*rsp = serializedTrade

	// TODO: IF an error occurs, void the ledger transaction
	_ = transaction

	bytes, err := json.Marshal(&serializedTrade)
	if err != nil {
		return err
	}
	brokerErr := h.broker.Publish("kytra-srv-v1-trades-trade-created", &broker.Message{Body: bytes})
	if brokerErr != nil {
		fmt.Printf("Error Sending Msg to broker: %v\n", err)
	} else {
		fmt.Printf("Message sent to broker\n")
	}

	return nil
}

// ListTrades fetches the trades placed over a time window
func (h *Handler) ListTrades(ctx context.Context, req *proto.ListTradesRequest, rsp *proto.ListTradesResponse) error {
	if req.StartTime == 0 || req.EndTime == 0 {
		return errors.BadRequest("MISSING_TIME", "A start time and end time is required")
	}

	startTime := time.Unix(req.StartTime, 0)
	endTime := time.Unix(req.EndTime, 0)

	trades, err := h.db.ListTrades(startTime, endTime, req.PortfolioUuids)
	if err != nil {
		return err
	}
	rsp.Trades = h.serializeTrades(trades, false)

	return nil
}

// GetTrade looksup a trade given a UUID
func (h *Handler) GetTrade(ctx context.Context, req *proto.Trade, rsp *proto.Trade) error {
	// Validate the request
	if req.Uuid == "" {
		return errors.BadRequest("MISSING_UUID", "A trade UUID is required")
	}

	// Query the DB
	trade, err := h.db.GetTrade(storage.Trade{UUID: req.Uuid})
	if err != nil {
		return err
	}

	// Serialize the data
	*rsp = h.serializeTrade(trade, true)
	return nil
}

// SetTradeMetadata sets the notes and target price for a trade, looked up using the UUID
func (h *Handler) SetTradeMetadata(ctx context.Context, req *proto.Trade, rsp *proto.Trade) error {
	// Validate the request
	if req.Uuid == "" {
		return errors.BadRequest("MISSING_UUID", "A trade UUID is required")
	}

	// Get the trade
	trade, err := h.db.GetTrade(storage.Trade{UUID: req.Uuid})
	if err != nil {
		return err
	}

	// Update the DB
	query := storage.Trade{UUID: req.Uuid, Notes: req.Notes, TargetPrice: req.TargetPrice}
	trade, err = h.db.SetTradeMetadata(query)
	if err != nil {
		return err
	}

	// Serialize the data
	*rsp = h.serializeTrade(trade, true)
	return nil
}

// ListTradesForPosition returns all the trades made for the position (asset / portfolio relationship)
func (h *Handler) ListTradesForPosition(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Validate the request
	if req.PortfolioUuid == "" {
		return errors.BadRequest("MISSING_PORTFOLIO_UUID", "A portfolio UUID is required")
	} else if req.Asset == nil {
		return errors.BadRequest("MISSING_ASSET", "An asset is required")
	} else if req.Asset.Uuid == "" || req.Asset.Type == "" {
		return errors.BadRequest("MISSING_ASSET", "An asset is required")
	}

	// Check for a custom time
	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	// Query the DB
	position := storage.Position{
		AssetUUID:     req.Asset.Uuid,
		AssetType:     req.Asset.Type,
		PortfolioUUID: req.PortfolioUuid,
	}
	trades, err := h.db.ListTradesForPosition(position, time)
	if err != nil {
		return err
	}

	// Serialize the data
	rsp.BookCost = helpers.BookCost(trades)
	rsp.Quantity = helpers.SumQuantity(trades)
	rsp.Trades = h.serializeTrades(trades, req.IncludeMetadata)
	return nil
}

// ListTradesForPortfolio returns all the trades made within the given portfolio
func (h *Handler) ListTradesForPortfolio(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Validate the request
	if req.PortfolioUuid == "" {
		return errors.BadRequest("MISSING_PORTFOLIO_UUID", "A portfolio UUID is required")
	}

	// Check for a custom time
	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	// Query the DB
	trades, err := h.db.ListTradesForPortfolio(req.PortfolioUuid, time)
	if err != nil {
		return err
	}

	// Serialize the data
	rsp.BookCost = helpers.BookCost(trades)
	rsp.Quantity = helpers.SumQuantity(trades)
	rsp.Trades = h.serializeTrades(trades, req.IncludeMetadata)
	return nil
}

// ListPositions returns all the active positions for the given portfolio
func (h *Handler) ListPositions(ctx context.Context, req *proto.BulkListRequest, rsp *proto.ListResponse) error {
	// Validate the request
	if len(req.PortfolioUuids) == 0 {
		return errors.BadRequest("MISSING_PORTFOLIO_UUIDS", "One or more portfolio UUIDs are required")
	}

	// Check for a custom time
	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	// Query the DB
	positions, err := h.db.ListPositions(req.PortfolioUuids, req.AssetType, req.AssetUuids, time)
	if err != nil {
		return err
	}

	// Serialize the data
	rsp.Positions = make([]*proto.Position, len(positions))
	for i, pos := range positions {
		rsp.Positions[i] = &proto.Position{
			Quantity:      pos.Quantity,
			PortfolioUuid: pos.PortfolioUUID,
			BookCost:      pos.BookCost,
			Asset:         &proto.Asset{Type: pos.AssetType, Uuid: pos.AssetUUID},
		}
	}
	return nil
}

// ListPositionsForPortfolio returns all the active positions for the given portfolio
func (h *Handler) ListPositionsForPortfolio(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Validate the request
	if req.PortfolioUuid == "" {
		return errors.BadRequest("MISSING_PORTFOLIO_UUID", "A portfolio UUID is required")
	}

	// Check for a custom time
	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	// Query the DB
	positions, err := h.db.ListPositionsForPortfolio(req.PortfolioUuid, time)
	if err != nil {
		return err
	}

	// Serialize the data
	rsp.Positions = make([]*proto.Position, len(positions))
	for i, pos := range positions {
		rsp.Positions[i] = &proto.Position{
			Quantity:      pos.Quantity,
			PortfolioUuid: pos.PortfolioUUID,
			BookCost:      pos.BookCost,
			Asset:         &proto.Asset{Type: pos.AssetType, Uuid: pos.AssetUUID},
		}
	}
	return nil
}

// AllAssets returns all the assets that have ever been trades
func (h *Handler) AllAssets(ctx context.Context, req *proto.AllAssetsRequest, rsp *proto.AllAssetsResponse) error {
	assets, err := h.db.AllAssets()
	if err != nil {
		return err
	}

	rsp.Assets = make([]*proto.Asset, len(assets))
	for i, a := range assets {
		rsp.Assets[i] = &proto.Asset{Uuid: a.UUID, Type: a.Type}
	}

	return nil
}

func (h *Handler) serializeTrades(trades []storage.Trade, includeMetadata bool) []*proto.Trade {
	res := make([]*proto.Trade, len(trades))
	for i, trade := range trades {
		x := h.serializeTrade(trade, includeMetadata)
		res[i] = &x
	}
	return res
}

func (h *Handler) serializeTrade(trade storage.Trade, includeMetadata bool) proto.Trade {
	var tradeType proto.TradeType
	switch trade.Type {
	case "BUY":
		tradeType = proto.TradeType_BUY
	case "SELL":
		tradeType = proto.TradeType_SELL
	}

	res := proto.Trade{
		Uuid:          trade.UUID,
		Type:          tradeType,
		PortfolioUuid: trade.PortfolioUUID,
		Quantity:      trade.Quantity,
		UnitPrice:     trade.UnitPrice,
		CreatedAt:     trade.CreatedAt.Unix(),
		Asset:         &proto.Asset{Type: trade.AssetType, Uuid: trade.AssetUUID},
	}

	if includeMetadata {
		res.Notes = trade.Notes
		res.TargetPrice = trade.TargetPrice
	}

	return res
}
