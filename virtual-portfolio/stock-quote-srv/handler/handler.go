package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	iex "github.com/kytra-app/helpers/iex-cloud"
	proto "github.com/kytra-app/stock-quote-srv/proto"
)

// New returns an instance of Handler
func New(iex iex.Service) *Handler {
	return &Handler{iex: iex, quotes: make(map[string]quote)}
}

// Handler is an object can process RPC requests
type Handler struct {
	iex    iex.Service
	quotes map[string]quote
}

type quote struct {
	cachedAt time.Time
	price    int32
}

// GetQuote gets a proto.Stock and finds the latest quote for it
func (h *Handler) GetQuote(ctx context.Context, req *proto.Stock, rsp *proto.Quote) error {
	price, err := h.fetchPrice(req.Symbol)
	if err != nil {
		return err
	}

	rsp.Price = price
	return nil
}

// ListQuotes gets a proto.ListRequest (which has a slice of symbols) and finds the latest quote for them
func (h *Handler) ListQuotes(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	rsp.Quotes = make([]*proto.Quote, len(req.Symbols))

	for i, symbol := range req.Symbols {
		price, err := h.fetchPrice(symbol)
		if err != nil {
			return err
		}

		rsp.Quotes[i] = &proto.Quote{Symbol: symbol, Price: price}
	}

	return nil
}

const (
	maxCacheMins = 15
)

var (
	errNoCacheResult = errors.New("Quote not found in cache")
)

// RefreshCache gets the lastest prices from IEX for all symbols currently in
// the cache.
func (h *Handler) RefreshCache() error {
	for symbol := range h.quotes {
		price, err := h.loadPriceFromIEX(symbol)
		if err != nil {
			return err
		}
		h.setPriceInCache(symbol, price)
	}
	return nil
}

func (h *Handler) fetchPrice(symbol string) (int32, error) {
	if price, err := h.loadPriceFromCache(symbol); err == nil {
		return price, nil
	}

	price, err := h.loadPriceFromIEX(symbol)
	if err != nil {
		return 0, err
	}
	h.setPriceInCache(symbol, price)

	return price, nil
}

func (h *Handler) loadPriceFromIEX(symbol string) (int32, error) {
	fmt.Printf("Loading Price From IEX: %v\n", symbol)

	quote, err := h.iex.Quote(symbol)
	if err != nil {
		return 0, err
	}

	return int32(quote.LatestPrice * 100), nil
}

func (h *Handler) loadPriceFromCache(symbol string) (int32, error) {
	quote, ok := h.quotes[symbol]
	if !ok {
		return 0, errNoCacheResult
	}

	if time.Now().Sub(quote.cachedAt).Minutes() > maxCacheMins {
		delete(h.quotes, symbol)
		return 0, errNoCacheResult
	}

	return quote.price, nil
}

func (h *Handler) setPriceInCache(symbol string, price int32) {
	h.quotes[symbol] = quote{price: price, cachedAt: time.Now()}
}
