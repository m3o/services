package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the stock-quote-srv db
type Service interface {
	Create(Quote) (Quote, error)
	Get(time.Time, string, bool) (Quote, error)
	List(time.Time, []string, bool) ([]Quote, error)
	Close() error
}

// Quote represents the price of a stock
type Quote struct {
	ID uint `gorm:"primary_key"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	StockUUID     string
	Price         int32
	ChangePercent float32
	MarketOpen    bool
}

// BeforeSave performs the validations
func (m *Quote) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(m,
		validation.Field(&m.Price, validation.Required),
		validation.Field(&m.StockUUID, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
