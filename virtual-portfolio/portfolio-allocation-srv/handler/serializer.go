package handler

import (
	"context"
	"fmt"

	"github.com/kytra-app/helpers/unique"
	proto "github.com/kytra-app/portfolio-allocation-srv/proto"
	valuation "github.com/kytra-app/portfolio-value-tracking-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	quotes "github.com/kytra-app/stock-quote-srv-v2/proto"
	stocks "github.com/kytra-app/stocks-srv/proto"
	trades "github.com/kytra-app/trades-srv/proto"
)

func (h Handler) serializePortfolios(ctx context.Context, data []*portfolios.Portfolio) ([]*proto.Portfolio, error) {
	// Step 0. Validate the input
	if len(data) == 0 {
		return []*proto.Portfolio{}, nil
	}

	// Step 1. Get the UUIDs of the portfolios
	uuids := make([]string, len(data))
	portfoliosByUUID := make(map[string]*portfolios.Portfolio, len(data))
	for i, p := range data {
		uuids[i] = p.Uuid
		portfoliosByUUID[p.Uuid] = p
	}

	// Step 2. Get the positions and group by the portfolio UUID
	tRsp, err := h.trades.ListPositions(ctx, &trades.BulkListRequest{
		PortfolioUuids: uuids,
	})
	if err != nil {
		return []*proto.Portfolio{}, err
	}
	positionsByPortfolioUUID := map[string][]*trades.Position{}
	for _, p := range tRsp.GetPositions() {
		x, ok := positionsByPortfolioUUID[p.PortfolioUuid]
		if !ok {
			x = []*trades.Position{}
		}
		positionsByPortfolioUUID[p.PortfolioUuid] = append(x, p)
	}

	// Step 3. Get the stock metadata (we need the name and sector)
	stockUUIDs := []string{}
	for _, p := range tRsp.GetPositions() {
		asset := p.GetAsset()

		if asset != nil && asset.Type == "Stock" {
			stockUUIDs = append(stockUUIDs, asset.Uuid)
		}
	}
	sRsp, err := h.stocks.List(ctx, &stocks.ListRequest{
		Uuids: unique.Strings(stockUUIDs),
	})
	if err != nil {
		return []*proto.Portfolio{}, err
	}
	stocksByUUID := make(map[string]*stocks.Stock, len(sRsp.GetStocks()))
	for _, s := range sRsp.GetStocks() {
		stocksByUUID[s.Uuid] = s
	}

	// Step 4. Get the quotes for the stocks
	qRsp, err := h.quotes.ListQuotes(ctx, &quotes.ListRequest{
		Uuids: stockUUIDs,
	})
	if err != nil {
		return []*proto.Portfolio{}, err
	}
	quotesByStockUUID := make(map[string]*quotes.Quote, len(qRsp.GetQuotes()))
	for _, q := range qRsp.GetQuotes() {
		quotesByStockUUID[q.StockUuid] = q
	}

	// Step 5. Get the valuations for the portfolios
	vRsp, err := h.valuation.ListValuations(ctx, &valuation.ListValuationsRequest{
		PortfolioUuids: uuids,
	})
	if err != nil {
		return []*proto.Portfolio{}, err
	}
	valuationsByPortfolioUUID := make(map[string]int64, len(vRsp.GetValuations()))
	for _, v := range vRsp.GetValuations() {
		valuationsByPortfolioUUID[v.PortfolioUuid] = v.Amount
	}

	// Step 6. Serialize the portfolios
	result := []*proto.Portfolio{}
	for _, uuid := range uuids {
		portfolio, ok := portfoliosByUUID[uuid]
		if !ok {
			fmt.Printf("Missing Portfolio: %v\n", uuid)
			continue
		}

		valuation, ok := valuationsByPortfolioUUID[uuid]
		if !ok {
			fmt.Printf("Missing Valuation: %v\n", uuid)
			continue
		}

		positions, ok := positionsByPortfolioUUID[uuid]
		if !ok {
			positions = []*trades.Position{}
		}

		var stockValue int64
		valueBySector := map[string]int64{}

		holdings := make([]*proto.Holding, len(positions))
		for i, pos := range positions {
			stock, ok := stocksByUUID[pos.GetAsset().GetUuid()]
			if !ok {
				continue
			}

			quote, ok := quotesByStockUUID[pos.GetAsset().GetUuid()]
			if !ok {
				continue
			}

			posValue := pos.Quantity * quote.Price
			stockValue += posValue

			sectorVal := valueBySector[stock.Sector]
			valueBySector[stock.Sector] = sectorVal + posValue

			holdings[i] = &proto.Holding{
				Type:      "Stock",
				Uuid:      stock.Uuid,
				Name:      stock.Name,
				Sector:    stock.Sector,
				UnitPrice: quote.Price,
				Value:     pos.Quantity * quote.Price,
			}
		}

		holdingsBySector := map[string][]*proto.Holding{}

		for _, h := range holdings {
			h.PercentOfSector = calcPercent(h.Value, valueBySector[h.Sector])
			h.PercentOfPortfolio = calcPercent(h.Value, valuation)
			h.PercentOfAssetClass = calcPercent(h.Value, stockValue)

			holdings, ok := holdingsBySector[h.Sector]
			if !ok {
				holdings = []*proto.Holding{}
			}
			holdingsBySector[h.Sector] = append(holdings, h)
		}

		result = append(result, &proto.Portfolio{
			Uuid:     uuid,
			UserUuid: portfolio.UserUuid,
			AssetClasses: []*proto.Sector{
				&proto.Sector{
					Name:                     "Stocks",
					TargetPercentOfPortfolio: portfolio.AssetClassTargetStocks,
					Value:                    stockValue,
					PercentOfPortfolio:       calcPercent(stockValue, valuation),
					Holdings:                 holdings,
				},
				&proto.Sector{
					Name:                     "Cash",
					Value:                    valuation - stockValue,
					PercentOfPortfolio:       calcPercent(valuation-stockValue, valuation),
					TargetPercentOfPortfolio: portfolio.AssetClassTargetCash,
				},
			},
			Sectors: []*proto.Sector{
				&proto.Sector{
					Name:                     "Information Technology",
					Value:                    valueBySector["Information Technology"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetInformationTechnology(),
					PercentOfPortfolio:       calcPercent(valueBySector["Information Technology"], valuation),
					Holdings:                 holdingsBySector["Information Technology"],
				},
				&proto.Sector{
					Name:                     "Financials",
					Value:                    valueBySector["Financials"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetFinancials(),
					PercentOfPortfolio:       calcPercent(valueBySector["Financials"], valuation),
					Holdings:                 holdingsBySector["Financials"],
				},
				&proto.Sector{
					Name:                     "Energy",
					Value:                    valueBySector["Energy"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetEnergy(),
					PercentOfPortfolio:       calcPercent(valueBySector["Energy"], valuation),
					Holdings:                 holdingsBySector["Energy"],
				},
				&proto.Sector{
					Name:                     "Health Care",
					Value:                    valueBySector["Health Care"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetHealthCare(),
					PercentOfPortfolio:       calcPercent(valueBySector["Health Care"], valuation),
					Holdings:                 holdingsBySector["Health Care"],
				},
				&proto.Sector{
					Name:                     "Materials",
					Value:                    valueBySector["Materials"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetMaterials(),
					PercentOfPortfolio:       calcPercent(valueBySector["Materials"], valuation),
					Holdings:                 holdingsBySector["Materials"],
				},
				&proto.Sector{
					Name:                     "Utilities",
					Value:                    valueBySector["Utilities"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetUtilities(),
					PercentOfPortfolio:       calcPercent(valueBySector["Utilities"], valuation),
					Holdings:                 holdingsBySector["Utilities"],
				},
				&proto.Sector{
					Name:                     "Real Estate",
					Value:                    valueBySector["Real Estate"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetRealEstate(),
					PercentOfPortfolio:       calcPercent(valueBySector["Real Estate"], valuation),
					Holdings:                 holdingsBySector["Real Estate"],
				},
				&proto.Sector{
					Name:                     "Consumer Discretionary",
					Value:                    valueBySector["Consumer Discretionary"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetConsumerDiscretionary(),
					PercentOfPortfolio:       calcPercent(valueBySector["Consumer Discretionary"], valuation),
					Holdings:                 holdingsBySector["Consumer Discretionary"],
				},
				&proto.Sector{
					Name:                     "Consumer Staples",
					Value:                    valueBySector["Consumer Staples"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetConsumerStaples(),
					PercentOfPortfolio:       calcPercent(valueBySector["Consumer Staples"], valuation),
					Holdings:                 holdingsBySector["Consumer Staples"],
				},
				&proto.Sector{
					Name:                     "Communication Services",
					Value:                    valueBySector["Communication Services"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetCommunicationServices(),
					PercentOfPortfolio:       calcPercent(valueBySector["Communication Services"], valuation),
					Holdings:                 holdingsBySector["Communication Services"],
				},
				&proto.Sector{
					Name:                     "Industrials",
					Value:                    valueBySector["Industrials"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetIndustrials(),
					PercentOfPortfolio:       calcPercent(valueBySector["Industrials"], valuation),
					Holdings:                 holdingsBySector["Industrials"],
				},
				&proto.Sector{
					Name:                     "Miscellaneous",
					Value:                    valueBySector["Miscellaneous"],
					TargetPercentOfPortfolio: portfolio.GetIndustryTargetMiscellaneous(),
					PercentOfPortfolio:       calcPercent(valueBySector["Miscellaneous"], valuation),
					Holdings:                 holdingsBySector["Miscellaneous"],
				},
			},
		})
	}

	return result, nil
}

func calcPercent(valA, valB int64) float32 {
	if valB == 0 {
		return 0
	}

	return float32(valA) * 100 / float32(valB)
}
