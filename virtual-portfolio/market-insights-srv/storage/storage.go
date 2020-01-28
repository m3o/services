package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the db
type Service interface {
	List(time.Time) ([]Insight, error)
	Create(Insight) (Insight, error)
	Close() error
}

// Insight is a market insight for a given stock on a given day
type Insight struct {
	gorm.Model

	Date                         time.Time
	AssetUUID                    string `gorm:"type:uuid"`
	AssetType                    string `gorm:"size:25"`
	EarningsToday                bool   `gorm:"type:boolean;DEFAULT=FALSE"`
	Score                        float32
	PrevDayPriceChangePercentage float32
}

// BeforeSave performs the validations
func (i *Insight) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(i,
		validation.Field(&i.Date, validation.Required),
		validation.Field(&i.AssetUUID, validation.Required),
		validation.Field(&i.AssetType, validation.Required))

	if err != nil {
		db.AddError(err)
		return
	}
}
