package handler

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	proto "github.com/kytra-app/daily-summary-api/proto"
	feed "github.com/kytra-app/feed-items-srv/proto"
	reactlink "github.com/kytra-app/helpers/reactlink"
	"github.com/kytra-app/helpers/unique"
	valuation "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
)

// Evening returns the evening summary for today
func (h *Handler) Evening(ctx context.Context, req *proto.Request, rsp *proto.EveningSummary) (err error) {
	// Step 1. Get the user
	u, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return err
	}
	user, err := h.users.Find(ctx, &users.User{Uuid: u.UUID})
	if err != nil {
		return err
	}

	// Step 2. Attach the introduction
	if intro, err := h.eveningIntroduction(ctx, user); err == nil {
		rsp.Introduction = intro
	} else {
		fmt.Println(err)
	}

	// Step 3. Attach the trades
	if trades, err := h.eveningTrades(ctx, user); err == nil {
		rsp.Trades = trades
	} else {
		fmt.Println(err)
	}

	// Step 4. Attach the performance
	if performance, err := h.eveningPerformance(ctx, user); err == nil {
		rsp.Performance = performance
	} else {
		fmt.Println(err)
	}

	// Step 5. Attach the posts
	if posts, err := h.eveningPosts(ctx, user); err == nil {
		rsp.Posts = posts
	} else {
		fmt.Println(err)
	}

	// Step 6. Attach the benchmarking
	if benchmarking, err := h.eveningBenchmarking(ctx, user); err == nil {
		rsp.Benchmarking = benchmarking
	} else {
		fmt.Println(err)
	}

	return nil
}

func (h *Handler) eveningIntroduction(ctx context.Context, user *users.User) (*proto.Section, error) {
	// Step 1. Write the summary (e.g. "The maket rose 1.3%...")
	summary, err := h.summaryForDate(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	// Step 2. Return the result
	body := fmt.Sprintf("Good evening %v,\n\n%v", user.FirstName, summary)
	return &proto.Section{Body: body}, nil
}

var noPerformanceSection = &proto.Section{
	Title: "Performance",
	Body:  "You have no active positions.",
}

type mover struct {
	stockUUID     string
	change        float32
	price         int64
	chgPercentage float32
}

func (h Handler) performanceAttribution(ctx context.Context, portfolio *portfolios.Portfolio) ([]mover, error) {
	// Step 1. Get the positions for the portfolio
	tRsp, err := h.trades.ListPositionsForPortfolio(ctx, &trades.ListRequest{PortfolioUuid: portfolio.Uuid})
	if err != nil {
		return nil, err
	}
	if len(tRsp.GetPositions()) == 0 {
		return []mover{}, nil
	}

	// Step 2. Get the quotes and 1 day changes
	stockUUIDs := make([]string, len(tRsp.GetPositions()))
	for i, p := range tRsp.GetPositions() {
		stockUUIDs[i] = p.GetAsset().GetUuid()
	}
	qRsp, err := h.quotes.ListQuotes(ctx, &quotes.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return nil, err
	}
	oneDayChanges := make(map[string]float32)
	prices := make(map[string]int64)
	for _, q := range qRsp.GetQuotes() {
		if q.CreatedAt > time.Now().Truncate(time.Hour*24).Unix() {
			oneDayChanges[q.StockUuid] = q.PercentageChange
		}
		prices[q.StockUuid] = q.Price
	}

	// Step 3. Construct the movers
	movers := make([]mover, len(tRsp.GetPositions()))
	for i, pos := range tRsp.GetPositions() {
		stockUUID := pos.GetAsset().GetUuid()

		newVal := float32(prices[stockUUID] * pos.GetQuantity())
		prevVal := newVal / (1 + (oneDayChanges[stockUUID] / 100))
		chg := newVal - prevVal

		movers[i] = mover{
			stockUUID:     stockUUID,
			change:        chg,
			price:         prices[stockUUID],
			chgPercentage: oneDayChanges[stockUUID],
		}
	}

	// Step 4. Sort and return the result
	sort.SliceStable(movers, func(i, j int) bool {
		return movers[i].change > movers[j].change
	})
	return movers, nil
}

func (h *Handler) eveningPerformance(ctx context.Context, user *users.User) (*proto.Section, error) {
	// Step 1. Get the users portfolios
	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.Uuid})
	if err != nil {
		return nil, err
	}

	// Step 2. Get the movers
	movers, err := h.performanceAttribution(ctx, portfolio)
	if err != nil || len(movers) == 0 {
		return nil, err
	}

	// Step 3. Calculate the total change
	var totalChange float32
	for _, m := range movers {
		totalChange += m.change
	}

	// Step 4. Get the top and bottom movers
	gainMovers := []mover{}
	for i := 0; i < 3; i++ {
		if movers[i].change <= 0 {
			break
		}
		gainMovers = append(gainMovers, movers[i])
	}
	fallMovers := []mover{}
	for i := 0; i < 3; i++ {
		m := movers[len(movers)-i-1]
		if m.change >= 0 {
			break
		}
		fallMovers = append(fallMovers, m)
	}
	if len(gainMovers) == 0 && len(fallMovers) == 0 {
		return nil, nil
	}

	// Step 5. Fetch the metadata for the top and bottom movers
	stockUUIDs := make([]string, len(gainMovers)+len(fallMovers))
	for i, mover := range gainMovers {
		stockUUIDs[i] = mover.stockUUID
	}
	for i, mover := range fallMovers {
		stockUUIDs[i+len(gainMovers)] = mover.stockUUID
	}
	sRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return nil, err
	}
	stkMetadata := make(map[string]*stocks.Stock, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		stkMetadata[s.Uuid] = s
	}

	// Step 6. Encode the string, e.g. "Your portfolio gain was driven by Facebook (+1.5%), Amazon (2.5%) and Tesla (1.2%). You did have a few detractors Dow Chemical (-5%), Google (-3%)."
	gainStrings := make([]string, len(gainMovers))
	for i, mover := range gainMovers {
		pct := mover.chgPercentage
		md, ok := stkMetadata[mover.stockUUID]
		if !ok {
			continue
		}

		link := reactlink.Create("Stock", md.Uuid, md.Name)
		gainStrings[i] = fmt.Sprintf("%v (+%v%%)", link, math.Round(float64(pct)*10)/10)
	}
	fallStrings := make([]string, len(fallMovers))
	for i, mover := range fallMovers {
		pct := mover.chgPercentage
		md, ok := stkMetadata[mover.stockUUID]
		if !ok {
			continue
		}

		link := reactlink.Create("Stock", md.Uuid, md.Name)
		fallStrings[i] = fmt.Sprintf("%v (%v%%)", link, math.Round(float64(pct)*10)/10)
	}

	joinWithAnd := func(data []string) string {
		result := strings.Join(data[0:len(data)-1], ", ")
		return strings.Join([]string{result, data[len(data)-1]}, " and ")
	}

	var body string
	if totalChange > 0 {
		body = fmt.Sprintf("Your portfolio gain was driven by %v.", joinWithAnd(gainStrings))

		if len(fallStrings) > 0 {
			body = fmt.Sprintf("%v You did have a few detractors %v.", body, joinWithAnd(fallStrings))
		}
	} else {
		body = fmt.Sprintf("Your portfolio fall was caused by %v.", joinWithAnd(fallStrings))

		if len(fallStrings) > 0 {
			body = fmt.Sprintf("%v You did have a few saviours: %v.", body, joinWithAnd(gainStrings))
		}
	}

	return &proto.Section{Title: "Performance", Body: body}, nil
}

func (h *Handler) eveningTrades(ctx context.Context, user *users.User) ([]*proto.Trade, error) {
	// Step 1. Get the other users who the user follows
	userUUIDs, err := h.followingForResource(ctx, user, "User")
	if err != nil || len(userUUIDs) == 0 {
		return []*proto.Trade{}, err
	}

	// Step 2. Get the users portfolios
	pRsp, err := h.portfolios.List(ctx, &portfolios.ListRequest{UserUuids: userUUIDs})
	if err != nil {
		return []*proto.Trade{}, err
	}
	portfolioUUIDs := make([]string, len(pRsp.GetPortfolios()))
	for i, p := range pRsp.GetPortfolios() {
		portfolioUUIDs[i] = p.Uuid
	}

	// Step 3. Get the trades made in those portfolios today
	tRsp, err := h.trades.ListTrades(ctx, &trades.ListTradesRequest{
		PortfolioUuids: portfolioUUIDs,
		StartTime:      time.Now().Truncate(time.Hour * 24).Unix(),
		EndTime:        time.Now().Unix(),
	})
	if err != nil || len(tRsp.GetTrades()) == 0 {
		return []*proto.Trade{}, err
	}

	// Step 4. Get the portfolio UUIDs for the first X trades
	max := math.Min(5, float64(len(tRsp.GetTrades())))
	trades := tRsp.GetTrades()[0:int(max)]

	portfolioUUIDs = make([]string, int(max))
	for i, t := range trades {
		portfolioUUIDs[i] = t.PortfolioUuid
	}
	portfolioUUIDs = unique.Strings(portfolioUUIDs)

	// Step 5. Get the user UUIDs for the portfolios
	userUUIDs = []string{}
	userUUIDByPortfolioUUID := map[string]string{}

	for _, portfolioUUID := range portfolioUUIDs {
		for _, p := range pRsp.GetPortfolios() {
			if p.Uuid == portfolioUUID {
				userUUIDs = append(userUUIDs, p.UserUuid)
				userUUIDByPortfolioUUID[p.Uuid] = p.UserUuid
			}
		}
	}

	// Step 6. Get the metadata for the users
	uRsp, err := h.users.List(ctx, &users.ListRequest{Uuids: userUUIDs})
	if err != nil {
		return []*proto.Trade{}, err
	}
	usersByUUID := make(map[string]*users.User, len(uRsp.GetUsers()))
	for _, u := range uRsp.GetUsers() {
		usersByUUID[u.Uuid] = u
	}

	// Step 7. Get the stock metatdata
	stockUUIDs := make([]string, int(max))
	for i, t := range trades {
		stockUUIDs[i] = t.GetAsset().Uuid
	}
	stockUUIDs = unique.Strings(stockUUIDs)

	sRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return []*proto.Trade{}, err
	}
	stocksByUUID := make(map[string]*stocks.Stock, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		stocksByUUID[s.Uuid] = s
	}

	// Step 8. Get the value of the users portfoliso
	date := time.Now().Add(time.Hour * -24).Unix()
	vRsp, err := h.valuation.ListPriceMovements(ctx, &valuation.ListPriceMovementsRequest{
		PortfolioUuids: portfolioUUIDs, StartDate: date, EndDate: date,
	})
	if err != nil {
		return []*proto.Trade{}, err
	}
	valueByPortfolio := map[string]int64{}
	for _, v := range vRsp.GetPriceMovements() {
		valueByPortfolio[v.PortfolioUuid] = v.GetLatestValue()
	}
	fmt.Println(valueByPortfolio)

	// Step 9. Get the values of the stocks trades
	qRsp, err := h.quotes.ListQuotes(ctx, &quotes.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return []*proto.Trade{}, err
	}
	quoteByStockUUID := map[string]int64{}
	for _, q := range qRsp.GetQuotes() {
		quoteByStockUUID[q.StockUuid] = q.Price
	}

	// Step 10. Serialize the trades
	rsp := make([]*proto.Trade, len(trades))
	for i, trade := range trades {
		userUUID := userUUIDByPortfolioUUID[trade.PortfolioUuid]
		user, ok := usersByUUID[userUUID]
		if !ok {
			continue
		}

		userPortfolioValue, ok := valueByPortfolio[trade.PortfolioUuid]
		if !ok {
			continue
		}

		stock := stocksByUUID[trade.GetAsset().Uuid]
		quote := quoteByStockUUID[trade.GetAsset().Uuid]
		value := quote * trade.Quantity

		tradePercentage := math.Abs(math.Round(float64(value*100*10/userPortfolioValue)) / 10)

		tradeType := "bought"
		if trade.Type.String() == "SELL" {
			tradeType = "sold"
		}

		link := reactlink.Create("Stock", stock.Uuid, stock.Name)
		rsp[i] = &proto.Trade{
			Description: fmt.Sprintf("%v %v a %v%% position in %v (%v shares).",
				user.FirstName, tradeType, tradePercentage, link, trade.Quantity),
		}
	}
	return rsp, nil
}

func (h *Handler) eveningPosts(ctx context.Context, user *users.User) ([]*proto.Post, error) {
	// Step 1. Get recent feed items
	fRsp, err := h.feed.GetFeed(ctx, &feed.GetFeedRequest{
		Type: "User", Uuid: user.Uuid, Page: 0, Limit: 5,
	})
	if err != nil {
		return []*proto.Post{}, err
	}

	// Step 2. Get the post UUIDs
	postUUIDs := []string{}
	for _, item := range fRsp.GetItems() {
		today := time.Now().Truncate(time.Hour * 24).Unix()

		if item.CreatedAt >= today {
			postUUIDs = append(postUUIDs, item.PostUuid)
		}
	}
	if len(postUUIDs) == 0 {
		return []*proto.Post{}, nil
	}

	// Step 3. Get the post metadata
	pRsp, err := h.posts.List(ctx, &posts.ListRequest{Uuids: postUUIDs})
	if err != nil {
		return []*proto.Post{}, err
	}
	postsByUUID := make(map[string]*posts.Post, len(pRsp.GetPosts()))
	for _, p := range pRsp.GetPosts() {
		postsByUUID[p.Uuid] = p
	}

	// Step 4. Get the users who made the posts
	userUUIDs := make([]string, len(pRsp.GetPosts()))
	for i, p := range pRsp.GetPosts() {
		userUUIDs[i] = p.GetUserUuid()
	}
	userUUIDs = unique.Strings(userUUIDs)
	uRsp, err := h.users.List(ctx, &users.ListRequest{Uuids: userUUIDs})
	if err != nil {
		return []*proto.Post{}, err
	}
	usersByUUID := make(map[string]*users.User, len(uRsp.GetUsers()))
	for _, u := range uRsp.GetUsers() {
		usersByUUID[u.Uuid] = u
	}

	// Step 5. Get the assets the posts were made about
	stocksByUUID := map[string]*stocks.Stock{}
	stockUUIDs := []string{}
	for _, post := range pRsp.GetPosts() {
		if post.FeedType == "Stock" {
			stockUUIDs = append(stockUUIDs, post.FeedUuid)
		}
	}
	stockUUIDs = unique.Strings(stockUUIDs)
	if len(stockUUIDs) > 0 {
		sRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
		if err != nil {
			return []*proto.Post{}, err
		}
		for _, s := range sRsp.GetStocks() {
			stocksByUUID[s.Uuid] = s
		}
	}

	// Step 6. Serialize the posts
	rsp := []*proto.Post{}
	for _, md := range pRsp.GetPosts() {
		user, ok := usersByUUID[md.UserUuid]
		if !ok {
			fmt.Println("Missing User")
			continue
		}

		var description string
		if md.FeedType == "Stock" {
			stock, ok := stocksByUUID[md.FeedUuid]
			if !ok {
				fmt.Println("Missing Stock")
				continue
			}

			stockLink := reactlink.Create("Stock", stock.Uuid, stock.Name)
			userLink := reactlink.Create("User", user.Uuid, user.FirstName)
			description = fmt.Sprintf("%v shared a post about %v: \"%v\".", userLink, stockLink, md.Title)
		} else {
			userLink := reactlink.Create("User", user.Uuid, user.FirstName)
			description = fmt.Sprintf("%v shared a post: \"%v\".", userLink, md.Title)
		}

		rsp = append(rsp, &proto.Post{Uuid: md.Uuid, Description: description})
	}

	return rsp, nil
}

func (h *Handler) eveningBenchmarking(ctx context.Context, user *users.User) (*proto.Section, error) {
	// Your portfolio ganied 1.2% today, outperforming your investor network by 0.2%.
	// The top performer in your network today was Barry, who gained 4.5% largely driven
	// by Facebook (+11%), which they hold a 12% position in.

	// Step 1. Get the connection UUIDs
	connectionUUIDs, err := h.followingForResource(ctx, user, "User")
	if err != nil || len(connectionUUIDs) == 0 {
		return nil, err
	}

	// Step 2. Get the portfolios
	pRsp, err := h.portfolios.List(ctx, &portfolios.ListRequest{
		UserUuids: append(connectionUUIDs, user.Uuid),
	})
	fmt.Printf("There are %v connection UUIDs and %v portfolios\n", len(connectionUUIDs), len(pRsp.GetPortfolios()))
	if err != nil {
		return nil, err
	}
	portfolioUUIDs := make([]string, len(pRsp.GetPortfolios()))
	portfoliosByUserUUID := make(map[string]*portfolios.Portfolio, len(pRsp.GetPortfolios()))
	for i, p := range pRsp.GetPortfolios() {
		portfolioUUIDs[i] = p.Uuid
		portfoliosByUserUUID[p.UserUuid] = p
	}

	// Step 3. Get the one day changes
	vRsp, err := h.valuation.ListPriceMovements(ctx, &valuation.ListPriceMovementsRequest{
		PortfolioUuids: portfolioUUIDs,
		StartDate:      time.Now().Unix(),
		EndDate:        time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}

	// Step 4. Get the average & top portfolio performance, and the current users change
	var totalPctChanges float32
	var topPerformerPctChange float32
	var topPerformerPortfolioUUID string
	var userPctChange float32
	for _, m := range vRsp.GetPriceMovements() {

		if m.PortfolioUuid == portfoliosByUserUUID[user.Uuid].Uuid {
			userPctChange = m.PercentageChange
			continue
		}

		totalPctChanges += m.PercentageChange

		if m.PercentageChange > topPerformerPctChange {
			topPerformerPctChange = m.PercentageChange
			topPerformerPortfolioUUID = m.PortfolioUuid
		}
	}
	avgPctChange := totalPctChanges / float32(len(vRsp.GetPriceMovements()))

	// Step 5. Get the metadata for the top performing user
	var topPerformingUserUUID string
	for _, p := range pRsp.GetPortfolios() {
		if p.Uuid == topPerformerPortfolioUUID {
			topPerformingUserUUID = p.UserUuid
			break
		}
	}
	topPerformer, err := h.users.Find(ctx, &users.User{Uuid: topPerformingUserUUID})
	if err != nil {
		return nil, err
	}

	// Step 6. Performance attribution for top performing user
	movers, err := h.performanceAttribution(ctx, portfoliosByUserUUID[topPerformingUserUUID])
	if err != nil || len(movers) == 0 {
		return nil, err
	}
	mainMover := movers[0]
	sRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: mainMover.stockUUID})
	if err != nil {
		return nil, err
	}

	// Step XXX. Serialize the response
	userDirection := "gained"
	if userPctChange < 0 {
		userDirection = "fell"
	}

	mainMoverDirection := "is up"
	if mainMover.chgPercentage < 0 {
		mainMoverDirection = "fell just"
	}

	userPerformanceDiff := userPctChange - avgPctChange
	userPerformance := "outperforming"
	if userPerformanceDiff < 0 {
		userPerformance = "underperforming"
	}

	formatPct := func(val float32) float64 {
		return math.Round(float64(val*100)) / 100
	}

	stockLink := reactlink.Create("Stock", sRsp.GetStock().Uuid, sRsp.GetStock().Name)
	userLink := reactlink.Create("User", topPerformingUserUUID, topPerformer.FirstName)
	body := fmt.Sprintf("Your portfolio %v %v%% today, %v your investor network by %v%%. The top performer in your network today was %v who gained %v%%, largely driven by their holding in %v which %v %v%% today.",
		userDirection, formatPct(userPctChange), userPerformance, formatPct(userPerformanceDiff), userLink,
		formatPct(topPerformerPctChange), stockLink, mainMoverDirection, math.Abs(formatPct(mainMover.chgPercentage)))

	return &proto.Section{Title: "Benchmarking", Body: body}, nil
}
