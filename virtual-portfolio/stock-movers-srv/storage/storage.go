package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the stock-movers-srv db
type Service interface {
	Create(Movement) (Movement, error)
	List(time.Time) ([]*Movement, error)
	Close() error
}

// Movement represents the price movement of a stock
type Movement struct {
	gorm.Model

	Percentage float32   `json:"percentage"`
	Date       time.Time `gorm:"unique_index:idx_name_code" json:"date"`
	StockUUID  string    `gorm:"type:uuid;unique_index:idx_name_code" json:"stock_uuid"`
}

// BeforeSave performs the validations
func (m *Movement) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(m,
		validation.Field(&m.Percentage, validation.Required),
		validation.Field(&m.StockUUID, validation.Required),
		validation.Field(&m.Date, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
