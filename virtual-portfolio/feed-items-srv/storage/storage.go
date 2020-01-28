package storage

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Service is an wrapper of the followers-srv db
type Service interface {
	Create(FeedItem) (FeedItem, error)
	Get(FeedItem) (FeedItem, error)
	Delete(FeedItem) error
	BulkDelete(string, string, []string) error
	GetFeed(string, string, int32, int32) ([]*FeedItem, error)
	Close() error
}

// FeedItem is a item in a feed
type FeedItem struct {
	gorm.Model

	UUID        string `gorm:"type:uuid;primary_key;"`
	FeedType    string `gorm:"type:varchar(100)"`
	FeedUUID    string `gorm:"type:varchar(100)"`
	Tag         string `gorm:"type:varchar(25)"`
	PostUUID    string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:varchar(250)"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (i *FeedItem) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (i *FeedItem) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(i,
		validation.Field(&i.FeedType, validation.Required),
		validation.Field(&i.FeedUUID, validation.Required),
		validation.Field(&i.PostUUID, validation.Required),
		validation.Field(&i.Tag, validation.Required, validation.In("POST", "LIKE", "COMMENT", "SHARE")))

	if err != nil {
		db.AddError(err)
	}
}
