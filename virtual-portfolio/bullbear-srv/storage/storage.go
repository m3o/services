package storage

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the bullbear-srv db
type Service interface {
	Get(Resource) (Resource, error)
	List(string, []string, string) ([]*Resource, error)
	Create(Opinion) error
	Close() error
}

// Resource is a object which can users have opinions on
type Resource struct {
	UUID       string
	Type       string
	Bulls      []string
	Bears      []string
	BullsCount int
	BearsCount int
	Opinion    string
}

// Opinion is the opinion held on a resource by a user
type Opinion struct {
	gorm.Model

	Opinion      string `gorm:"type:varchar(10)"`
	ResourceUUID string `gorm:"type:varchar(100)"`
	ResourceType string `gorm:"type:varchar(100)"`
	UserUUID     string `gorm:"type:varchar(100)"`
}

// BeforeSave performs the validations
func (o *Opinion) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(o,
		validation.Field(&o.ResourceUUID, validation.Required),
		validation.Field(&o.ResourceType, validation.Required),
		validation.Field(&o.UserUUID, validation.Required),
		validation.Field(&o.Opinion, validation.Required, validation.In("NONE", "BULLISH", "BEARISH")))

	if err != nil {
		db.AddError(err)
		return
	}
}
