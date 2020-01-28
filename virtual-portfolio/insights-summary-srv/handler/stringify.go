package handler

import (
	"fmt"
	"math"
	"time"

	trades "github.com/kytra-app/trades-srv/proto"
)

func (data *summaryData) Stringify(stockUUID string) (result string) {
	if priceMovement := data.stringifyPriceMovement(stockUUID); len(priceMovement) > 0 {
		result = result + priceMovement + " "
	}

	if priceTarget := data.stringifyPriceTargets(stockUUID); len(priceTarget) > 0 {
		result = result + priceTarget + ". "
	}

	if position := data.stringifyPosition(stockUUID); len(position) > 0 {
		if result != "" {
			result = result + "\n\n"
		}

		result = result + position + ". "
	}

	if connections := data.stringifyConnectionsWithPositions(stockUUID); len(connections) > 0 {
		if result != "" {
			result = result + "\n\n"
		}

		result = result + connections + ". "
	}

	return result
}

// stringifyPriceMovement generates the price movement part of the string for a stock, e.g:
// "Microsoft gained 5.4% yesterday" or "Facebook lost 3.2% Friday"
func (data *summaryData) stringifyPriceMovement(stockUUID string) string {
	metadata, ok := data.StkMetadata[stockUUID]
	if !ok {
		return ""
	}

	quote, ok := data.StkQuote[stockUUID]
	if !ok {
		return ""
	}

	// Ignore quote from previous days. Add 4 hours to account for timezones (after hours trading)
	// of the nasdaq stops at 8pm EST (1am GMT)/
	t := time.Unix(quote.CreatedAt, 0)
	todayStartedAt := time.Now().Truncate(time.Hour * 24).Add(4 * time.Hour)
	if t.Unix() < todayStartedAt.Unix() {
		return ""
	}

	pChg := quote.GetPercentageChange()
	changeAbs := math.Abs(float64(pChg))
	changeRounded := math.Round(changeAbs*10) / 10

	var result string

	fmt.Println(quote)
	if quote.MarketClosed {
		direction := "up"
		if pChg < 0 {
			direction = "down"
		}
		result = fmt.Sprintf("%v is %v %v%% in pre market trading", metadata.Name, direction, changeRounded)
	}

	direction := "gained"
	if pChg < 0 {
		direction = "fell"
	}
	result = fmt.Sprintf("%v %v %v%% today", metadata.Name, direction, changeRounded)

	if reason := data.stringifyReasoning(stockUUID); len(reason) > 0 {
		return fmt.Sprintf("%v %v.", result, reason)
	}
	return result + "."
}

// stringifyReasoning generates the reasoning part of the string for a stock, e.g.
// "after it come out with its Q4 earnings which beat expectations" or "ahead of its
// Q1 results today" or "" (if no reasoning).
func (data *summaryData) stringifyReasoning(stockUUID string) string {
	earningsToday, ok := data.StkHasEarning[stockUUID]
	if !ok || !earningsToday {
		return ""
	}

	quote, ok := data.StkQuote[stockUUID]
	if !ok {
		return "after it released earnings"
	}

	if quote.MarketClosed {
		return "in advance of earnings which are due to be released today"
	}

	result := "exceeded"
	if quote.GetPercentageChange() < 0 {
		result = "missed"
	}
	return fmt.Sprintf(" after it released earnings which %v expectations", result)
}

// stringifyPosition generates the position part of the string for a stock, e.g.:
// “2.5% of your portfolio is invested in Apple and your portfolio is overweight Consumer Discretionary”.
// or "You do not hold a position in Facebook. Your portfolio is underweight Information Technology."
func (data *summaryData) stringifyPosition(stockUUID string) string {
	metadata, ok := data.StkMetadata[stockUUID]
	if !ok {
		return ""
	}

	quote, ok := data.StkQuote[stockUUID]
	if !ok {
		return ""
	}

	portfolio := data.UsersPortfolio()
	var sectorTargetPercent float32
	switch metadata.Sector {
	case "Information Technology":
		sectorTargetPercent = portfolio.IndustryTargetInformationTechnology
	case "Financials":
		sectorTargetPercent = portfolio.IndustryTargetFinancials
	case "Energy":
		sectorTargetPercent = portfolio.IndustryTargetEnergy
	case "HealthCare":
		sectorTargetPercent = portfolio.IndustryTargetHealthCare
	case "Materials":
		sectorTargetPercent = portfolio.IndustryTargetMaterials
	case "Utilities":
		sectorTargetPercent = portfolio.IndustryTargetUtilities
	case "Real Estate":
		sectorTargetPercent = portfolio.IndustryTargetRealEstate
	case "Consumer Discretionary":
		sectorTargetPercent = portfolio.IndustryTargetConsumerDiscretionary
	case "Consumer Staples":
		sectorTargetPercent = portfolio.IndustryTargetConsumerStaples
	case "Communication Services":
		sectorTargetPercent = portfolio.IndustryTargetCommunicationServices
	case "Industrials":
		sectorTargetPercent = portfolio.IndustryTargetIndustrials
	default:
		sectorTargetPercent = portfolio.IndustryTargetMiscellaneous
	}

	sectorCurrentValue := data.SectorAllocations[metadata.GetSector()]
	sectorCurrentPercent := float32(100 * sectorCurrentValue / data.TotalPortfolioValue)
	sectorDiffPercentage := sectorCurrentPercent - sectorTargetPercent
	humanizeDiff := math.Abs(math.Round(float64(sectorDiffPercentage)*10) / 10)

	var result string
	if humanizeDiff < -1 {
		result = fmt.Sprintf("Your portfolio is underweight the %v sector by %v%% and", metadata.GetSector(), humanizeDiff)
	} else if sectorDiffPercentage > 1 {
		result = fmt.Sprintf("Your portfolio is overweight the %v sector by %v%% and", metadata.GetSector(), humanizeDiff)
	} else {
		result = fmt.Sprintf("Your portfolio is correctly weighted to the %v sector and", metadata.GetSector())
	}

	var position *trades.Position
	for _, p := range data.AllPositions {
		if p.PortfolioUuid != portfolio.Uuid {
			continue
		}

		if p.GetAsset().GetUuid() != stockUUID {
			continue
		}

		position = p
		break
	}
	if position == nil {
		return fmt.Sprintf("%v you do not hold a position in %v", result, metadata.GetSymbol())
	}

	posValue := float64(position.Quantity * quote.Price)
	totalVal := float64(data.TotalPortfolioValue)
	currentPercentage := posValue * 100 / totalVal

	if stockUUID == "0e7a0749-dd02-4489-b1df-d94157e10847" {
		fmt.Printf("Quantity: %v \n", position.Quantity)
		fmt.Printf("Price: %v \n", quote.Price)
		fmt.Printf("Portfolio Val: %v \n", data.TotalPortfolioValue)
		fmt.Printf("Percentage: %v \n", currentPercentage)
	}

	if currentPercentage < 1 {
		result = fmt.Sprintf("%v less than 1%% of your portfolio is invested in %v", result, metadata.GetSymbol())
	} else {
		rounded := math.Round(currentPercentage*10) / 10
		result = fmt.Sprintf("%v %v%% of your portfolio is invested in %v", result, rounded, metadata.GetSymbol())
	}

	return result
}

// stringifyPriceTargets generates the price target part of the string. e.g:
// "Price targets average $213 (+24%), rating this a strong buy."
func (data *summaryData) stringifyPriceTargets(stockUUID string) string {
	target, ok := data.StkPriceTarget[stockUUID]
	if !ok {
		return ""
	}

	quote, ok := data.StkQuote[stockUUID]
	if !ok {
		return ""
	}

	metadata, ok := data.StkMetadata[stockUUID]
	if !ok {
		return ""
	}

	priceDiffPercent := 100 * float32(target.PriceTarget-quote.Price) / float32(quote.Price)

	prefix := ""
	if priceDiffPercent >= 0 {
		prefix = "+"
	}

	var rating string
	if priceDiffPercent >= 40 {
		rating = "strong buy"
	} else if priceDiffPercent >= 20 {
		rating = "buy"
	} else if priceDiffPercent >= 10 {
		rating = "weak buy"
	} else if priceDiffPercent >= -10 {
		rating = "hold"
	} else if priceDiffPercent >= -20 {
		rating = "weak Sell"
	} else if priceDiffPercent >= -40 {
		rating = "sell"
	} else {
		rating = "strong sell"
	}

	percentage := math.Round(float64(priceDiffPercent))

	if target.NumberOfAnalysts == 1 {
		format := "%v is covered by an analyst who has set a price target of $%v (%v%v%%), rating it a %v"
		return fmt.Sprintf(format, metadata.GetSymbol(), target.PriceTarget/100, prefix, percentage, rating)
	}

	format := "%v is covered by %v analysts who have set an average price target of $%v (%v%v%%), rating it a %v"
	return fmt.Sprintf(format, metadata.GetSymbol(), target.NumberOfAnalysts, target.PriceTarget/100, prefix, percentage, rating)

}

// stringifyConnectionsWithPositions generates the connections part of the string. e.g:
// "Jack and two other investors are also invested."
func (data *summaryData) stringifyConnectionsWithPositions(stockUUID string) string {
	users := data.StkUsersWithPosition[stockUUID]
	if len(users) == 0 {
		return "None of your investor network holds a position"
	}

	firstUser := users[0]
	if len(users) == 1 {
		return fmt.Sprintf("%v holds a position", firstUser.FirstName)
	}

	suffix := ""
	if len(users) > 2 {
		suffix = "s"
	}

	format := "%v and %v other member%v of your investor network holds a position"
	return fmt.Sprintf(format, firstUser.FirstName, len(users)-1, suffix)
}
