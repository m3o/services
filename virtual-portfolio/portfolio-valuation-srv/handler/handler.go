package handler

import (
	"context"

	proto "github.com/kytra-app/portfolio-valuation-srv/proto"

	ledger "github.com/kytra-app/ledger-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(client client.Client) *Handler {
	return &Handler{
		stocks: stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades: trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		ledger: ledger.NewLedgerService("kytra-srv-v1-ledger:8080", client),
		quotes: quotes.NewStockQuoteService("kytra-srv-v1-stock-quote:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	stocks stocks.StocksService
	trades trades.TradesService
	ledger ledger.LedgerService
	quotes quotes.StockQuoteService
}

// GetPortfolio calculates the current value of a portfolio, taking into account the cash balance and the
// value of the current shares (calculated by summing of the quanaities * the value, as provided by IEX
func (h Handler) GetPortfolio(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	if req.Uuid == "" {
		return errors.BadRequest("MISSING_UUID", "A UUID is required")
	}

	// Determine the Cash Balance
	cashBalance, err := h.currentCashBalance(ctx, req.Uuid)
	if err != nil {
		return err
	}

	// Find the active positions (Asset & Quantity)
	positions, err := h.getCurrentPositions(ctx, req.Uuid)
	if err != nil {
		return err
	}

	// Value the positions
	positionValue, err := h.valuePositions(ctx, positions)
	if err != nil {
		return err
	}

	// Return the result
	*rsp = proto.Portfolio{
		Uuid:        req.Uuid,
		CashValue:   cashBalance,
		AssetsValue: positionValue,
		TotalValue:  cashBalance + positionValue,
	}

	return nil
}

func (h Handler) currentCashBalance(ctx context.Context, portfolioUUID string) (int64, error) {
	rsp, err := h.ledger.GetPortfolio(ctx, &ledger.Portfolio{Uuid: portfolioUUID})
	if err != nil {
		return 0, err
	}
	return rsp.CurrentBalance, nil
}

func (h Handler) getCurrentPositions(ctx context.Context, portfolioUUID string) ([]*trades.Position, error) {
	rsp, err := h.trades.ListPositionsForPortfolio(ctx, &trades.ListRequest{PortfolioUuid: portfolioUUID})
	if err != nil {
		return []*trades.Position{}, err
	}
	return rsp.Positions, nil
}

func (h Handler) valuePositions(ctx context.Context, positions []*trades.Position) (int64, error) {
	// Fetch the stocks which comprise the positions (we only support stocks at present)
	stockUUIDs := make([]string, len(positions))
	for i, p := range positions {
		stockUUIDs[i] = p.Asset.Uuid
	}
	stocksRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return 0, err
	}

	// Get the stock prices
	symbols := make([]string, len(stocksRsp.Stocks))
	uuidForSymbols := make(map[string]string, len(stocksRsp.Stocks))

	for i, s := range stocksRsp.Stocks {
		symbols[i] = s.Symbol
		uuidForSymbols[s.Symbol] = s.Uuid
	}
	quoteRsp, err := h.quotes.ListQuotes(ctx, &quotes.ListRequest{Symbols: symbols})
	if err != nil {
		return 0, err
	}
	stockPrices := make(map[string]int64, len(stocksRsp.Stocks))
	for _, quote := range quoteRsp.Quotes {
		stockPrices[uuidForSymbols[quote.Symbol]] = int64(quote.Price)
	}

	// Calculate the total value
	var total int64
	for _, position := range positions {
		unitPrice := stockPrices[position.Asset.Uuid]
		price := unitPrice * position.Quantity
		total = total + price
	}

	return total, nil
}
