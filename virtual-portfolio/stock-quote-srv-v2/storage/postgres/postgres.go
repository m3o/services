package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/stock-quote-srv/storage"

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
	db.AutoMigrate(&storage.Quote{})

	// CREATE INDEX index_stock_uuid_desc ON quotes (stock_uuid DESC)

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

// Create inserts a new quote into the database
func (p postgres) Create(quote storage.Quote) (storage.Quote, error) {
	req := p.db.Create(&quote)
	return quote, microgorm.TranslateErrors(req)
}

// List retrieves the latest quotes for each stock
func (p postgres) List(date time.Time, stockUUIDs []string, includeOutOfHours bool) ([]storage.Quote, error) {
	var quotes []storage.Quote

	req := p.db.Table("quotes")
	req = req.Where("stock_uuid IN (?)", stockUUIDs)
	req = req.Where("created_at < ? AND created_at >= ?", date, date.Add(time.Hour*24*-7))
	req = req.Select("DISTINCT ON (stock_uuid) stock_uuid, price, created_at, change_percent, market_open")
	if !includeOutOfHours {
		req = req.Where("market_open = ?", true)
	}
	req = req.Order("stock_uuid, created_at DESC").Find(&quotes)

	return quotes, microgorm.TranslateErrors(req)
}

// Get retrieves the latest quote for a stock
func (p postgres) Get(date time.Time, uuid string, includeOutOfHours bool) (storage.Quote, error) {
	var quote storage.Quote

	req := p.db.Table("quotes")
	req = req.Where("stock_uuid = ?", uuid)
	req = req.Where("created_at < ?", date)
	if !includeOutOfHours {
		req = req.Where("market_open = ?", true)
	}
	req = req.Order("created_at DESC").First(&quote)

	return quote, microgorm.TranslateErrors(req)
}
