package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the stock-target-price-srv db
type Service interface {
	Create(Stock) (Stock, error)
	List(time.Time, []string) ([]*Stock, error)
	StockExistsToday(string) (bool, error)
	Close() error
}

// Stock represents the target price of a stock
type Stock struct {
	ID uint `gorm:"primary_key"`

	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	UUID             string `gorm:"type:uuid"`
	PriceTarget      int64
	NumberOfAnalysts int64
}

// BeforeSave performs the validations
func (m *Stock) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(m,
		validation.Field(&m.UUID, validation.Required),
		validation.Field(&m.PriceTarget, validation.Required),
		validation.Field(&m.NumberOfAnalysts, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
