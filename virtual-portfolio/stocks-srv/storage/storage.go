package storage

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Service is an wrapper of the followers-srv db
type Service interface {
	Create(Stock) (Stock, error)
	Get(Stock) (Stock, error)
	Update(Stock) (Stock, error)
	List([]string) ([]*Stock, error)
	ListBySymbol([]string) ([]*Stock, error)
	ListByIndustries([]string) ([]*Stock, error)
	Query(string, int32) ([]*Stock, error)
	All() ([]*Stock, error)
	Delete(Stock) error
	Close() error
}

// Stock is an equity asset class
type Stock struct {
	gorm.Model

	UUID             string `gorm:"type:uuid;primary_key;"`
	Name             string `gorm:"type:varchar(250)"`
	Symbol           string `gorm:"type:varchar(15)"`
	Exchange         string `gorm:"type:varchar(6)"`
	Type             string `gorm:"type:varchar(6)"`
	Region           string `gorm:"type:varchar(6)"`
	Currency         string `gorm:"type:varchar(6)"`
	Color            string `gorm:"type:varchar(10)"`
	ProfilePictureID string `gorm:"type:varchar(100)"`
	Industry         string `gorm:"type:varchar(100)"`
	Website          string `gorm:"type:varchar(250)"`
	Description      string `gorm:"type:text"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (s *Stock) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (s *Stock) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Symbol, validation.Required),
		validation.Field(&s.Exchange, validation.Required),
		validation.Field(&s.Type, validation.Required),
		validation.Field(&s.Region, validation.Required),
		validation.Field(&s.Currency, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
