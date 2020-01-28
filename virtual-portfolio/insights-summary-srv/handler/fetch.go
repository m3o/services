package handler

import (
	followers "github.com/kytra-app/followers-srv/proto"
	valuation "github.com/kytra-app/portfolio-valuation-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	earnings "github.com/kytra-app/stock-earnings-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	target "github.com/kytra-app/stock-target-price-srv/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
)

func (data *summaryData) Fetch() error {
	if err := data.FetchStocks(); err != nil {
		return err
	}

	if err := data.FetchStkPriceTarget(); err != nil {
		return err
	}

	if err := data.FetchConnections(); err != nil {
		return err
	}

	if err := data.FetchEarnings(); err != nil {
		return err
	}

	if err := data.FetchQuotes(); err != nil {
		return err
	}

	// Requires Connections
	if err := data.FetchPortfolios(); err != nil {
		return err
	}

	// Requires Portfolios
	if err := data.FetchValuation(); err != nil {
		return err
	}

	// Requires Portfolios
	if err := data.FetchPositions(); err != nil {
		return err
	}

	// Requires Quotes, Stocks & Portfolios
	if err := data.FetchCurrentSectorAllocation(); err != nil {
		return err
	}

	return nil
}

func (data *summaryData) FetchCurrentSectorAllocation() error {
	tRsp, err := data.Handler.trades.ListPositionsForPortfolio(data.Context, &trades.ListRequest{
		PortfolioUuid: data.UsersPortfolio().GetUuid(),
	})

	if err != nil {
		return err
	}

	data.SectorAllocations = map[string]int64{}

	for _, p := range tRsp.GetPositions() {
		stockUUID := p.GetAsset().GetUuid()

		md, ok := data.StkMetadata[stockUUID]
		if !ok {
			continue
		}

		quote, ok := data.StkQuote[stockUUID]
		if !ok {
			continue
		}

		posValue := quote.GetPrice() * p.GetQuantity()

		value, ok := data.SectorAllocations[md.GetSector()]
		if !ok {
			value = posValue
		} else {
			value += posValue
		}

		data.SectorAllocations[md.GetSector()] = value
	}

	return nil
}

func (data *summaryData) FetchConnections() error {
	user := &followers.Resource{Uuid: data.UserUUID, Type: "User"}
	fRsp, err := data.Handler.followers.Get(data.Context, user)
	if err != nil {
		return err
	}

	uuids := []string{}
	for _, f := range fRsp.GetFollowing() {
		if f.Type == "User" {
			uuids = append(uuids, f.Uuid)
		}
	}

	uRsp, err := data.Handler.users.List(data.Context, &users.ListRequest{
		Uuids: append(uuids, data.UserUUID),
	})

	if err != nil {
		return err
	}
	data.AllUsers = uRsp.GetUsers()

	return nil
}

func (data *summaryData) FetchPortfolios() error {
	allUserUUIDs := []string{}
	for _, user := range data.AllUsers {
		allUserUUIDs = append(allUserUUIDs, user.Uuid)
	}
	allUserUUIDs = append(allUserUUIDs, data.UserUUID)

	pRsp, err := data.Handler.portfolios.List(data.Context, &portfolios.ListRequest{
		UserUuids: allUserUUIDs,
	})

	if err != nil {
		return err
	}

	data.AllPortfolios = pRsp.GetPortfolios()
	return nil
}

func (data *summaryData) FetchPositions() error {
	portfolioUUIDs := make([]string, len(data.AllPortfolios))
	for i, p := range data.AllPortfolios {
		portfolioUUIDs[i] = p.GetUuid()
	}

	tRsp, err := data.Handler.trades.ListPositions(data.Context, &trades.BulkListRequest{
		PortfolioUuids: portfolioUUIDs,
		AssetUuids:     data.StockUUIDs,
		AssetType:      "Stock",
	})

	if err != nil {
		return err
	}

	result := make(map[string][]*users.User)
	for _, p := range tRsp.GetPositions() {
		user := data.UserForPosition(p)

		if user.Uuid == data.UserUUID {
			continue
		}

		usrz, ok := result[p.GetAsset().GetUuid()]
		if !ok {
			usrz = []*users.User{user}
		} else {
			usrz = append(usrz, user)
		}
		result[p.GetAsset().GetUuid()] = usrz
	}

	data.AllPositions = tRsp.GetPositions()
	data.StkUsersWithPosition = result

	return nil
}

func (data *summaryData) FetchStocks() error {
	sRsp, err := data.Handler.stocks.List(data.Context, &stocks.ListRequest{
		Uuids: data.StockUUIDs,
	})

	if err != nil {
		return err
	}

	data.StkMetadata = make(map[string]*stocks.Stock, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		data.StkMetadata[s.Uuid] = s
	}

	return nil
}

func (data *summaryData) FetchEarnings() error {
	eRsp, err := data.Handler.earnings.List(data.Context, &earnings.ListRequest{
		StockUuids: data.StockUUIDs,
	})

	if err != nil {
		return err
	}

	data.StkHasEarning = make(map[string]bool, len(data.StockUUIDs))
	for _, uuid := range data.StockUUIDs {
		data.StkHasEarning[uuid] = false
	}

	for _, e := range eRsp.GetEarnings() {
		data.StkHasEarning[e.GetStockUuid()] = true
	}

	return nil
}

func (data *summaryData) FetchStkPriceTarget() error {
	ptRsp, err := data.Handler.targets.List(data.Context, &target.ListRequest{
		StockUuids: data.StockUUIDs,
	})
	if err != nil {
		return err
	}

	data.StkPriceTarget = make(map[string]*target.Stock, len(data.StockUUIDs))
	for _, s := range ptRsp.GetStocks() {
		data.StkPriceTarget[s.GetUuid()] = s
	}

	return nil
}

func (data *summaryData) FetchQuotes() error {
	qRsp, err := data.Handler.quotes.ListQuotes(data.Context, &quotes.ListRequest{
		Uuids: data.StockUUIDs, IncludeOutOfHours: true,
	})
	if err != nil {
		return err
	}

	data.StkQuote = make(map[string]*quotes.Quote, len(data.StockUUIDs))
	for _, q := range qRsp.GetQuotes() {
		data.StkQuote[q.GetStockUuid()] = q
	}
	return nil
}

func (data *summaryData) FetchValuation() error {
	vRsp, err := data.Handler.valuation.GetPortfolio(data.Context, &valuation.Portfolio{
		Uuid: data.UsersPortfolio().GetUuid(),
	})

	if err != nil {
		return err
	}
	data.TotalPortfolioValue = vRsp.TotalValue

	return nil
}

func (data *summaryData) UserForPosition(pos *trades.Position) *users.User {
	var portfolio *portfolios.Portfolio

	for _, p := range data.AllPortfolios {
		if p.Uuid == pos.GetPortfolioUuid() {
			portfolio = p
			break
		}
	}

	for _, u := range data.AllUsers {
		if u.Uuid == portfolio.GetUserUuid() {
			return u
		}
	}

	return nil
}

func (data *summaryData) UsersPortfolio() *portfolios.Portfolio {
	for _, p := range data.AllPortfolios {
		if p.UserUuid == data.UserUUID {
			return p
		}
	}

	return nil
}
