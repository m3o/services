package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the portfolios-srv db
type Service interface {
	Create(Valuation) (Valuation, error)
	GetDailyHistory(string) ([]Valuation, error)
	GetIntradayHistory(string, time.Time) ([]Valuation, error)
	GetPriceMovements([]string, time.Time, time.Time) ([]DailyPriceChance, error)
	ListValuations([]string, time.Time) ([]*Valuation, error)
	Close() error
}

// Valuation is a snapshot of a portfolios  value
type Valuation struct {
	gorm.Model

	PortfolioUUID string `gorm:"type:uuid"`
	Value         int64  `gorm:"type:int"`
	Date          time.Time
}

// DailyPriceChance is the change in a portfolios value over a given day
type DailyPriceChance struct {
	PortfolioUUID string
	EarliestValue int64
	LatestValue   int64
}

// PercentageChange from EarliestValue to LatestValue
func (d DailyPriceChance) PercentageChange() float32 {
	return float32(d.LatestValue-d.EarliestValue) * 100 / float32(d.EarliestValue)
}

// BeforeSave performs the validations
func (s *Valuation) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(s,
		validation.Field(&s.PortfolioUUID, validation.Required),
		validation.Field(&s.Date, validation.Required),
		validation.Field(&s.Value, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
