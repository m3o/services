package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/portfolio-value-tracking-srv/storage"

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
	db.AutoMigrate(&storage.Valuation{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(valuation storage.Valuation) (storage.Valuation, error) {
	req := p.db.Create(&valuation)
	return valuation, microgorm.TranslateErrors(req)
}

func (p postgres) GetDailyHistory(portfolioUUID string) (res []storage.Valuation, err error) {
	req := p.db.Table("valuations")
	req = req.Where("portfolio_uuid = ?", portfolioUUID)
	req = req.Select("DISTINCT ON(created_at::DATE) date, value")
	req = req.Order("created_at::DATE, created_at DESC").Limit(90).Find(&res)

	return res, microgorm.TranslateErrors(req)
}

func (p postgres) GetIntradayHistory(portfolioUUID string, date time.Time) (res []storage.Valuation, err error) {
	startTime := date.Truncate(time.Hour * 24)
	endTime := startTime.Add(time.Hour * 24)

	req := p.db.Table("valuations")
	req = req.Where("portfolio_uuid = ?", portfolioUUID)
	req = req.Where("created_at >= ? AND created_at < ?", startTime, endTime)
	req = req.Order("created_at DESC").Find(&res)

	return res, microgorm.TranslateErrors(req)
}

func (p postgres) ListValuations(portfolioUUIDs []string, t time.Time) ([]*storage.Valuation, error) {
	var prices []*storage.Valuation
	q := p.db.Table("valuations")
	q = q.Where("portfolio_uuid IN (?)", portfolioUUIDs)
	q = q.Where("created_at < ? AND created_at > ?", t, t.Add(time.Hour*-24))
	q = q.Order("portfolio_uuid, created_at DESC")
	q = q.Select("DISTINCT ON(portfolio_uuid) portfolio_uuid, created_at, value").Find(&prices)
	return prices, microgorm.TranslateErrors(q)
}

func (p postgres) GetPriceMovements(portfolioUUIDs []string, startDate time.Time, endDate time.Time) ([]storage.DailyPriceChance, error) {
	startTime := startDate.Truncate(time.Hour * 24)
	endTime := endDate.Truncate(time.Hour * 24).Add(time.Hour * 24)

	var earliestPrices []*storage.Valuation
	reqOne := p.db.Table("valuations")
	reqOne = reqOne.Where("portfolio_uuid IN (?)", portfolioUUIDs)
	reqOne = reqOne.Where("created_at < ? AND created_at > ?", startTime, startTime.Add(time.Hour*-48))
	reqOne = reqOne.Order("portfolio_uuid, created_at DESC")
	reqOne = reqOne.Select("DISTINCT ON(portfolio_uuid) portfolio_uuid, date, value").Find(&earliestPrices)
	if err := microgorm.TranslateErrors(reqOne); err != nil {
		return []storage.DailyPriceChance{}, err
	}

	var latestPrices []*storage.Valuation
	reqTwo := p.db.Table("valuations")
	reqTwo = reqTwo.Where("portfolio_uuid IN (?)", portfolioUUIDs)
	reqTwo = reqTwo.Where("created_at < ? AND created_at > ?", endTime, endTime.Add(time.Hour*-48))
	reqTwo = reqTwo.Order("portfolio_uuid, created_at DESC")
	reqTwo = reqTwo.Select("DISTINCT ON(portfolio_uuid) portfolio_uuid, date, value").Find(&latestPrices)
	if err := microgorm.TranslateErrors(reqTwo); err != nil {
		return []storage.DailyPriceChance{}, err
	}

	result := []storage.DailyPriceChance{}

	for _, uuid := range portfolioUUIDs {
		var earliestPrice *storage.Valuation
		for _, p := range earliestPrices {
			if p.PortfolioUUID == uuid {
				earliestPrice = p
				break
			}
		}
		if earliestPrice == nil {
			continue
		}

		var endPrice *storage.Valuation
		for _, p := range latestPrices {
			if p.PortfolioUUID == uuid {
				endPrice = p
				break
			}
		}
		if endPrice == nil {
			continue
		}

		result = append(result, storage.DailyPriceChance{
			PortfolioUUID: uuid,
			EarliestValue: earliestPrice.Value,
			LatestValue:   endPrice.Value,
		})
	}

	return result, nil
}
