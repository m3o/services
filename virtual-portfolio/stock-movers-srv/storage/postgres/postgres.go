package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/stock-movers-srv/storage"

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
	db.AutoMigrate(&storage.Movement{})
	db.Model(&storage.Movement{}).AddUniqueIndex("idx_movement_stock_uuid_date", "stock_uuid", "date")

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

// Create inserts a new movement into the database
func (p postgres) Create(movement storage.Movement) (storage.Movement, error) {
	req := p.db.Create(&movement)
	return movement, microgorm.TranslateErrors(req)
}

// List retrieves all movements which occured on the date provided
func (p postgres) List(date time.Time) ([]*storage.Movement, error) {
	var movers []*storage.Movement
	req := p.db.Where(storage.Movement{Date: date}).Find(&movers)
	return movers, microgorm.TranslateErrors(req)
}
