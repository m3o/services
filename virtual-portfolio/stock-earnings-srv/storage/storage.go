package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the stock-earnings-srv db
type Service interface {
	Create(Earning) (Earning, error)
	List(time.Time, []string) ([]*Earning, error)
	Close() error
}

// Earning is an earnings event
type Earning struct {
	gorm.Model

	Date      time.Time `gorm:"unique_index:idx_event_stock_uuid_date"`
	StockUUID string    `gorm:"type:uuid;unique_index:idx_event_stock_uuid_date"`
}

// BeforeSave performs the validations
func (e *Earning) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(e,
		validation.Field(&e.StockUUID, validation.Required),
		validation.Field(&e.Date, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
