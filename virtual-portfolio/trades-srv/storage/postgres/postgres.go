package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/trades-srv/helpers"
	"github.com/kytra-app/trades-srv/storage"
	"github.com/micro/go-micro/errors"

	// The PG driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type postgres struct{ db *gorm.DB }

// New returns an instance of storage.Service backed by a PG database
func New(host, dbname, user, password string) (storage.Service, error) {
	// Construct the DB Source
	src := fmt.Sprintf("host=%v dbname=%v user=%v password=%v port=5432 sslmode=disable",
		host, dbname, user, password)

	// Connect to the DB
	db, err := gorm.Open("postgres", src)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)
	db.AutoMigrate(&storage.Trade{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) CreateTrade(trade storage.Trade) (storage.Trade, error) {
	if err := p.ensureNoShortPosition(trade); err != nil {
		return trade, err
	}

	req := p.db.Create(&trade)
	return trade, microgorm.TranslateErrors(req)
}

func (p postgres) GetTrade(query storage.Trade) (storage.Trade, error) {
	var trade storage.Trade
	req := p.db.Where(&query).First(&trade)
	return trade, microgorm.TranslateErrors(req)
}

func (p postgres) SetTradeMetadata(query storage.Trade) (storage.Trade, error) {
	trade, err := p.GetTrade(storage.Trade{UUID: query.UUID})
	if err != nil {
		return query, err
	}

	req := p.db.Model(&trade).Updates(&query)
	return trade, microgorm.TranslateErrors(req)
}

func (p postgres) ListTrades(startTime time.Time, endTime time.Time, portfolioUUIDs []string) ([]storage.Trade, error) {
	var trades []storage.Trade

	req := p.db.Table("trades")
	req = req.Where("created_at >= ? AND created_at <= ?", startTime, endTime)
	if len(portfolioUUIDs) > 0 {
		req = req.Where("portfolio_uuid IN (?)", portfolioUUIDs)
	}
	req.Find(&trades)

	return trades, microgorm.TranslateErrors(req)
}

func (p postgres) ListTradesForPosition(position storage.Position, time time.Time) ([]storage.Trade, error) {
	query := storage.Trade{
		AssetType:     position.AssetType,
		AssetUUID:     position.AssetUUID,
		PortfolioUUID: position.PortfolioUUID,
	}

	var trades []storage.Trade
	req := p.db.Where(&query).Order("id ASC")
	req = req.Where("created_at < ?", time)
	req.Find(&trades)

	return trades, microgorm.TranslateErrors(req)
}

func (p postgres) ListTradesForPortfolio(portfolioUUID string, time time.Time) ([]storage.Trade, error) {
	var trades []storage.Trade
	req := p.db.Where(&storage.Trade{PortfolioUUID: portfolioUUID})
	req = req.Where("created_at < ?", time)
	req = req.Order("id ASC").Find(&trades)
	return trades, microgorm.TranslateErrors(req)
}

func (p postgres) ListPositionsForPortfolio(portfolioUUID string, time time.Time) ([]storage.Position, error) {
	trades, err := p.ListTradesForPortfolio(portfolioUUID, time)
	if err != nil {
		return []storage.Position{}, err
	}

	return p.groupTradesIntoPositions(trades), nil
}

func (p postgres) ListPositions(portfolioUUIDs []string, assetType string, assetUUIDs []string, time time.Time) ([]storage.Position, error) {
	var trades []storage.Trade
	tradesReq := p.db.Table("trades")
	tradesReq = tradesReq.Where("portfolio_uuid IN (?)", portfolioUUIDs)

	if assetType != "" {
		tradesReq = tradesReq.Where("asset_uuid IN (?)", assetUUIDs)
		tradesReq = tradesReq.Where("asset_type = ?", assetType)
	}

	tradesReq = tradesReq.Where("created_at < ?", time)
	tradesReq = tradesReq.Order("id ASC").Find(&trades)
	if err := microgorm.TranslateErrors(tradesReq); err != nil {
		return []storage.Position{}, err
	}

	positions := []storage.Position{}

	for _, uuid := range portfolioUUIDs {
		tradesForPortfolio := []storage.Trade{}

		for _, trade := range trades {
			if trade.PortfolioUUID == uuid {
				tradesForPortfolio = append(tradesForPortfolio, trade)
			}
		}

		positions = append(positions, p.groupTradesIntoPositions(tradesForPortfolio)...)
	}

	return p.groupTradesIntoPositions(trades), nil
}

// AllAssets returns a list of every asset that's ever been traded
func (p postgres) AllAssets() ([]storage.Asset, error) {
	var trades []storage.Trade
	query := p.db.Table("trades").Select("DISTINCT asset_uuid, asset_type")
	if err := microgorm.TranslateErrors(query.Find(&trades)); err != nil {
		return []storage.Asset{}, err
	}

	assets := make([]storage.Asset, len(trades))
	for i, t := range trades {
		assets[i] = storage.Asset{UUID: t.AssetUUID, Type: t.AssetType}
	}

	return assets, nil
}

// PrevalidateTrade runs multiple checks, to see if a trade is likely to succeed
func (p postgres) PrevalidateTrade(trade storage.Trade) error {
	if err := p.ensureNoShortPosition(trade); err != nil {
		return err
	}

	if err := trade.GetValidationError(p.db); err != nil {
		return err
	}

	return nil
}

// ensureNoShortPosition will return an error if a trade will result in a short position (negative quanity of shares)
func (p postgres) ensureNoShortPosition(trade storage.Trade) error {
	// BUY trades cannot result in a short position, assuming a current valid position
	if trade.Type == "BUY" {
		return nil
	}

	position, err := p.getPosition(trade.PortfolioUUID, trade.AssetType, trade.AssetUUID)
	if err != nil {
		return err
	} else if position.Quantity < trade.Quantity {
		return errors.BadRequest("INVALID_POSITION", "You do not have enough shares to execute this trade")
	}

	return nil
}

// getPosition returns the current quanity of shares held for a given position (Portfolio / Asset)
func (p postgres) getPosition(porfolioUUID, assetType, assetUUID string) (storage.Position, error) {
	var trades []storage.Trade
	search := storage.Trade{PortfolioUUID: porfolioUUID, AssetType: assetType, AssetUUID: assetUUID}

	query := p.db.Where(&search).Select("type, quantity").Find(&trades)
	if err := microgorm.TranslateErrors(query); err != nil {
		return storage.Position{}, err
	}

	var quantity int64
	for _, trade := range trades {
		switch trade.Type {
		case "BUY":
			quantity = quantity + trade.Quantity
		case "SELL":
			quantity = quantity - trade.Quantity
		}
	}

	result := storage.Position{
		AssetType: assetType,
		AssetUUID: assetUUID,
		Quantity:  quantity,
	}

	return result, nil
}

func (p postgres) groupTradesIntoPositions(trades []storage.Trade) (results []storage.Position) {
	// Group the trades by assets
	tradesForAsset := make(map[storage.Position][]storage.Trade)
	for _, trade := range trades {
		position := storage.Position{
			AssetUUID:     trade.AssetUUID,
			AssetType:     trade.AssetType,
			PortfolioUUID: trade.PortfolioUUID,
		}

		if res, posExists := tradesForAsset[position]; !posExists {
			tradesForAsset[position] = []storage.Trade{trade}
		} else {
			tradesForAsset[position] = append(res, trade)
		}
	}

	// Serialize the data
	for position, trades := range tradesForAsset {
		// Don't return blank positions (where number of shares bought == number of shares sold)
		if position.Quantity = helpers.SumQuantity(trades); position.Quantity == 0 {
			continue
		}

		position.BookCost = helpers.BookCost(trades)
		results = append(results, position)
	}

	return results
}
