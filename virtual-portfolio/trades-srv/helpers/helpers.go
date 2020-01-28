package helpers


import (
	"github.com/kytra-app/trades-srv/storage"
)

// SumQuantity determines the size of a position, given a slice of trades
func SumQuantity(trades []storage.Trade) (result int64) {
	for _, trade := range trades {
		switch trade.Type {
		case "BUY":
			result = result + trade.Quantity
		case "SELL":
			result = result - trade.Quantity
		}
	}
	return result
}

// BookCost determines the book cost for a position, given a slice of trades
func BookCost(trades []storage.Trade) (result int64) {
	var activeBuys []storage.Trade

	for _, trade := range trades {
		switch trade.Type {
		case "BUY":
			result = result + (trade.Quantity * trade.UnitPrice)
			activeBuys = append(activeBuys, trade)
		case "SELL":
			for {
				// An invalid value input was provided
				if len(activeBuys) == 0 {
					return -1
				}

				buy := activeBuys[len(activeBuys)-1]

				if buy.Quantity <= trade.Quantity {
					activeBuys = activeBuys[:len(activeBuys)-1]
					trade.Quantity = trade.Quantity - buy.Quantity
					result = result - (buy.UnitPrice * buy.Quantity)
				} else {
					buy.Quantity = buy.Quantity - trade.Quantity
					result = result - (buy.UnitPrice * trade.Quantity)
					break
				}
			}
		}
	}

	return result
}