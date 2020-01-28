package handler

import (
	"context"
	"encoding/json"
	"fmt"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/helpers/microtime"
	proto "github.com/kytra-app/stock-target-price-srv/proto"
	"github.com/kytra-app/stock-target-price-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

// New returns an instance of Handler
func New(iex iex.Service, db storage.Service, client client.Client, broker broker.Broker) *Handler {
	return &Handler{
		db:     db,
		iex:    iex,
		broker: broker,
		stocks: stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	iex    iex.Service
	db     storage.Service
	stocks stocks.StocksService
	broker broker.Broker
}

// List returns all the price targets for a stock
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	t, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	stocks, err := h.db.List(t, req.StockUuids)
	if err != nil {
		return err
	}

	rsp.Stocks = make([]*proto.Stock, len(stocks))
	for i, s := range stocks {
		rsp.Stocks[i] = &proto.Stock{
			Uuid:             s.UUID,
			PriceTarget:      s.PriceTarget,
			NumberOfAnalysts: s.NumberOfAnalysts,
		}
	}

	return nil
}

// Insight is the JSON object published by the insight-srv
type Insight struct {
	StockUUID string `json:"asset_uuid"`
}

// HandleNewInsight processes the messages publlushed by
func (h *Handler) HandleNewInsight(e broker.Event) error {
	fmt.Printf("[HandleNewInsight] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var i Insight
	if err := json.Unmarshal(e.Message().Body, &i); err != nil {
		return err
	}

	// Check for an existng insight
	exists, err := h.db.StockExistsToday(i.StockUUID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if exists {
		fmt.Println("Stock price target already created today")
		return nil
	}

	// Fetch the stock (we need the symbol for IEX query)
	sRsp, err := h.stocks.Get(context.Background(), &stocks.Stock{Uuid: i.StockUUID})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Fetch the price target
	target, err := h.iex.GetPriceTarget(sRsp.GetStock().Symbol)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create the DB price target
	_, err = h.db.Create(storage.Stock{
		UUID:             i.StockUUID,
		PriceTarget:      int64(target.PriceTargetAverage * 100),
		NumberOfAnalysts: int64(target.NumberOfAnalysts),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
