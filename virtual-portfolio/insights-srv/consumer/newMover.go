package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/kytra-app/insights-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/broker"
)

// Mover is the JSON object published by the stock-movers-srv
type Mover struct {
	StockUUID  string  `json:"stock_uuid"`
	Percentage float32 `json:"percentage"`
}

// HandleNewMover handles the event when a stock mover is created
func (h *Handler) HandleNewMover(e broker.Event) error {
	fmt.Printf("[HandleNewMover] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var mover Mover
	if err := json.Unmarshal(e.Message().Body, &mover); err != nil {
		return err
	}

	// Find the stock
	sRsp, err := h.stocks.Get(context.Background(), &stocks.Stock{Uuid: mover.StockUUID})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	stock := sRsp.Stock

	// Create a post in the feed
	direction := "down"
	if mover.Percentage > 0 {
		direction = "up"
	}

	percentage := math.Round(float64(mover.Percentage)*100) / 100

	i, err := h.db.CreateInsight(storage.Insight{
		Title:     fmt.Sprintf("%v is %v %v%% since markets opened", stock.Symbol, direction, percentage),
		Type:      "PRICE_MOVEMENT",
		AssetUUID: mover.StockUUID,
		AssetType: "Stock",
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	h.publishNewInsight(i)

	return nil
}
