package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	// "errors"
	// "fmt"
	// "time"

	iex "github.com/kytra-app/helpers/iex-cloud"
	proto "github.com/kytra-app/stock-movers-srv/proto"
	"github.com/kytra-app/stock-movers-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

const percentageChangeRequired = 5.0

// New returns an instance of Handler
func New(iex iex.Service, db storage.Service, client client.Client, broker broker.Broker) *Handler {
	return &Handler{
		db:     db,
		iex:    iex,
		broker: broker,
		stocks: stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades: trades.NewTradesService("kytra-srv-v1-trades:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	iex    iex.Service
	db     storage.Service
	stocks stocks.StocksService
	trades trades.TradesService
	broker broker.Broker
}

// List returns all the stocks which have moved more than X% in the requested date
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	var date time.Time

	if req.Date != 0 {
		date = time.Unix(req.Date, 0).Truncate(time.Hour * 24)
	} else {
		date = time.Now().Truncate(time.Hour * 24)
	}

	movers, err := h.db.List(date)
	if err != nil {
		return err
	}

	rsp.Movers = make([]*proto.Mover, len(movers))
	for i, m := range movers {
		rsp.Movers[i] = &proto.Mover{
			Date:       m.Date.Unix(),
			StockUuid:  m.StockUUID,
			Percentage: m.Percentage,
		}
	}

	return nil
}

// FetchMovements retrieves and stores all price movements
func (h *Handler) FetchMovements() {
	// Step 1. Find the assets being traded
	tRsp, err := h.trades.AllAssets(context.Background(), &trades.AllAssetsRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 2. Fetch the stocks
	stockUUIDs := make([]string, len(tRsp.GetAssets()))
	for i, a := range tRsp.GetAssets() {
		stockUUIDs[i] = a.Uuid
	}
	sRsp, err := h.stocks.List(context.Background(), &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 3. Fetch the price movements
	movements := make(map[string]float32, len(sRsp.GetStocks()))
	for _, stock := range sRsp.GetStocks() {
		quote, err := h.iex.Quote(stock.Symbol)
		if err != nil {
			continue
		}
		if !quote.USMarketOpen {
			fmt.Println("US Market isn't open")
			return
		}
		movements[stock.Uuid] = quote.ChangePercent * 100
	}

	// Step 4. Create the movements
	for uuid, change := range movements {
		// Step 4.1 Check for min change
		if math.Abs(float64(change)) < percentageChangeRequired {
			continue
		}

		// Step 4.2 Write to the DB
		movement := storage.Movement{
			StockUUID:  uuid,
			Percentage: change,
			Date:       time.Now().Truncate(time.Hour * 24),
		}
		_, err := h.db.Create(movement)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Step 4.3 Publish to broker
		bytes, err := json.Marshal(&movement)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = h.broker.Publish("kytra-srv-v1-stock-movers-mover-created", &broker.Message{Body: bytes})
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
