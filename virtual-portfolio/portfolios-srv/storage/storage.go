package storage

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Service is an wrapper of the portfolios-srv db
type Service interface {
	Create(Portfolio) (Portfolio, error)
	Get(Portfolio) (Portfolio, error)
	Update(Portfolio) (Portfolio, error)
	All(int64, int64) ([]Portfolio, error)
	ListByUUIDs([]string) ([]Portfolio, error)
	ListByUserUUIDs([]string) ([]Portfolio, error)
	Close() error
}

// Portfolio is a group of investments
type Portfolio struct {
	gorm.Model

	UUID     string `gorm:"type:uuid;primary_key;"`
	UserUUID string `gorm:"type:uuid"`

	AssetClassTargetStocks              float32 `gorm:"type:numeric(4,2)"`
	AssetClassTargetCash                float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetInformationTechnology float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetFinancials            float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetEnergy                float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetHealthCare            float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetMaterials             float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetUtilities             float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetRealEstate            float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetConsumerDiscretionary float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetConsumerStaples       float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetCommunicationServices float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetIndustrials           float32 `gorm:"type:numeric(4,2)"`
	IndustryTargetMiscellaneous         float32 `gorm:"type:numeric(4,2)"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (s *Portfolio) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (s *Portfolio) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(s,
		validation.Field(&s.UserUUID, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
