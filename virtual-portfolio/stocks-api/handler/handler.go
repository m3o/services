package handler

import (
	"context"
	"math/rand"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/helpers/photos"
	summary "github.com/kytra-app/insights-summary-srv/proto"
	proto "github.com/kytra-app/stocks-api/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/client"
)

var popularStockUUIDs = []string{"SIRI", "AAPL", "AMD", "CSCO", "INTC", "LYFT", "MSFT", "MU", "NVDA", "ZNGA"}

// Handler is an object can process RPC requests
type Handler struct {
	auth      auth.Authenticator
	iex       iex.Service
	photos    photos.Service
	stocks    stocks.StocksService
	summary   summary.InsightsSummaryService
	followers followers.FollowersService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, iex iex.Service, pics photos.Service, client client.Client) Handler {
	return Handler{
		auth:      auth,
		iex:       iex,
		photos:    pics,
		stocks:    stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		summary:   summary.NewInsightsSummaryService("kytra-srv-v1-insights-summary:8080", client),
		followers: followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
	}
}

// Search returns the stocks which match the query provided
func (h Handler) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.ListResponse) error {
	query := &stocks.SearchRequest{Query: req.Query, Limit: req.Limit}
	stocksRsp, err := h.stocks.Search(ctx, query)
	if err != nil {
		return err
	}

	rsp.Stocks, err = h.serializeStocks(ctx, stocksRsp.Stocks)
	return err
}

// Popular returns the 10 most popular stocks
func (h Handler) Popular(ctx context.Context, req *proto.PopularRequest, rsp *proto.ListResponse) error {
	var stocksRsp *stocks.ListResponse
	var err error

	if req.Sector != "" {
		stocksRsp, err = h.stocks.List(ctx, &stocks.ListRequest{Sector: req.Sector})
	} else {
		stocksRsp, err = h.stocks.List(ctx, &stocks.ListRequest{Symbols: popularStockUUIDs})
	}

	if err != nil {
		return err
	}

	// Shuffle Data
	stocks := stocksRsp.Stocks
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(stocks), func(i, j int) { stocks[i], stocks[j] = stocks[j], stocks[i] })

	rsp.Stocks, err = h.serializeStocks(ctx, stocksRsp.Stocks)
	return err
}

// Get retrieves the stock and recent posts
func (h Handler) Get(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	sRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: req.Uuid})
	if err != nil {
		return err
	}

	// Is this user following this stock
	following, _ := h.getFollowingStatus(ctx, sRsp.Stock)

	// Get the summary (non-critical)
	var stockSummary string
	if user, err := h.auth.UserFromContext(ctx); err == nil {
		sCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
		defer cancel()

		sRsp, err := h.summary.Get(sCtx, &summary.GetRequest{
			UserUuid: user.UUID, AssetType: "Stock", AssetUuid: req.Uuid,
		})

		if err == nil {
			stockSummary = sRsp.GetSummary()
		}
	}

	rsp.Stock = &proto.Stock{
		Uuid:              sRsp.Stock.Uuid,
		Name:              sRsp.Stock.Name,
		Symbol:            sRsp.Stock.Symbol,
		Exchange:          sRsp.Stock.Exchange,
		Type:              sRsp.Stock.Type,
		Region:            sRsp.Stock.Region,
		Currency:          sRsp.Stock.Currency,
		Sector:            sRsp.Stock.Sector,
		Industry:          sRsp.Stock.Industry,
		Website:           sRsp.Stock.Website,
		Description:       sRsp.Stock.Description,
		Color:             sRsp.Stock.Color,
		ProfilePictureUrl: h.photos.GetURL(sRsp.Stock.ProfilePictureId),
		Following:         following,
		Summary:           stockSummary,
	}

	// Get the basic market info
	if s, err := h.getMarketInfo(rsp.Stock); err != nil {
		rsp.Stock = s
	}

	return nil
}

// Follow creates a follower relationship between the user and the stock
func (h Handler) Follow(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		rsp.Error = &proto.Error{Code: 401, Message: err.Error()}
		return nil
	}

	followReq := followers.Request{
		Followee: &followers.Resource{Uuid: req.Uuid, Type: "Stock"},
		Follower: &followers.Resource{Uuid: u.UUID, Type: "User"},
	}

	_, err = h.followers.Follow(ctx, &followReq)
	return err
}

// Unfollow deletes the follower relationship between the user and the stock
func (h Handler) Unfollow(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	// Authenticate the user using the JWT
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		rsp.Error = &proto.Error{Code: 401, Message: err.Error()}
		return nil
	}

	unFollowReq := followers.Request{
		Followee: &followers.Resource{Uuid: req.Uuid, Type: "Stock"},
		Follower: &followers.Resource{Uuid: u.UUID, Type: "User"},
	}

	_, err = h.followers.Unfollow(ctx, &unFollowReq)
	return err
}

func (h Handler) getFollowingStatus(ctx context.Context, followee *stocks.Stock) (bool, error) {
	rsp, err := h.getFollowingStatuses(ctx, []*stocks.Stock{followee})

	if following, ok := rsp[followee.Uuid]; ok {
		return following, nil
	}

	return false, err
}

func (h Handler) getFollowingStatuses(ctx context.Context, followees []*stocks.Stock) (map[string]bool, error) {
	follower, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return map[string]bool{}, nil
	}

	fContext, fCancel := context.WithTimeout(ctx, 150*time.Millisecond)
	defer fCancel()

	uuids := make([]string, len(followees))
	for i, user := range followees {
		uuids[i] = user.Uuid
	}

	fRsp, err := h.followers.List(fContext, &followers.ListRequest{
		Follower:      &followers.Resource{Type: "User", Uuid: follower.UUID},
		FolloweeType:  "Stock",
		FolloweeUuids: uuids,
	})

	if err != nil {
		return map[string]bool{}, err
	}

	rsp := make(map[string]bool, len(fRsp.Resources))
	for _, r := range fRsp.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}

func (h Handler) serializeStocks(ctx context.Context, stocks []*stocks.Stock) ([]*proto.Stock, error) {
	followings, err := h.getFollowingStatuses(ctx, stocks)
	if err != nil {
		return []*proto.Stock{}, err
	}

	rsp := make([]*proto.Stock, len(stocks))
	for i, stock := range stocks {
		rsp[i] = &proto.Stock{
			Uuid:              stock.Uuid,
			Name:              stock.Name,
			Symbol:            stock.Symbol,
			Exchange:          stock.Exchange,
			Sector:            stock.Sector,
			Type:              stock.Type,
			Region:            stock.Region,
			Currency:          stock.Currency,
			Color:             stock.Color,
			Description:       stock.Description,
			ProfilePictureUrl: h.photos.GetURL(stock.ProfilePictureId),
			Following:         followings[stock.Uuid],
		}
	}

	return rsp, nil
}

func (h Handler) getMarketInfo(stock *proto.Stock) (*proto.Stock, error) {
	stats, err := h.iex.KeyStats(stock.Symbol)
	if err != nil {
		return stock, err
	}

	stock.MarketCap = stats.MarketCap
	stock.Week_52High = stats.Week52High
	stock.Week_52Low = stats.Week52Low
	stock.Avg_10Volume = stats.Avg10Volume
	stock.TtmEps = stats.TTMEPS
	stock.TtmDividendRate = stats.TTMDividendRate
	stock.DividendYield = stats.DividendYield
	stock.PeRatio = stats.PERatio
	stock.Beta = stats.Beta
	stock.YtdChangePercent = stats.YtdChangePercent
	stock.Month_1ChangePercent = stats.Month1ChangePercent
	stock.Day_5ChangePercent = stats.Day5ChangePercent
	stock.ExDividendDate = stats.ExDividendDate
	stock.NextEarningsDate = stats.NextEarningsDate
	stock.NextDividendDate = stats.NextDividendDate

	prevDay, err := h.iex.PreviousDayPrice(stock.Symbol)
	if err != nil {
		return stock, err
	}

	stock.PrevDayOpen = prevDay.Open
	stock.PrevDayClose = prevDay.Close
	stock.PrevDayHigh = prevDay.High
	stock.PrevDayLow = prevDay.Low
	stock.PrevDayVolume = prevDay.Volume

	quote, err := h.iex.Quote(stock.Symbol)
	if err != nil {
		return stock, err
	}
	stock.CurrentPrice = quote.LatestPrice

	return stock, nil
}
