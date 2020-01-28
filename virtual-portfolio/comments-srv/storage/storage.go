package storage

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Service is an wrapper of the comments-srv db
type Service interface {
	Get(Comment) (Comment, error)
	Delete(string) error
	Create(Comment) (Comment, error)
	GetResource(Resource) (Resource, error)
	ListResources(string, []string) ([]*Resource, error)
	Close() error
}

// Resource is a object which users can comment on
type Resource struct {
	UUID     string
	Type     string
	Comments []Comment
}

// Comment is a user-generated comment
type Comment struct {
	UUID         string `gorm:"type:uuid;primary_key;"`
	Text         string `gorm:"type:text"`
	ResourceUUID string `gorm:"type:varchar(100)"`
	ResourceType string `gorm:"type:varchar(100)"`
	UserUUID     string `gorm:"type:varchar(100)"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (c *Comment) BeforeCreate(scope *gorm.Scope) error {
	u := uuid.NewV4()
	return scope.SetColumn("UUID", u.String())
}

// BeforeSave performs the validations
func (c *Comment) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(c,
		validation.Field(&c.ResourceUUID, validation.Required),
		validation.Field(&c.ResourceType, validation.Required),
		validation.Field(&c.UserUUID, validation.Required),
		validation.Field(&c.Text, validation.Required))

	if err != nil {
		db.AddError(err)
		return
	}
}
