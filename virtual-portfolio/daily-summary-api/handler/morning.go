package handler

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	proto "github.com/kytra-app/daily-summary-api/proto"
	reactlink "github.com/kytra-app/helpers/reactlink"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	earnings "github.com/kytra-app/stock-earnings-srv/proto"
	news "github.com/kytra-app/stock-news-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
)

// Morning returns the morning summary for today
func (h *Handler) Morning(ctx context.Context, req *proto.Request, rsp *proto.MorningSummary) (err error) {
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
	if intro, err := h.morningIntroduction(ctx, user); err == nil {
		rsp.Introduction = intro
	} else {
		fmt.Println(err)
	}

	// Step 3. Attach the headlines
	if headlines, err := h.morningHeadlines(ctx); err == nil {
		rsp.Headlines = headlines
	} else {
		fmt.Println(err)
	}

	// Step 4. Attach the events
	if events, err := h.morningEvents(ctx, user); err == nil {
		rsp.Events = events
	} else {
		fmt.Println(err)
	}
	// Step 4. Attach looking back
	if lb, err := h.morningLookingBack(ctx, user); err == nil {
		rsp.LookingBack = lb
	} else {
		fmt.Println(err)
	}

	// Step 4. Attach lucky dip
	if ld, err := h.morningLuckyDip(ctx); err == nil {
		rsp.LuckyDip = ld
	} else {
		fmt.Println(err)
	}

	return nil
}

func (h *Handler) morningIntroduction(ctx context.Context, user *users.User) (*proto.Section, error) {
	// Step 1. Get the previous day (if Monday, go back to Friday)
	prevDay := time.Now().Add(time.Hour * -24)
	if prevDay.Weekday().String() == "Sunday" {
		prevDay = prevDay.Add(time.Hour * -48)
	}

	// Step 2. Write the summary (e.g. "The maket rose 1.3%...")
	summary, err := h.summaryForDate(ctx, prevDay)
	if err != nil {
		return nil, err
	}

	// Step 3. Return the result
	var format string
	switch time.Now().Weekday().String() {
	case "Monday":
		format = "Happy monday %v, we hope you had a great weekend 😁\n\n%v"
	case "Tuesday":
		format = "Good morning %v, have a great day 👍\n\n%v"
	case "Wednesday":
		format = "Good morning %v, have a wonderful Wednesday ✨\n\n%v"
	case "Thursday":
		format = "Good morning %v, just two more days until the weekend 🎉\n\n%v"
	case "Friday":
		format = "Happy Friday %v! Have a great weekend 🥳\n\n%v"
	default:
		format = "Good morning %v, have a great day. \n\n%v"
	}

	body := fmt.Sprintf(format, user.FirstName, summary)
	return &proto.Section{Body: body}, nil
}

func (h *Handler) morningHeadlines(ctx context.Context) ([]*proto.Headline, error) {
	// Step 1. Fetch the headlines
	nRsp, err := h.news.ListMarketNews(ctx, &news.ListRequest{})
	if err != nil {
		return []*proto.Headline{}, err
	}

	// Step 2. Serialize the news and return
	max := int(math.Min(3, float64(len(nRsp.Articles))))
	headlines := make([]*proto.Headline, max)
	for i, a := range nRsp.GetArticles() {
		if i >= max {
			continue
		}

		headlines[i] = &proto.Headline{
			Url:       a.ArticleUrl,
			Title:     a.Title,
			Source:    a.Source,
			CreatedAt: a.CreatedAt,
		}
	}
	return headlines, nil
}

var noEventsSection = &proto.Section{
	Title: "Events",
	Body:  "None of the stocks you're following have events today.",
}

func (h *Handler) morningEvents(ctx context.Context, user *users.User) (*proto.Section, error) {
	// Step 1. Get the stocks the user follows
	stockUUIDs, err := h.followingForResource(ctx, user, "Stock")
	if err != nil || len(stockUUIDs) == 0 {
		return nil, err
	}

	// Step 2. Get the events happening for those stocks today
	eRsp, err := h.earnings.List(ctx, &earnings.ListRequest{StockUuids: stockUUIDs})
	if err != nil {
		return nil, err
	}
	earnings := eRsp.GetEarnings()
	if len(earnings) == 0 {
		return noEventsSection, nil
	}

	// Step 3. Get the metadata for those stocks
	stockUUIDs = make([]string, len(earnings))
	for i, e := range earnings {
		stockUUIDs[i] = e.StockUuid
	}
	sRsp, err := h.stocks.List(ctx, &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		return nil, err
	}
	md := make(map[string]*stocks.Stock, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		md[s.Uuid] = s
	}

	// Step 4. Serialize the result
	var body string
	if len(earnings) == 1 {
		stk := md[earnings[0].StockUuid]
		link := reactlink.Create("Stock", stk.Uuid, stk.Name)
		body = fmt.Sprintf("%v has earnings being released today. We'll keep you updated on the markets reaction through your insights.", link)
	} else {
		names := make([]string, len(earnings))
		for i, e := range earnings {
			stk := md[e.StockUuid]
			names[i] = reactlink.Create("Stock", stk.Uuid, stk.Name)
		}

		lastName := names[len(names)-1]
		otherNames := strings.Join(names[0:len(names)-2], ", ")

		body = fmt.Sprintf(
			"%v of the stocks you follow have earnings being released today; %v and %v. We’ll keep you updated through your insights.",
			len(earnings), otherNames, lastName)
	}

	return &proto.Section{Title: "Earnings", Body: body}, nil
}

func (h *Handler) morningLookingBack(ctx context.Context, user *users.User) (*proto.Section, error) {
	// Step 1. Get the users portfolio
	portfolio, err := h.portfolios.Get(ctx, &portfolios.Portfolio{UserUuid: user.Uuid})
	if err != nil {
		return nil, err
	}

	// Step 2. Get the positions the user has in their portfolio
	tRsp, err := h.trades.ListPositionsForPortfolio(ctx, &trades.ListRequest{
		PortfolioUuid: portfolio.Uuid, IncludeMetadata: true,
	})
	if err != nil {
		return nil, err
	}
	if len(tRsp.GetPositions()) == 0 {
		return nil, nil
	}

	// Step 3. Choose a random position
	// Step 3.1 Sort the positions to keep consistency between requests
	positions := tRsp.GetPositions()
	sort.Slice(positions, func(i, j int) bool {
		posI, posJ := positions[i], positions[j]
		return posI.GetAsset().GetUuid() > posJ.GetAsset().GetUuid()
	})

	// Step 3.2 Use todays date as the key to keep consistency per day
	// and ensure repetition is kept to a minimum.
	today := time.Now().Truncate(time.Hour * 24).Unix()
	index := today % int64(len(tRsp.GetPositions()))
	position := positions[index]

	// Step 4. Get the trades for thaat position
	tRsp, err = h.trades.ListTradesForPosition(ctx, &trades.ListRequest{
		PortfolioUuid: portfolio.Uuid, Asset: position.Asset,
	})
	if err != nil {
		return nil, err
	}

	// Step 5. Get the latest quote for the stock
	quote, err := h.quotes.GetQuote(ctx, &quotes.Stock{Uuid: position.Asset.Uuid})
	if err != nil {
		return nil, err
	}

	// Step 6. Get the metadata for the stock
	sRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: position.Asset.Uuid})
	if err != nil {
		return nil, err
	}

	// Step 7. Serialize the response
	value := position.GetQuantity() * quote.GetPrice()
	profit := value - position.GetBookCost()
	profitPct := profit * 100 / position.GetBookCost()
	profitRounded := math.Round(float64(profitPct))

	firstTradeTime := time.Now().Unix()
	for _, t := range tRsp.GetTrades() {
		if t.CreatedAt < firstTradeTime {
			firstTradeTime = t.CreatedAt
		}
	}
	firstTradeDate := time.Unix(firstTradeTime, 0).Format("Mon Jan 2 2006")

	// Step 7.1 Embed the metadata
	link := reactlink.Create("Stock", sRsp.GetStock().Uuid, sRsp.GetStock().Name)

	var body string
	if profit > 0 {
		body = fmt.Sprintf("On %v you bought a position in %v. Well done %v, this has made you %v%% (+$%v).",
			firstTradeDate, link, user.FirstName, profitRounded, profit/100)
	} else {
		body = fmt.Sprintf("On %v you bought a position in %v. So far %v, this has cost you %v%% (-$%v).",
			firstTradeDate, link, user.FirstName, profitRounded, math.Abs(float64(profit)/100))
	}

	return &proto.Section{Title: "Looking Back", Body: body}, nil
}

func (h *Handler) morningLuckyDip(ctx context.Context) (*proto.Section, error) {
	// Step 1. Determine the day (don't include weekends)
	endTime := time.Now().Truncate(time.Hour * 24)
	if endTime.Weekday().String() == "Sunday" {
		// go back to friday
		endTime = endTime.Add(time.Hour * -48)
	}
	startTime := endTime.Add(time.Hour * -24)

	// Step 2. Check the cache
	if body, ok := h.luckyDipCache[startTime.Unix()]; ok {
		return &proto.Section{Title: "Lucky Dip", Body: body}, nil
	}

	// Step 3. Get all the trades which occcured on the date
	tRsp, err := h.trades.ListTrades(ctx, &trades.ListTradesRequest{
		StartTime: startTime.Unix(), EndTime: endTime.Unix(),
	})
	if err != nil {
		return nil, err
	}
	if len(tRsp.GetTrades()) == 0 {
		return nil, nil
	}

	// Step 4. Determine the most traded stock
	tradesPerStock := map[string]int{}
	for _, t := range tRsp.GetTrades() {
		val := tradesPerStock[t.Asset.Uuid]
		tradesPerStock[t.Asset.Uuid] = val + 1
	}
	var mostTradedStockUUID string
	var mostTradedStockQuantity int
	for uuid, quant := range tradesPerStock {
		if quant > mostTradedStockQuantity {
			mostTradedStockUUID = uuid
			mostTradedStockQuantity = quant
		}
	}

	// Step5. Get the metadata for the stock
	sRsp, err := h.stocks.Get(ctx, &stocks.Stock{Uuid: mostTradedStockUUID})
	if err != nil {
		return nil, err
	}
	md := sRsp.GetStock()

	// Step 6. Seriaize the result
	percent := mostTradedStockQuantity * 100 / len(tRsp.GetTrades())
	link := reactlink.Create("Stock", md.Uuid, md.Name)
	body := fmt.Sprintf(
		"%v was the most traded stock on Kytra yesterday, accounting for ~%v%% of all trades.",
		link, math.Round(float64(percent)),
	)

	// Step 7. Cache and return
	h.luckyDipCache[startTime.Unix()] = body
	return &proto.Section{Title: "Lucky Dip", Body: body}, nil
}
