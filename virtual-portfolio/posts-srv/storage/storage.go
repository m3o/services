package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Service is an wrapper of the followers-srv db
type Service interface {
	Create(Post) (Post, error)
	Update(Post) (Post, error)
	Get(Post) (Post, error)
	Count(Post) (int32, error)
	CountByUser([]string, time.Time, time.Time) (map[string]int32, error)
	List([]string) ([]*Post, error)
	Recent(int32, int32) ([]*Post, error)
	ListFeed(string, string) ([]*Post, error)
	ListUser(string, int32, int32) ([]*Post, error)
	Delete(Post) error
	Close() error
}

// Post is a user generated 'message'
type Post struct {
	gorm.Model

	UUID                string `gorm:"type:uuid;primary_key;"`
	Text                string `gorm:"type:text"`
	Title               string `gorm:"type:text"`
	UserUUID            string `gorm:"type:varchar(100)"`
	FeedType            string `gorm:"type:varchar(100)"`
	FeedUUID            string `gorm:"type:varchar(100)"`
	AttachmentLinkURL   string `gorm:"type:text"`
	AttachmentPictureID string `gorm:"type:varchar(100)"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (p *Post) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (p *Post) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Text, validation.Required),
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.UserUUID, validation.Required),
		validation.Field(&p.FeedType, validation.Required),
		validation.Field(&p.FeedUUID, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
