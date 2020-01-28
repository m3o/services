package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	followers "github.com/kytra-app/followers-srv/proto"
	"github.com/kytra-app/insights-srv/storage"
	valuation "github.com/kytra-app/portfolio-valuation-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/broker"
)

// Trade is the JSON object published by the trades-srv
type Trade struct {
	TypeInt       trades.TradeType `json:"type"`
	PortfolioUUID string           `json:"portfolio_uuid"`
	Quantity      int64            `json:"quantity"`
	UnitPrice     int64            `json:"unit_price"`
	Asset         struct {
		UUID string `json:"uuid"`
		Type string `json:"type"`
	} `json:"asset"`
}

// Type of the trade, BUY or SELL
func (t Trade) Type() string {
	if t.TypeInt == trades.TradeType_BUY {
		return "BUY"
	}
	return "SELL"
}

// HandleNewTrade handles the event when a stock mover is created
func (h *Handler) HandleNewTrade(e broker.Event) error {
	fmt.Printf("[HandleNewTrade] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var trade Trade
	if err := json.Unmarshal(e.Message().Body, &trade); err != nil {
		fmt.Println(err)
		return err
	}

	// Get the stock
	sRsp, err := h.stocks.Get(context.Background(), &stocks.Stock{Uuid: trade.Asset.UUID})
	if err != nil {
		fmt.Println(err)
		return err
	}
	stock := sRsp.GetStock()

	// Get the portfolio
	portfolio, err := h.portfolios.Get(context.Background(), &portfolios.Portfolio{Uuid: trade.PortfolioUUID})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get the users new position
	position, err := h.trades.ListTradesForPosition(context.Background(), &trades.ListRequest{
		IncludeMetadata: true,
		PortfolioUuid:   portfolio.Uuid,
		Asset:           &trades.Asset{Uuid: trade.Asset.UUID, Type: trade.Asset.Type},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get the value of the users portfolio
	value, err := h.valuation.GetPortfolio(context.Background(), &valuation.Portfolio{Uuid: portfolio.Uuid})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Determine how many shares were previously owned
	prevQuantity := position.Quantity
	if trade.Type() == "BUY" {
		prevQuantity = prevQuantity - trade.Quantity
	} else {
		prevQuantity = prevQuantity + trade.Quantity
	}

	// Determine the value of the position before and after
	prevValue := prevQuantity * trade.UnitPrice
	newValue := position.Quantity * trade.UnitPrice

	// Determine the percentage of the portfolio before and after (plus the change)
	prevPercentage := prevValue * 100 / value.TotalValue
	newPercentage := newValue * 100 / value.TotalValue
	changePercentage := math.Abs(float64(newPercentage - prevPercentage))

	// Get the user who placed the trade
	user, err := h.users.Find(context.Background(), &users.User{Uuid: portfolio.UserUuid})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Compose the title
	var title string
	if prevPercentage > 0 && newPercentage > prevPercentage {
		title = fmt.Sprintf(
			"%v added %v%% to their %v position, taking it to %v%%",
			user.FirstName, changePercentage, stock.Name, newPercentage,
		)
	} else if prevPercentage == 0 && newPercentage > prevPercentage {
		title = fmt.Sprintf(
			"%v bought a new %v%% position in %v",
			user.FirstName, newPercentage, stock.Name,
		)
	} else if prevPercentage > newPercentage && newPercentage > 0 {
		title = fmt.Sprintf(
			"%v reduced their %v position by %v%%, taking it to %v%%",
			user.FirstName, stock.Name, changePercentage, newPercentage,
		)
	} else if prevPercentage > newPercentage && newPercentage == 0 {
		title = fmt.Sprintf(
			"%v sold their %v%% %v position",
			user.FirstName, prevPercentage, stock.Name,
		)
	}

	// Find the users who follow the user who made the trade
	userUUIDs, err := h.followersForUser(user.Uuid)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create an insight for each user
	for _, uuid := range userUUIDs {
		i, err := h.db.CreateInsight(storage.Insight{
			Title:     title,
			Type:      "TRADE",
			AssetUUID: stock.Uuid,
			AssetType: "Stock",
			UserUUID:  uuid,
		})

		if err != nil {
			fmt.Println(err)
			continue
		}

		h.publishNewInsight(i)
	}

	return nil
}

func (h Handler) followersForUser(userUUID string) ([]string, error) {
	query := followers.Resource{Type: "User", Uuid: userUUID}

	rsp, err := h.followers.Get(context.Background(), &query)
	if err != nil {
		return []string{}, err
	}

	uuids := make([]string, len(rsp.Followers))
	for i, f := range rsp.Followers {
		uuids[i] = f.Uuid
	}

	return uuids, nil
}
