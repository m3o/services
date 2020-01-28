package handler

import (
	"context"
	"fmt"
	"math"
	"time"

	feed "github.com/kytra-app/feed-items-srv/proto"
	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/microtime"
	valuation "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
	earnings "github.com/kytra-app/stock-earnings-srv/proto"
	news "github.com/kytra-app/stock-news-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	target "github.com/kytra-app/stock-target-price-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(auth auth.Authenticator, client client.Client) *Handler {
	return &Handler{
		auth:          auth,
		news:          news.NewStockNewsService("kytra-srv-v1-stock-news:8080", client),
		feed:          feed.NewFeedItemsService("kytra-srv-v1-feed-items:8080", client),
		posts:         posts.NewPostsService("kytra-srv-v1-posts:8080", client),
		users:         users.NewUsersService("kytra-srv-v1-users:8080", client),
		stocks:        stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades:        trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		quotes:        quotes.NewStockQuoteService("kytra-srv-v2-stock-quote:8080", client),
		targets:       target.NewStockTargetPriceService("kytra-srv-v1-stock-target-price:8080", client),
		earnings:      earnings.NewStockEarningsService("kytra-srv-v1-stock-earnings:8080", client),
		valuation:     valuation.NewPortfolioValueTrackingService("kytra-srv-v1-portfolio-value-tracking:8080", client),
		followers:     followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		portfolios:    portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		luckyDipCache: map[int64]string{},
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	auth       auth.Authenticator
	news       news.StockNewsService
	feed       feed.FeedItemsService
	users      users.UsersService
	posts      posts.PostsService
	stocks     stocks.StocksService
	trades     trades.TradesService
	quotes     quotes.StockQuoteService
	targets    target.StockTargetPriceService
	earnings   earnings.StockEarningsService
	followers  followers.FollowersService
	valuation  valuation.PortfolioValueTrackingService
	portfolios portfolios.PortfoliosService

	luckyDipCache map[int64]string
}

func (h *Handler) summaryForDate(ctx context.Context, t time.Time) (string, error) {
	portfolio, err := h.fetchPortfolioPrice(ctx, t)
	if err != nil {
		return "", err
	}

	index, err := h.fetchIndexPrice(ctx, t)
	if err != nil {
		return "", err
	}

	diff := portfolio - index
	var format string
	if index > 0 && diff > 0 {
		format = "The markets had a good day %v, rising %v. In addition to this, your portfolio outperformed the market by %v 🎉"
	} else if index > 0 && diff < 0 {
		emoji := "😬"
		if diff < 0.5 {
			emoji = "😭"
		}
		format = "The markets had a bullish day %v, gaining %v. Your portfolio lagged behind the market by %v " + emoji
	} else if index < 0 && diff > 0 {
		format = "The markets took a bit of a beating %v, falling %v. However, your portfolio weathered the storm well and beat the market by %v 😄"
	} else {
		format = "The markets had a rough day %v, falling %v. Your portfolio was hit especially hard and fell by an additional %v 😬"
	}

	date := t.Truncate(time.Hour * 24)
	today := time.Now().Truncate(time.Hour * 24)
	yesterday := today.Add(time.Hour * -24)

	var dayStr string
	if date.Unix() >= today.Unix() {
		dayStr = "today"
	} else if date.Unix() >= yesterday.Unix() {
		dayStr = "yesterday"
	} else {
		dayStr = t.Weekday().String()
	}

	formatFloat32 := func(val float32) string {
		abs := math.Abs(float64(val))
		rounded := math.Round(abs*100) / 100

		if val > 0.1 {
			return fmt.Sprintf("%v%%", rounded)
		}
		return "less than 0.1%"
	}

	return fmt.Sprintf(format, dayStr, formatFloat32(index), formatFloat32(diff)), nil
}

// fetchIndexPrice gets the day % chg for the Nasdaq Composite Index
func (h *Handler) fetchIndexPrice(ctx context.Context, t time.Time) (float32, error) {
	t = t.Truncate(time.Hour * 24)

	ctx = microtime.ContextWithTime(ctx, t.Add(time.Hour*24))
	formattedErr := errors.InternalServerError("MISSING_INDEX_DATA", "No index data was found")

	quote, err := h.quotes.GetQuote(ctx, &quotes.Stock{Uuid: "^IXIC"})
	if err != nil || quote == nil {
		return 0, formattedErr
	}
	// Quote wasn't for the day requested (it most likely doesn't exist yet)
	ca := time.Unix(quote.CreatedAt, 0).Truncate(time.Hour * 24)
	if ca != t {
		return 0, formattedErr
	}

	return quote.GetPercentageChange(), nil
}

// fetchPortfolioPrice gets the % chg for the users portfolios
func (h *Handler) fetchPortfolioPrice(ctx context.Context, t time.Time) (float32, error) {
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return 0, err
	}

	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.UUID})
	if err != nil {
		return 0, err
	}

	rsp, err := h.valuation.GetPriceMovement(ctx, &valuation.GetPriceMovementsRequest{
		PortfolioUuid: portfolio.Uuid, StartDate: t.Unix(), EndDate: t.Unix(),
	})
	if err != nil {
		return 0, err
	}

	return rsp.GetPriceMovement().GetPercentageChange(), nil
}

func (h *Handler) followingForResource(ctx context.Context, user *users.User, rType string) ([]string, error) {
	// Step 1. Get the resources the users are following
	fRsp, err := h.followers.Get(ctx, &followers.Resource{Uuid: user.Uuid, Type: "User"})
	if err != nil {
		return nil, err
	}

	// Step 2. Filter down to the requested type
	uuids := []string{}
	for _, r := range fRsp.GetFollowing() {
		if r.Type == rType {
			uuids = append(uuids, r.Uuid)
		}
	}

	return uuids, nil
}
