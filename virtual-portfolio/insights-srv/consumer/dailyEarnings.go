package consumer

import (
	"context"
	"fmt"

	"github.com/kytra-app/insights-srv/storage"
	earnings "github.com/kytra-app/stock-earnings-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
)

// HandleDailyEarnings looks for earnings events happening today. It should be called
// daily by a worker
func (h Handler) HandleDailyEarnings() {
	// Get the earnings
	eRsp, err := h.earnings.List(context.Background(), &earnings.ListRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the stocks and group them by UUID
	stockUUIDs := make([]string, len(eRsp.Earnings))
	for i, e := range eRsp.Earnings {
		stockUUIDs[i] = e.StockUuid
	}
	sRsp, err := h.stocks.List(context.Background(), &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		fmt.Println(err)
		return
	}
	stocksByUUID := make(map[string]*stocks.Stock, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		stocksByUUID[s.Uuid] = s
	}

	// Create an insight for each earning
	for _, e := range eRsp.Earnings {
		stock := stocksByUUID[e.StockUuid]

		i, _ := h.db.CreateInsight(storage.Insight{
			Title:     fmt.Sprintf("%v has an earnings release today", stock.Symbol),
			AssetUUID: e.StockUuid,
			AssetType: "Stock",
			Type:      "EVENT",
		})
		h.publishNewInsight(i)
	}
}
