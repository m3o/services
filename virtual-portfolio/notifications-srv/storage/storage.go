package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Service is an wrapper of the notifications-srv db
type Service interface {
	Create(Notification) (Notification, error)
	List(Query) ([]*Notification, error)
	SetNotificationsSeen(string) error
	Close() error
}

// Query is an object containing the various attributes which can be
// used to query the Notification database.
type Query struct {
	UserUUID   string
	StartTime  *time.Time
	EndTime    *time.Time
	Page       int64
	Limit      int64
	OnlyUnseen bool
}

// Notification is a message sent to the user to notify them of
// an event on the Kytra plaltform
type Notification struct {
	gorm.Model

	UUID         string `gorm:"type:uuid;primary_key;"`
	UserUUID     string `gorm:"type:uuid"`
	Title        string `gorm:"type:text"`
	Description  string `gorm:"type:text"`
	Seen         bool   `gorm:"type:text;default:false"`
	Emoji        string `gorm:"type:varchar(10)"`
	ResourceType string `gorm:"type:varchar(100)"`
	ResourceUUID string `gorm:"type:uuid"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (p *Notification) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (p *Notification) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(p,
		validation.Field(&p.UserUUID, validation.Required),
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.ResourceType, validation.Required),
		validation.Field(&p.ResourceUUID, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
