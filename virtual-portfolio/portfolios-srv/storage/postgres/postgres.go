package postgres

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"github.com/kytra-app/helpers/microgorm"
	"github.com/kytra-app/portfolios-srv/storage"

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
	db.AutoMigrate(&storage.Portfolio{})

	return postgres{db}, nil
}

// Close will terminate the database connection
func (p postgres) Close() error {
	return p.db.Close()
}

func (p postgres) Create(portfolio storage.Portfolio) (storage.Portfolio, error) {
	req := p.db.Create(&portfolio)
	return portfolio, microgorm.TranslateErrors(req)
}

func (p postgres) Get(query storage.Portfolio) (storage.Portfolio, error) {
	var portfolio storage.Portfolio
	req := p.db.Where(query).First(&portfolio)
	return portfolio, microgorm.TranslateErrors(req)
}

func (p postgres) Update(query storage.Portfolio) (storage.Portfolio, error) {
	portfolio, err := p.Get(storage.Portfolio{UUID: query.UUID})
	if err != nil {
		return query, err
	}

	req := p.db.Model(&portfolio).Update(mapWithoutBlank(query))
	return portfolio, microgorm.TranslateErrors(req)
}

func (p postgres) All(page, limit int64) (res []storage.Portfolio, err error) {
	req := p.db.Table("portfolios").Limit(limit).Offset(limit * page).Find(&res)
	return res, microgorm.TranslateErrors(req)
}

func (p postgres) ListByUUIDs(uuids []string) (res []storage.Portfolio, err error) {
	req := p.db.Table("portfolios").Where("uuid IN (?)", uuids).Find(&res)
	return res, microgorm.TranslateErrors(req)
}

func (p postgres) ListByUserUUIDs(uuids []string) (res []storage.Portfolio, err error) {
	req := p.db.Table("portfolios").Where("user_uuid IN (?)", uuids).Find(&res)
	return res, microgorm.TranslateErrors(req)
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
