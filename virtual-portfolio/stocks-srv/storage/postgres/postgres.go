package postgres

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"github.com/kytra-app/stocks-srv/storage"

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

	db.Debug()
	db.AutoMigrate(&storage.Stock{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(stock storage.Stock) (storage.Stock, error) {
	if errs := p.db.Create(&stock).GetErrors(); len(errs) > 0 {
		return stock, errs[0]
	}

	return stock, nil
}

func (p postgres) Get(query storage.Stock) (storage.Stock, error) {
	var stock storage.Stock
	errs := p.db.Table("stocks").Where(&query).First(&stock).GetErrors()

	if len(errs) > 0 {
		return stock, errs[0]
	}

	return stock, nil
}

func (p postgres) All() ([]*storage.Stock, error) {
	var stocks []*storage.Stock
	query := p.db.Table("stocks").Find(&stocks)

	if errs := query.GetErrors(); len(errs) > 0 {
		return stocks, errs[0]
	}

	return stocks, nil
}

func (p postgres) List(uuids []string) ([]*storage.Stock, error) {
	var stocks []*storage.Stock
	query := p.db.Table("stocks").Where("uuid IN (?)", uuids).Find(&stocks)

	if errs := query.GetErrors(); len(errs) > 0 {
		return stocks, errs[0]
	}

	return stocks, nil
}

func (p postgres) ListBySymbol(symbols []string) ([]*storage.Stock, error) {
	var stocks []*storage.Stock
	query := p.db.Table("stocks").Where("symbol IN (?)", symbols).Find(&stocks)

	if errs := query.GetErrors(); len(errs) > 0 {
		return stocks, errs[0]
	}

	return stocks, nil
}

func (p postgres) ListByIndustries(industries []string) ([]*storage.Stock, error) {
	var stocks []*storage.Stock
	query := p.db.Table("stocks").Where("industry IN (?)", industries).Find(&stocks)

	if errs := query.GetErrors(); len(errs) > 0 {
		return stocks, errs[0]
	}

	return stocks, nil
}

func (p postgres) Query(query string, limit int32) ([]*storage.Stock, error) {
	q := fmt.Sprintf("lower(name) LIKE lower('%%%v%%') OR lower(symbol) LIKE lower('%%%v%%')",
		query, query,
	)

	var stocks []*storage.Stock
	request := p.db.Table("stocks").Where(q).Limit(limit).Find(&stocks)

	if errs := request.GetErrors(); len(errs) > 0 {
		return stocks, errs[0]
	}

	return stocks, nil
}

func (p postgres) Delete(query storage.Stock) error {
	if errs := p.db.Delete(&query).GetErrors(); len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (p postgres) Update(params storage.Stock) (storage.Stock, error) {
	stock, err := p.Get(storage.Stock{UUID: params.UUID})
	if err != nil {
		return params, err
	}

	errs := p.db.Model(&stock).Updates(mapWithoutBlank(params)).GetErrors()
	if len(errs) > 0 {
		return stock, errs[0]
	}

	return stock, nil
}

func mapWithoutBlank(data interface{}) map[string]interface{} {
	in := structs.Map(data)
	out := make(map[string]interface{}, len(in))

	for key, val := range in {
		if val != "" && val != nil {
			out[key] = val
		}
	}

	return out
}
