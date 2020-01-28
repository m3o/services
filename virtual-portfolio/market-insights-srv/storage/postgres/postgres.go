package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/market-insights-srv/storage"

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

	db.AutoMigrate(&storage.Insight{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(i storage.Insight) (storage.Insight, error) {
	req := p.db.Create(&i)
	return i, microgorm.TranslateErrors(req)
}

func (p postgres) List(date time.Time) ([]storage.Insight, error) {
	var insights []storage.Insight
	req := p.db.Where(&storage.Insight{Date: date}).Find(&insights)
	return insights, microgorm.TranslateErrors(req)
}
