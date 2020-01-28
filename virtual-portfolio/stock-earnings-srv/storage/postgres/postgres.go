package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/stock-earnings-srv/storage"

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
	db.AutoMigrate(&storage.Earning{})
	db.Model(&storage.Earning{}).AddUniqueIndex("idx_event_stock_uuid_date", "stock_uuid", "date")

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

// Create inserts a new Earning into the database
func (p postgres) Create(Earning storage.Earning) (storage.Earning, error) {
	req := p.db.Create(&Earning)
	return Earning, microgorm.TranslateErrors(req)
}

// List retrieves all earnings which occured on the date provided
func (p postgres) List(date time.Time, stockUUIDs []string) ([]*storage.Earning, error) {
	var earnings []*storage.Earning

	req := p.db.Where(storage.Earning{Date: date})
	if stockUUIDs != nil {
		req = req.Where("stock_uuid IN (?)", stockUUIDs)
	}
	req = req.Find(&earnings)

	return earnings, microgorm.TranslateErrors(req)
}
