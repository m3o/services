package handler

import (
	"context"
	"math"
	"time"

	followers "github.com/kytra-app/followers-srv/proto"
	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/microtime"
	"github.com/kytra-app/helpers/photos"
	valuation "github.com/kytra-app/portfolio-valuation-srv/proto"
	valueTracking "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	proto "github.com/kytra-app/portfolios-api/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	quote "github.com/kytra-app/stock-quote-srv-v2/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth          auth.Authenticator
	pics          photos.Service
	stocks        stocks.StocksService
	trades        trades.TradesService
	portfolios    portfolios.PortfoliosService
	valuation     valuation.PortfolioValuationService
	followers     followers.FollowersService
	quote         quote.StockQuoteService
	valueTracking valueTracking.PortfolioValueTrackingService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, pics photos.Service, client client.Client) Handler {
	return Handler{
		auth:          auth,
		pics:          pics,
		quote:         quote.NewStockQuoteService("kytra-srv-v2-stock-quote:8080", client),
		stocks:        stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		trades:        trades.NewTradesService("kytra-srv-v1-trades:8080", client),
		followers:     followers.NewFollowersService("kytra-srv-v1-followers:8080", client),
		portfolios:    portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
		valuation:     valuation.NewPortfolioValuationService("kytra-srv-v1-portfolio-valuation:8080", client),
		valueTracking: valueTracking.NewPortfolioValueTrackingService("kytra-srv-v1-portfolio-value-tracking:8080", client),
	}
}

// SetTargets sets the asset class and industry targets for a users portfolio
func (h Handler) SetTargets(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	portfolioUUID, err := h.portfolioUUIDForUser(ctx)
	if err != nil {
		return err
	}

	params := &portfolios.Portfolio{
		Uuid:                                portfolioUUID,
		AssetClassTargetStocks:              req.AssetClassTargetStocks,
		AssetClassTargetCash:                req.AssetClassTargetCash,
		IndustryTargetInformationTechnology: req.IndustryTargetInformationTechnology,
		IndustryTargetFinancials:            req.IndustryTargetFinancials,
		IndustryTargetEnergy:                req.IndustryTargetEnergy,
		IndustryTargetHealthCare:            req.IndustryTargetHealthCare,
		IndustryTargetMaterials:             req.IndustryTargetMaterials,
		IndustryTargetUtilities:             req.IndustryTargetUtilities,
		IndustryTargetRealEstate:            req.IndustryTargetRealEstate,
		IndustryTargetConsumerDiscretionary: req.IndustryTargetConsumerDiscretionary,
		IndustryTargetConsumerStaples:       req.IndustryTargetConsumerStaples,
		IndustryTargetCommunicationServices: req.IndustryTargetCommunicationServices,
		IndustryTargetIndustrials:           req.IndustryTargetIndustrials,
	}

	_, err = h.portfolios.Update(ctx, params)
	return err
}

// GetPortfolio retrieves the value and positions for the current user
func (h Handler) GetPortfolio(ctx context.Context, req *proto.Portfolio, rsp *proto.Portfolio) error {
	portfolioUUID, err := h.portfolioUUIDForUser(ctx)
	if err != nil {
		return err
	}

	return h.serializePortfolio(ctx, portfolioUUID, rsp)
}

// GetInvestor retrieves the value and positions for the given investor
func (h Handler) GetInvestor(ctx context.Context, req *proto.Investor, rsp *proto.Portfolio) error {
	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: req.Uuid})
	if err != nil {
		return err
	}

	return h.serializePortfolio(ctx, portfolio.Uuid, rsp)
}

func (h Handler) serializePortfolio(ctx context.Context, portfolioUUID string, rsp *proto.Portfolio) error {
	// Get the current time
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	// Get the positions
	positions, err := h.positionsForPortfolio(ctx, portfolioUUID)
	if err != nil {
		return err
	}

	// Get the portfolios value
	valuation, err := h.valuation.GetPortfolio(ctx, &valuation.Portfolio{Uuid: portfolioUUID})
	if err != nil {
		return err
	}
	for i, position := range positions {
		positions[i].PercentageOfPortfolio = toPercentage(position.Value, valuation.AssetsValue)
	}

	// Get the one day price movement
	vRsp, err := h.valueTracking.GetPriceMovement(ctx, &valueTracking.GetPriceMovementsRequest{
		PortfolioUuid: portfolioUUID, StartDate: date.Unix(), EndDate: date.Unix(),
	})
	if err != nil {
		return err
	}

	// TODO: Refactor once we move away from simulated portfolios (with a fixed start balance of $100k)
	const startingBalance = 100000 * 100
	lifetimeGain := valuation.TotalValue - startingBalance
	lifetimeGainPercentage := toPercentage(lifetimeGain, startingBalance)

	// Return the result
	*rsp = proto.Portfolio{
		Uuid:                   portfolioUUID,
		Positions:              positions,
		TotalValue:             valuation.TotalValue,
		CashValue:              valuation.CashValue,
		AssetsValue:            valuation.AssetsValue,
		LifetimeGain:           lifetimeGain,
		LifetimeGainPercentage: float32(lifetimeGainPercentage),
	}

	if mov := vRsp.GetPriceMovement(); mov != nil {
		rsp.OneDayGain = mov.GetLatestValue() - mov.GetEarliestValue()
		rsp.OneDayGainPercentage = mov.GetPercentageChange()
	}

	rsp.AssetClasses = h.groupPositionsIntoAssetClasses(positions, rsp)
	rsp.AssetIndustries = h.groupPositionsIntoAssetIndustries(positions, rsp)

	return nil
}

func (h Handler) portfolioUUIDForUser(ctx context.Context) (string, error) {
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return "", errors.Unauthorized("UNAUTHENTICATED", "Invalid JWT")
	}

	porfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.UUID})
	if err != nil {
		return "", err
	}

	return porfolio.Uuid, nil
}

func (h Handler) positionsForPortfolio(ctx context.Context, portfolioUUID string) ([]*proto.Position, error) {
	// Lookup the trades
	positionsRsp, err := h.trades.ListPositionsForPortfolio(ctx, &trades.ListRequest{PortfolioUuid: portfolioUUID})
	if err != nil {
		return []*proto.Position{}, err
	}

	// Find the following status
	stockUUIDs := make([]string, len(positionsRsp.Positions))
	for i, p := range positionsRsp.Positions {
		stockUUIDs[i] = p.Asset.Uuid
	}
	followings, _ := h.fetchFollowingsForStocks(ctx, stockUUIDs)

	// Fetch the quotes
	qRsp, err := h.quote.ListQuotes(ctx, &quote.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return []*proto.Position{}, nil
	}
	quotesForUUID := make(map[string]*quote.Quote, len(qRsp.Quotes))
	for _, q := range qRsp.Quotes {
		quotesForUUID[q.StockUuid] = q
	}

	// Serialize the data
	positions := []*proto.Position{}
	for _, position := range positionsRsp.Positions {
		quote, ok := quotesForUUID[position.Asset.Uuid]
		if !ok {
			continue
		}

		val, err := h.serializePosition(ctx, position, quote)
		if err != nil {
			return positions, err
		}
		val.AssetFollowing = followings[val.AssetUuid]
		positions = append(positions, val)
	}

	return positions, nil
}

func (h Handler) serializePosition(ctx context.Context, position *trades.Position, quote *quote.Quote) (*proto.Position, error) {
	// Construct the result
	result := &proto.Position{
		AssetType: position.Asset.Type,
		AssetUuid: position.Asset.Uuid,
		Quantity:  position.Quantity,
		BookCost:  position.BookCost,
	}

	// Fetch the stock
	stockRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: position.Asset.Uuid})
	if err != nil {
		return result, err
	}
	result.AssetName = stockRsp.Stock.Name
	result.AssetColor = stockRsp.Stock.Color
	result.AssetSector = stockRsp.Stock.Sector
	result.AssetDescription = stockRsp.Stock.Description
	result.AssetProfilePictureUrl = h.pics.GetURL(stockRsp.Stock.ProfilePictureId)

	// Find the current price (non-critical)
	result.UnitPrice = quote.Price
	result.Value = quote.Price * position.Quantity
	result.GainLoss = result.Value - result.BookCost
	result.GainLossPercentage = toPercentage(result.GainLoss, result.BookCost)

	if quote.CreatedAt > time.Now().Truncate(time.Hour*24).Unix() {
		result.OneDayChangePercentage = quote.PercentageChange
	} else {
		result.PrevDayChangePercentage = quote.PercentageChange
	}

	return result, nil
}

func toPercentage(top, bottom int64) float32 {
	if bottom == 0 {
		return 0
	}

	value := float64(top) / float64(bottom)
	return float32(math.Round(10000*value)) / 100
}

func (h *Handler) groupPositionsIntoAssetClasses(positions []*proto.Position, rsp *proto.Portfolio) []*proto.Category {
	aggregate := struct {
		Stocks int64
		Bonds  int64
	}{}

	for _, p := range positions {
		switch p.AssetType {
		case "Stock":
			aggregate.Stocks += p.Value
		case "Bond":
			aggregate.Bonds += p.Value
		}
	}

	portfolio, err := h.portfolios.Get(context.Background(), &portfolios.Portfolio{Uuid: rsp.Uuid})
	if err != nil {
		return nil
	}

	return []*proto.Category{
		&proto.Category{
			Name:              "Stocks",
			TargetPercentage:  portfolio.AssetClassTargetStocks,
			CurrentPercentage: toPercentage(aggregate.Stocks, rsp.TotalValue),
		},
		&proto.Category{
			Name:              "Cash",
			TargetPercentage:  portfolio.AssetClassTargetCash,
			CurrentPercentage: toPercentage(rsp.TotalValue-aggregate.Stocks, rsp.TotalValue),
		},
	}
}

func (h *Handler) groupPositionsIntoAssetIndustries(positions []*proto.Position, rsp *proto.Portfolio) (res []*proto.Category) {
	aggregate := map[string]int64{
		"Information Technology": 0,
		"Financials":             0,
		"Energy":                 0,
		"Health Care":            0,
		"Materials":              0,
		"Utilities":              0,
		"Real Estate":            0,
		"Consumer Discretionary": 0,
		"Consumer Staples":       0,
		"Communication Services": 0,
		"Industrials":            0,
	}

	var total int64
	for _, p := range positions {
		total += p.Value
		aggregate[p.AssetSector] += p.Value
	}

	portfolio, err := h.portfolios.Get(context.Background(), &portfolios.Portfolio{Uuid: rsp.Uuid})
	if err != nil {
		return nil
	}

	for name, value := range aggregate {
		var target float32
		switch name {
		case "Information Technology":
			target = portfolio.IndustryTargetInformationTechnology
		case "Financials":
			target = portfolio.IndustryTargetFinancials
		case "Energy":
			target = portfolio.IndustryTargetEnergy
		case "Health Care":
			target = portfolio.IndustryTargetHealthCare
		case "Materials":
			target = portfolio.IndustryTargetMaterials
		case "Utilities":
			target = portfolio.IndustryTargetUtilities
		case "Real Estate":
			target = portfolio.IndustryTargetRealEstate
		case "Consumer Discretionary":
			target = portfolio.IndustryTargetConsumerDiscretionary
		case "Consumer Staples":
			target = portfolio.IndustryTargetConsumerStaples
		case "Communication Services":
			target = portfolio.IndustryTargetCommunicationServices
		case "Industrials":
			target = portfolio.IndustryTargetIndustrials
		}

		res = append(res, &proto.Category{
			Name:              name,
			TargetPercentage:  target,
			CurrentPercentage: toPercentage(value, total),
		})
	}

	return res
}

func (h Handler) fetchFollowingsForStocks(ctx context.Context, uuids []string) (map[string]bool, error) {
	// Setup response, default to not following users
	rsp := make(map[string]bool, len(uuids))
	for _, uuid := range uuids {
		rsp[uuid] = false
	}

	// Try and get the user from the context
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return rsp, nil
	}

	// Construct the query
	req := &followers.ListRequest{
		Follower:      &followers.Resource{Uuid: u.UUID, Type: "User"},
		FolloweeType:  "Stock",
		FolloweeUuids: uuids,
	}

	// Request the data
	data, err := h.followers.List(ctx, req)
	if err != nil {
		return rsp, err
	}

	// Update the response
	for _, r := range data.Resources {
		rsp[r.Uuid] = r.Following
	}

	return rsp, nil
}
