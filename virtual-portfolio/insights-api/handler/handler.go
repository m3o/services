package handler

import (
	"context"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/helpers/microtime"
	"github.com/kytra-app/helpers/photos"
	"github.com/kytra-app/helpers/unique"
	proto "github.com/kytra-app/insights-api/proto"
	insights "github.com/kytra-app/insights-srv/proto"
	summary "github.com/kytra-app/insights-summary-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth      auth.Authenticator
	iex       iex.Service
	photos    photos.Service
	stocks    stocks.StocksService
	insights  insights.InsightsService
	quotes    quotes.StockQuoteService
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
		insights:  insights.NewInsightsService("kytra-srv-v1-insights:8080", client),
		summary:   summary.NewInsightsSummaryService("kytra-srv-v1-insights-summary:8080", client),
		followers: followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		quotes:    quotes.NewStockQuoteService("kytra-srv-v2-stock-quote:8080", client),
	}
}

// List retrieves the insights for the given day
func (h Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Authenticate the user using the JWT
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHORIZED", "A valid JWT is required")
	}

	// Get the date from the request
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}
	date = date.Truncate(time.Hour * 24)

	// Fetch the individual insights
	iRsp, err := h.insights.ListInsights(ctx, &insights.ListInsightsRequest{UserUuid: user.UUID})
	if err != nil {
		return err
	}

	// Get the stock UUIDs
	stockUUIDs := make([]string, len(iRsp.GetInsights()))
	for i, insight := range iRsp.Insights {
		stockUUIDs[i] = insight.GetAsset().Uuid
	}
	stockUUIDs = unique.Strings(stockUUIDs)

	// Get the quotes
	qRsp, err := h.quotes.ListQuotes(ctx, &quotes.ListRequest{
		Uuids: stockUUIDs, IncludeOutOfHours: true,
	})
	if err != nil {
		return err
	}
	stkQuotes := make(map[string]*quotes.Quote, len(qRsp.GetQuotes()))
	for _, q := range qRsp.GetQuotes() {
		if q.CreatedAt > date.Unix() {
			stkQuotes[q.GetStockUuid()] = q
		}
	}

	// Get the stocks and group by UUID
	sRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return err
	}
	stocksByUUID := make(map[string]*stocks.Stock, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		stocksByUUID[s.Uuid] = s
	}

	// Get the following status' for the stocks
	followingStatuses, _ := h.getFollowingStatuses(ctx, sRsp.GetStocks())

	// Group the insights by stock UUID
	insightsByStockUUID := make(map[string][]*insights.Insight, len(stockUUIDs))
	for _, uuid := range stockUUIDs {
		stockInsights := []*insights.Insight{}

		for _, i := range iRsp.GetInsights() {
			if i.GetAsset().Uuid == uuid {
				stockInsights = append(stockInsights, i)
			}
		}

		insightsByStockUUID[uuid] = stockInsights
	}

	// The the summaries
	if len(stockUUIDs) == 0 {
		return nil
	}
	summaryRsp, err := h.summary.List(ctx, &summary.ListRequest{
		UserUuid:   user.UUID,
		AssetType:  "Stock",
		AssetUuids: stockUUIDs,
	})
	if err != nil {
		return err
	}
	summaryForUUID := make(map[string]string)
	for _, a := range summaryRsp.GetAssets() {
		summaryForUUID[a.Uuid] = a.Summary
	}

	// Serialize the insights
	rsp.Insights = make([]*proto.Insight, len(stockUUIDs))
	for i, uuid := range stockUUIDs {
		quote := stkQuotes[uuid]
		stock := stocksByUUID[uuid]
		insights := insightsByStockUUID[uuid]
		following := followingStatuses[uuid]

		events := make([]*proto.Event, len(insights))
		for i, insight := range insights {
			events[i] = &proto.Event{
				Title:     insight.Title,
				Subtitle:  insight.Subtitle,
				LinkUrl:   insight.LinkUrl,
				PostUuid:  insight.PostUuid,
				Type:      insight.Type,
				Seen:      insight.Seen,
				CreatedAt: time.Unix(insight.CreatedAt, 0).String(),
			}
		}

		summary, ok := summaryForUUID[uuid]
		if !ok {
			summary = ""
		}

		rsp.Insights[i] = &proto.Insight{
			Events:  events,
			Summary: summary,
			Asset: &proto.Asset{
				Type:              "Stock",
				Uuid:              stock.Uuid,
				Name:              stock.Name,
				Sector:            stock.Sector,
				Description:       stock.Description,
				Color:             stock.Color,
				ProfilePictureUrl: h.photos.GetURL(stock.ProfilePictureId),
				Following:         following,
				Symbol:            stock.Symbol,
			},
		}

		if quote != nil {
			rsp.Insights[i].Quote = &proto.Quote{
				Price:            quote.GetPrice(),
				CreatedAt:        quote.GetCreatedAt(),
				MarketClosed:     quote.GetMarketClosed(),
				PercentageChange: quote.GetPercentageChange(),
			}
		}
	}

	return nil
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

// Seen records the user having read the insights for an asset at the current time
func (h Handler) Seen(ctx context.Context, req *proto.SeenRequest, rsp *proto.SeenResponse) error {
	// Authenticate the user using the JWT
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHORIZED", "A valid JWT is required")
	}

	_, err = h.insights.CreateUserView(ctx, &insights.UserView{
		UserUuid: user.UUID,
		Asset: &insights.Asset{
			Uuid: req.StockUuid,
			Type: "Stock",
		},
	})

	return err
}
