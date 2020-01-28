package handler

import (
	"context"
	"fmt"
	"time"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/helpers/microtime"
	proto "github.com/kytra-app/stock-earnings-srv/proto"
	"github.com/kytra-app/stock-earnings-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
)

const percentageChangeRequired = 2.0

// New returns an instance of Handler
func New(iex iex.Service, db storage.Service, client client.Client) *Handler {
	return &Handler{
		db:     db,
		iex:    iex,
		stocks: stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	iex    iex.Service
	broker broker.Broker
	db     storage.Service
	stocks stocks.StocksService
}

// List returns all the earnings occuring on a given date
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}
	date = date.Truncate(time.Hour * 24)

	earnings, err := h.db.List(date, req.StockUuids)
	if err != nil {
		return err
	}

	rsp.Earnings = make([]*proto.Earning, len(earnings))
	for i, e := range earnings {
		rsp.Earnings[i] = &proto.Earning{
			Date:      e.Date.Unix(),
			StockUuid: e.StockUUID,
		}
	}

	return nil
}

// FetchEarnings retrieves and stores all earnings
func (h *Handler) FetchEarnings() {
	// Step 1. Fetch the earnings
	earnings, err := h.iex.ListUpcomingEarnings()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 2. Fetch the stocks
	symbols := make([]string, len(earnings))
	for i, e := range earnings {
		symbols[i] = e.Symbol
	}
	sRsp, err := h.stocks.List(context.Background(), &stocks.ListRequest{Symbols: symbols})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 3. Get the UUIDs for the symbols
	uuidBySymbol := make(map[string]string, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		uuidBySymbol[s.Symbol] = s.Uuid
	}

	// Step 4. Create the earnings
	for _, e := range earnings {
		// Step 4.1 Get the UUID
		uuid, ok := uuidBySymbol[e.Symbol]
		if !ok {
			fmt.Printf("No stock found for %v\n", e.Symbol)
			continue
		}

		// Step 4.2 Write to the DB
		date, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if _, err := h.db.Create(storage.Earning{StockUUID: uuid, Date: date}); err != nil {
			fmt.Println(err)
			continue
		}
	}
}
