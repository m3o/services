package handler

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	allocation "github.com/kytra-app/portfolio-allocation-srv/proto"
	valuation "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	posts "github.com/kytra-app/posts-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
)

// summariesForUsers produces a map of summaries, e.g.
//  Sam holds 24 investments across 10 sectors, they’ve netted gains of 7.2%
//	over the last month and are most heavily invested in Information Technology.
//	They have posted 12 times in the last month and traded 82 times in the past year.
//	Similar holdings to your portfolio includes Apple, Facebook, LinkedIn and 3 others.
func (h *Handler) summariesForUsers(ctx context.Context, uuids []string, allocations ...*allocation.Portfolio) (map[string]string, error) {
	// Step 0. Get the current users uuid
	user, err := h.auth.UserFromContext(ctx)
	if err != nil {
		return map[string]string{}, err
	}

	// Step 1. Get the users metadata
	uRsp, err := h.users.List(ctx, &users.ListRequest{Uuids: append(uuids, user.UUID)})
	if err != nil {
		return map[string]string{}, err
	}
	usersByUUID := make(map[string]*users.User, len(uuids))
	for _, u := range uRsp.GetUsers() {
		usersByUUID[u.Uuid] = u
	}

	// Step 2. Get the portfolio uuids
	pRsp, err := h.portfolios.List(ctx, &portfolios.ListRequest{
		UserUuids: append(uuids, user.UUID),
	})
	if err != nil {
		return map[string]string{}, err
	}
	portfolioUUIDs := make([]string, len(pRsp.GetPortfolios()))
	portfolioUUIDsByUserUUIDs := make(map[string]string, len(pRsp.GetPortfolios()))
	for i, p := range pRsp.GetPortfolios() {
		portfolioUUIDs[i] = p.Uuid
		portfolioUUIDsByUserUUIDs[p.UserUuid] = p.Uuid
	}

	// Step 3. Get the value changes for the portfolio
	startTime := time.Now().Add(time.Hour * 24 * -30).Unix()
	endTime := time.Now().Unix()
	vRsp, err := h.valuation.ListPriceMovements(ctx, &valuation.ListPriceMovementsRequest{
		PortfolioUuids: portfolioUUIDs,
		StartDate:      startTime,
		EndDate:        endTime,
	})
	if err != nil {
		return map[string]string{}, err
	}
	valueChangeByPortfolioUUID := make(map[string]float32, len(portfolioUUIDs))
	for _, v := range vRsp.GetPriceMovements() {
		valueChangeByPortfolioUUID[v.PortfolioUuid] = v.PercentageChange
	}

	// Step 4. Get the number of posts made by the users
	poRsp, err := h.posts.CountByUser(ctx, &posts.CountByUserRequest{
		UserUuids: uuids, StartTime: startTime, EndTime: endTime,
	})
	if err != nil {
		return map[string]string{}, err
	}
	postsByUserUUID := map[string]int32{}
	for _, c := range poRsp.GetCounts() {
		postsByUserUUID[c.UserUuid] = c.Count
	}

	// Step 5. Get the portfolio allocations
	if len(allocations) == 0 {
		aRsp, err := h.allocation.List(ctx, &allocation.ListRequest{Uuids: portfolioUUIDs})
		if err != nil {
			return map[string]string{}, err
		}
		allocations = aRsp.GetPortfolios()
	}
	allocationByPortfolioUUID := make(map[string]*allocation.Portfolio, len(allocations))
	for _, p := range allocations {
		allocationByPortfolioUUID[p.Uuid] = p
	}

	// Step 6. Get the trades placed in the last year
	tRsp, err := h.trades.ListTrades(ctx, &trades.ListTradesRequest{
		PortfolioUuids: portfolioUUIDs,
		StartTime:      time.Now().Add(time.Hour * 24 * -365).Unix(),
		EndTime:        time.Now().Unix(),
	})
	if err != nil {
		return map[string]string{}, err
	}
	tradesByPortfolioUUID := make(map[string]int64, len(portfolioUUIDs))
	for _, t := range tRsp.GetTrades() {
		val := tradesByPortfolioUUID[t.PortfolioUuid]
		tradesByPortfolioUUID[t.PortfolioUuid] = val + 1
	}

	// Step 7. Serialize the data
	currentUserPortfolioUUID := portfolioUUIDsByUserUUIDs[user.UUID]
	currentUserAllocation := allocationByPortfolioUUID[currentUserPortfolioUUID]
	currentUserStocks := []*allocation.Holding{}

	for _, s := range currentUserAllocation.GetSectors() {
		currentUserStocks = append(currentUserStocks, s.GetHoldings()...)
	}

	data := make(map[string]string, len(uuids))
	for _, userUUID := range uuids {
		portfolioUUID := portfolioUUIDsByUserUUIDs[userUUID]
		change := valueChangeByPortfolioUUID[portfolioUUID]
		trades := tradesByPortfolioUUID[portfolioUUID]
		posts := postsByUserUUID[userUUID]

		allocation, ok := allocationByPortfolioUUID[portfolioUUID]
		if !ok {
			continue
		}

		user, ok := usersByUUID[userUUID]
		if !ok {
			continue
		}

		var sectorsWithHoldings int
		var totalHoldings int

		var mostHeldSector string
		var mostHeldSectorValue int64

		commonHoldings := []string{}
		for _, s := range allocation.GetSectors() {
			for _, h := range s.GetHoldings() {
				for _, foo := range currentUserStocks {
					if foo.Uuid == h.Uuid {
						commonHoldings = append(commonHoldings, foo.Name)
						break
					}
				}
			}
		}

		for _, s := range allocation.GetSectors() {
			totalHoldings += len(s.GetHoldings())

			if s.Value > mostHeldSectorValue {
				mostHeldSector = s.Name
				mostHeldSectorValue = s.Value
			}

			if len(s.GetHoldings()) > 0 {
				sectorsWithHoldings++
			}
		}

		changeType := "losses"
		roundedChange := math.Abs(math.Round(float64(change)*100) / 100)
		if change >= 0 {
			changeType = "gains"
		}

		pluralize := func(val int64, str string) string {
			suffix := ""
			if val != 1 {
				suffix = "s"
			}
			return fmt.Sprintf("%v %v%v", val, str, suffix)
		}

		summary := fmt.Sprintf("%v holds %v investments across %v sectors, ", user.FirstName, totalHoldings, sectorsWithHoldings)
		summary += fmt.Sprintf("they’ve had %v of %v%% over the last month and are most heavily invested in %v. ", changeType, roundedChange, mostHeldSector)
		summary += fmt.Sprintf("They have shared %v in the last month and placed %v in the last year. ", pluralize(int64(posts), "post"), pluralize(trades, "trade"))

		if len(commonHoldings) == 0 {
			summary += "They share no common holdings with you. "
		} else {
			// TODO: Adjust to: similar holdings to your portfolio includes Apple, Facebook, LinkedIn and 3 others.
			summary += fmt.Sprintf("Shared holdings with your portfolio include %v. ", strings.Join(commonHoldings, ", "))
		}

		data[userUUID] = summary
	}

	return data, nil
}
