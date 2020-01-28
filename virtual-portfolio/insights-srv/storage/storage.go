package storage

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/gorm"
)

// Service is an wrapper of the followers-srv db
type Service interface {
	CreateInsight(Insight) (Insight, error)
	CreateUserView(UserView) (UserView, error)

	ListAssets(time.Time, bool) ([]Asset, error)
	ListInsightsForUser(string, time.Time) ([]Insight, error)
	// ListInsightsForAssets(string, []string, time.Time) ([]Insight, error)
	GetUserView(string, Asset, time.Time) (UserView, error)

	Close() error
}

// Asset is a resource an insight can be about
type Asset struct {
	UUID, Type string
}

// Insight is an insight into an asset for one or more usersx
type Insight struct {
	gorm.Model

	AssetUUID string `json:"asset_uuid" gorm:"type:uuid;"`
	AssetType string `json:"asset_type" gorm:"type:varchar(20)"`
	UserUUID  string `json:"user_uuid"`
	PostUUID  string `json:"post_uuid"`
	LinkURL   string `json:"link_url"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Type      string `json:"type"`
}

// UserView is a record of a user viewing the insights for an asset
type UserView struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	UserUUID  string `gorm:"type:uuid"`
	AssetUUID string `gorm:"type:uuid"`
	AssetType string `gorm:"type:varchar(20)"`
}

// BeforeSave performs the validations
func (i *Insight) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(i,
		validation.Field(&i.AssetUUID, validation.Required),
		validation.Field(&i.AssetType, validation.Required),
		validation.Field(&i.Title, validation.Required),
		validation.Field(&i.Type, validation.Required, validation.In("NEWS", "POST", "PRICE_MOVEMENT", "TRADE", "EVENT")))

	if err != nil {
		db.AddError(err)
	}
}

// BeforeSave performs the validations
func (u *UserView) BeforeSave(db *gorm.DB) {
	err := validation.ValidateStruct(u,
		validation.Field(&u.AssetUUID, validation.Required),
		validation.Field(&u.AssetType, validation.Required),
		validation.Field(&u.UserUUID, validation.Required))

	if err != nil {
		db.AddError(err)
	}
}
