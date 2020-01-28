package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	"github.com/kytra-app/insights-srv/storage"
	valuation "github.com/kytra-app/portfolio-valuation-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	earnings "github.com/kytra-app/stock-earnings-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	cron "github.com/robfig/cron/v3"
)

// New returns an instance of Handler
func New(db storage.Service, client client.Client) Handler {
	return Handler{
		db:         db,
		followers:  followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		users:      users.NewUsersService("kytra-srv-v1-users:8080", client),
		stocks:     stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades:     trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		earnings:   earnings.NewStockEarningsService("kytra-srv-v1-stock-earnings:8080", client),
		valuation:  valuation.NewPortfolioValuationService("kytra-srv-v1-portfolio-valuation:8080", client),
	}
}

// Handler processes incoming messages from the broker
type Handler struct {
	db         storage.Service
	followers  followers.FollowersService
	users      users.UsersService
	stocks     stocks.StocksService
	trades     trades.TradesService
	portfolios portfolios.PortfoliosService
	earnings   earnings.StockEarningsService
	valuation  valuation.PortfolioValuationService

	subNewPost    broker.Subscriber
	subNewMover   broker.Subscriber
	subNewTrade   broker.Subscriber
	subNewArticle broker.Subscriber
	cron          *cron.Cron
}

// Subscribe registeres the consumer to recieve events from the broker
func (h *Handler) Subscribe() (err error) {
	h.subNewPost, err = broker.Subscribe("kytra-srv-v1-posts-post-created", h.HandleNewPost, broker.Queue("insights-post-created"))
	if err != nil {
		return err
	}

	h.subNewMover, err = broker.Subscribe("kytra-srv-v1-stock-movers-mover-created", h.HandleNewMover, broker.Queue("insights-mover-created"))
	if err != nil {
		return err
	}

	h.subNewTrade, err = broker.Subscribe("kytra-srv-v1-trades-trade-created", h.HandleNewTrade, broker.Queue("insights-trade-created"))
	if err != nil {
		return err
	}

	h.subNewArticle, err = broker.Subscribe("kytra-srv-v1-stock-news-article-created", h.HandleNewArticle, broker.Queue("insights-article-created"))
	if err != nil {
		return err
	}

	h.cron = cron.New(cron.WithLocation(time.UTC))
	h.cron.AddFunc("0 7 * * *", h.HandleDailyEarnings)
	h.cron.Start()

	return nil
}

// Unsubscribe deregisters the consumer from recieving events from the broker
func (h *Handler) Unsubscribe() {
	h.cron.Stop()

	h.subNewPost.Unsubscribe()
	h.subNewMover.Unsubscribe()
	h.subNewTrade.Unsubscribe()
	h.subNewArticle.Unsubscribe()
}

func (h *Handler) publishNewInsight(i storage.Insight) {
	bytes, err := json.Marshal(&i)
	if err != nil {
		fmt.Println(err)
		return
	}
	broker.Publish("kytra-srv-v1-insights-insight-created", &broker.Message{Body: bytes})
	fmt.Println("Published Insight")
}
