package handler

import (
	"context"
	"fmt"
	"math"
	"time"

	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/microtime"
	proto "github.com/kytra-app/home-cards-api/proto"
	valuation "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	quote "github.com/kytra-app/stock-quote-srv-v2/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth       auth.Authenticator
	quote      quote.StockQuoteService
	trades     trades.TradesService
	valuation  valuation.PortfolioValueTrackingService
	portfolios portfolios.PortfoliosService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, client client.Client) Handler {
	return Handler{
		auth:       auth,
		quote:      quote.NewStockQuoteService("kytra-srv-v2-stock-quote:8080", client),
		trades:     trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		valuation:  valuation.NewPortfolioValueTrackingService("kytra-srv-v1-portfolio-value-tracking:8080", client),
	}
}

// List retrieves the cards for today
func (h Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Authenticate the user using the JWT
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return errors.Unauthorized("UNAUTHORIZED", "A valid JWT is required")
	}

	// Get the users portfolio
	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.UUID})
	if err != nil {
		return err
	}

	// Check to see if the user has active positions
	tRsp, err := h.trades.ListPositionsForPortfolio(ctx, &trades.ListRequest{
		PortfolioUuid: portfolio.Uuid,
	})
	if err != nil {
		return err
	}
	hasPositions := len(tRsp.GetPositions()) > 0

	// If the user has no position, insert the "Getting Started" card
	if !hasPositions {
		rsp.Cards = append(rsp.GetCards(), h.gettingStartedCard())
	}

	// If in the evening and the user has some positions, insert the "Evening Summary"
	// eveningSummaryTime := time.Now().Truncate(time.Hour * 24).Add(time.Hour * 21).Add(time.Minute * 30)
	if hasPositions { //&& eveningSummaryTime.Unix() <= time.Now().Unix() {
		card, err := h.eveningSummaryCard(ctx, user.UUID)
		if err == nil {
			rsp.Cards = append(rsp.GetCards(), card)
		} else {
			fmt.Println(err)
		}
	}

	// If the user has some positions, insert the "Morning Summary"
	if hasPositions {
		card, err := h.morningSummaryCard(ctx, user.UUID)
		if err == nil {
			rsp.Cards = append(rsp.GetCards(), card)
		} else {
			fmt.Println(err)
		}
	}

	return nil
}

func (h *Handler) gettingStartedCard() *proto.Card {
	return &proto.Card{
		Type:     "GETTING_STARTED",
		Title:    "Get Started",
		Subtitle: "Welcome to Kytra 👋 Tap here for notes on how to get started.",
	}
}

func (h *Handler) morningSummaryCard(ctx context.Context, userUUID string) (*proto.Card, error) {
	prevDay := time.Now().Add(time.Hour * -24)
	if prevDay.Weekday().String() == "Sunday" {
		prevDay = prevDay.Add(time.Hour * -48)
	}

	summary, err := h.summaryForDate(ctx, prevDay)
	if err != nil {
		return nil, err
	}

	return &proto.Card{
		Type:     "MORNING_SUMMARY",
		Title:    "Morning Summary",
		Subtitle: summary,
	}, nil
}

func (h *Handler) eveningSummaryCard(ctx context.Context, userUUID string) (*proto.Card, error) {
	summary, err := h.summaryForDate(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	return &proto.Card{
		Type:     "EVENING_SUMMARY",
		Title:    "Evening Summary",
		Subtitle: summary,
	}, nil
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
		if diff > 0.5 {
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

	quote, err := h.quote.GetQuote(ctx, &quote.Stock{Uuid: "^IXIC"})
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
