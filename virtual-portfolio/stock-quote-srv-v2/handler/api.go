package handler

import (
	"context"

	"github.com/kytra-app/helpers/microtime"
	proto "github.com/kytra-app/stock-quote-srv/proto"
)

// GetQuote the a stock and finds the latest quote for it
func (h *Handler) GetQuote(ctx context.Context, req *proto.Stock, rsp *proto.Quote) error {
	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	quote, err := h.db.Get(time, req.Uuid, req.IncludeOutOfHours)
	if err != nil {
		return err
	}

	*rsp = proto.Quote{
		StockUuid:        quote.StockUUID,
		Price:            int64(quote.Price),
		CreatedAt:        quote.CreatedAt.Unix(),
		PercentageChange: quote.ChangePercent,
		MarketClosed:     !quote.MarketOpen,
	}

	return nil
}

// ListQuotes gets a proto.ListRequest (which has a slice of Uuids) and finds the latest quote for them
func (h *Handler) ListQuotes(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	rsp.Quotes = make([]*proto.Quote, len(req.Uuids))

	time, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	quotes, err := h.db.List(time, req.Uuids, req.IncludeOutOfHours)
	if err != nil {
		return err
	}

	rsp.Quotes = make([]*proto.Quote, len(quotes))
	for i, quote := range quotes {
		rsp.Quotes[i] = &proto.Quote{
			StockUuid:        quote.StockUUID,
			Price:            int64(quote.Price),
			CreatedAt:        quote.CreatedAt.Unix(),
			PercentageChange: quote.ChangePercent,
			MarketClosed:     !quote.MarketOpen,
		}
	}

	return nil
}
