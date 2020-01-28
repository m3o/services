package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/stock-target-price-srv/storage"

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
	db.AutoMigrate(&storage.Stock{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

// Create inserts a new price target into the database
func (p postgres) Create(stock storage.Stock) (storage.Stock, error) {
	req := p.db.Create(&stock)
	return stock, microgorm.TranslateErrors(req)
}

// List retrieves the price targets stock which occured on or before the date provided
func (p postgres) List(t time.Time, stockUUIDs []string) ([]*storage.Stock, error) {
	var stocks []*storage.Stock

	req := p.db.Table("stocks")
	req = req.Where("uuid IN (?)", stockUUIDs)
	req = req.Where("created_at <= ?", t)
	req = req.Order("uuid, created_at DESC")
	req = req.Select("DISTINCT ON(uuid) uuid, price_target, number_of_analysts, created_at")
	req = req.Find(&stocks)

	return stocks, microgorm.TranslateErrors(req)
}

func (p postgres) StockExistsToday(uuid string) (bool, error) {
	var count int

	req := p.db.Table("stocks")
	req = req.Where("uuid = ?", uuid)
	req = req.Where("created_at >= ?", time.Now().Truncate(time.Hour*24))
	req = req.Count(&count)

	err := microgorm.TranslateErrors(req)
	if err == nil && count > 0 {
		return true, nil
	}
	return false, err
}
